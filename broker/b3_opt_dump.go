package main

import (
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	mathrand "math/rand"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"bufio"
	"os/exec"

	vegeta "github.com/tsenart/vegeta/lib"
)

func GetKubeClient(confPath string) *kubernetes.Clientset {

	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: confPath}
	loader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})

	clientConfig, err := loader.ClientConfig()
	if err != nil {
		panic(err)
	}

	kubeclient, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		panic(err)
	}
	return kubeclient
}

func tokenGenerator() string {
	b := make([]byte, 6)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func dump(info, file string) {

	f, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	check(err)
	defer f.Close()

	w := bufio.NewWriter(f)

	w.WriteString(info)
	w.Flush()
}

var (
	wg sync.WaitGroup
	layout    = "2006-01-02 15:04:05 +0000 UTC"
	clientset = GetKubeClient("/root/admin.conf")
)

func main() {

	start := time.Now()
	argsWithoutProg := os.Args[1:]

	inputFile := string(argsWithoutProg[0])
	var experimentDuration = 150 //time.Duration(150) * time.Second

	fmt.Printf("INPUTFILE: %s\n", inputFile)
	fmt.Printf("SECONDS: %s\n", argsWithoutProg[1])

	if len(argsWithoutProg) > 1 {
		experimentDurationInt, _ := strconv.Atoi(argsWithoutProg[1])
		experimentDuration = experimentDurationInt //time.Duration(experimentDurationInt) * time.Second
	}

	var serviceAttackRate = 50
	if len(argsWithoutProg) > 2 {
		serviceAttackRateInt, _ := strconv.Atoi(argsWithoutProg[2])
		serviceAttackRate = serviceAttackRateInt
		fmt.Printf("Service Attack rate (per instance): %s\n", argsWithoutProg[2])
	}

	var factorialNumber = 900000
	if len(argsWithoutProg) > 3 {
		factorialNumberInt, _ := strconv.Atoi(argsWithoutProg[3])
		factorialNumber = factorialNumberInt
		fmt.Printf("Factorial number to be computed: %s\n", argsWithoutProg[3])
	}

	var timeRef = 0
	file, err := ioutil.ReadFile(inputFile)

	if err == nil {

		r := csv.NewReader(strings.NewReader(string(file)))

		for {

			//time.Sleep(time.Duration(1) * time.Second)
			//time.Sleep(time.Duration(500) * time.Millisecond)
			record, err := r.Read()

			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			if string(record[0]) != "jid" {

				timestamp, _ := strconv.Atoi(string(record[0]))
				slo := string(record[11])
				cpuReq := string(record[8])
				memReq := string(record[9])
				taskID := string(record[1])
				class := string(record[10])

				controllerKind := string(record[12])
				replicaOrParallelism := string(record[13])
				completions := string(record[14])
				qosMeasuring := string(record[15])

				var qosMeasuringAux string
				if qosMeasuring == "time_aggregated" {
					qosMeasuringAux = "timeaggregated"
				} else if qosMeasuring == "task_aggregated" {
					qosMeasuringAux = "taskaggregated"
				} else {
					qosMeasuringAux = qosMeasuring
				}

				controllerName := class + "-" + controllerKind  + "-" + qosMeasuringAux + "-" + taskID + "-" + tokenGenerator()
				//dump(controllerName+"\n", "controllers.csv")



				replicaOrParallelismInt, _ := strconv.Atoi(replicaOrParallelism)
				replicaOrParallelismInt32 := int32(replicaOrParallelismInt)

				completionsInt, _ := strconv.Atoi(completions)
				completionsInt32 := int32(completionsInt)



				if controllerKind == "deployment" {
					//port, _ := strconv.Atoi(string(record[16]))
					nodePort,_ := strconv.Atoi(string(record[16]))

					deployment := getDeploymentSpec(controllerName, cpuReq, memReq, slo, class, &replicaOrParallelismInt32, qosMeasuring) // int32(port))

					service := getServiceOfDeployment(deployment,  int32(nodePort))
					fmt.Println("Reading new task...")
					fmt.Printf("Deployment %s, cpu: %v, mem: %v, slo: %s, replicas: %d, qosMeasuring: %s, nodePort %d\n", controllerName, cpuReq, memReq, slo, replicaOrParallelism, qosMeasuring, nodePort)

					fmt.Printf("Creating service %s-service\n", controllerName)
					_, err = clientset.CoreV1().Services("default").Create(service)
					if err != nil {
						fmt.Printf("[ERROR] While creating service: \n", err)
					}

					if timestamp == timeRef {
						fmt.Println("Time: ", timestamp)
						fmt.Println("Creating deployment ", controllerName)
						wg.Add(1)
						go CreateAndManageDeploymentTermination(controllerName, deployment, service, experimentDuration-timeRef, &wg)

					} else {
						waittime := int(timestamp - timeRef)
						timeRef = timestamp
						fmt.Println("")
						time.Sleep(time.Duration(waittime) * time.Second)
						fmt.Println("Time: ", timestamp)
						fmt.Println("Creating deployment ", controllerName)
						wg.Add(1)
						go CreateAndManageDeploymentTermination(controllerName, deployment, service, experimentDuration-timeRef, &wg)
					}

					wg.Add(1)
					//frequency := 100
					//go CreateAndManageDemand(controllerName, nodePort, frequency, replicaOrParallelismInt, experimentDuration-timeRef, &wg)
					go CreateAndManageDemand(controllerName, nodePort, serviceAttackRate, replicaOrParallelismInt, experimentDuration-timeRef, factorialNumber, &wg)

				} else if controllerKind == "job" {

					var taskDuration string

					fmt.Println("fields %d", r.FieldsPerRecord)

					if r.FieldsPerRecord > 17 {
						taskDuration = string(record[17])
					} else {
						taskDuration = "25"
					}

					job := getJobSpec(controllerName, cpuReq, memReq, slo, class, &replicaOrParallelismInt32, &completionsInt32, qosMeasuring, taskDuration)

					//job := getJobSpec(controllerName, cpuReq, memReq, slo, class, &replicaOrParallelismInt32, &completionsInt32, qosMeasuring)

					fmt.Println("Reading new task...")
					fmt.Printf("Job %s, cpu: %v, mem: %v, slo: %s, paralelism: %d, completions: %d, qosMeasuring: %s, taskDuration %s\n", controllerName, cpuReq, memReq, slo, replicaOrParallelism, completions, qosMeasuring, taskDuration)

					if timestamp == timeRef {
						fmt.Println("Time: ", timestamp)
						fmt.Println("Creating job ", controllerName)
						wg.Add(1)
						go CreateAndManageJobTermination(controllerName, job, experimentDuration-timeRef, &wg)

					} else {
						waittime := int(timestamp - timeRef)
						timeRef = timestamp
						fmt.Println("")
						time.Sleep(time.Duration(waittime) * time.Second)
						fmt.Println("Time: ", timestamp)
						fmt.Println("Creating job ", controllerName)
						wg.Add(1)
						go CreateAndManageJobTermination(controllerName, job, experimentDuration-timeRef, &wg)
					}
				}


			}
		}
	}

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Println("Finished - runtime: ", elapsed)
}

func CreateAndManageDeploymentTermination(controllerName string, deployment *appsv1beta2.Deployment, service *v1.Service, expectedRuntime int, wg *sync.WaitGroup) {

	_, err := clientset.AppsV1beta2().Deployments("default").Create(deployment)
	if err != nil {
		fmt.Printf("[ERROR] While creating deployment: ", err)
	}


	//_, err = clientset.CoreV1().Services("default").Create(service)
	//if err != nil {
	//	fmt.Printf("[ERROR] While creating service: ", err)
	//}

	time.Sleep(time.Duration(expectedRuntime) * time.Second)

	fmt.Println("Killing all deployments and respective services after ", expectedRuntime, "seconds")
	//out := fmt.Sprintf("Killing all deployments and services after %s", expectedRuntime)
	//dump(out, "/root/broker.log")

	//clientset.AppsV1beta2().Deployments("default").Delete(controllerName, &metav1.DeleteOptions{})
	cmd := exec.Command("/usr/bin/kubectl", "delete", "deploy", controllerName)
	cmd.Run()


	cmd = exec.Command("/usr/bin/kubectl", "delete", "service", service.Name)
	cmd.Run()

	wg.Done()
}

func CreateAndManageDemand(controllerName string, nodePort int, frequency int, numberOfReplicas int, expectedRuntime int, factorialNumber int, wg *sync.WaitGroup) {

	//time.Sleep(time.Duration(5) * time.Second) ate o dia 08/05 este valor que estava sendo considerado para warme-up do deployment
	//time.Sleep(time.Duration(20) * time.Second)

	var wgDump sync.WaitGroup

	time.Sleep(time.Duration(60) * time.Second)

	fmt.Printf("Creating the demmand to service on port %d\n", nodePort)

	//time.Sleep(time.Duration(expectedRuntime) * time.Second)
	freq := frequency * numberOfReplicas
	duration :=  time.Duration(expectedRuntime - 120) * time.Second

	rate := vegeta.Rate{Freq: freq, Per: time.Second}
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		//URL:    "http://localhost:" + strconv.Itoa(nodePort) + "/",
		URL:    "http://localhost:" + strconv.Itoa(nodePort) + "/factorial.php" + "?number="+ strconv.Itoa(factorialNumber),
		//URL:    "http://10.11.5.56:30555/factorial.php?number=900000",
	})
	//attacker := vegeta.NewAttacker()
	attacker := vegeta.NewAttacker(vegeta.Timeout(2 * time.Minute), vegeta.Connections(1000))

	var metrics vegeta.Metrics
	var out string
	//var outB string
	//var latencies = []float64{}
	var latencies = []string{}
	//var errors = []string{}

	var intermediaryDump = 0

	for res := range attacker.Attack(targeter, rate, duration, fmt.Sprintf("app %s on port %d\n", controllerName, nodePort)) {
		//fmt.Printf("Error: %v\n", res.Error)
		//fmt.Printf("latency: %f\n", res.Latency.Seconds())
		//out += fmt.Sprintf("%f\n", res.Latency.Seconds())
		//dump(out, "/home/" + controllerName + ".latency")
		metrics.Add(res)
		latencyLine := fmt.Sprintf("%d,%f", res.Timestamp.Unix(), res.Latency.Seconds())
		//latencies = append(latencies, res.Latency.Seconds())
		latencies = append(latencies, latencyLine)
		//if (res.Code != 200) {
		//        oneError := strconv.Itoa(int(res.Timestamp.Unix())) + "," + res.URL + "," + strconv.Itoa(int(res.Code))
		//        errors = append(errors, oneError)
		//}

		//TODO fazer dump das latencias paralelamente quando a quantidade atingir determinado threshold (wg.add e wd.done para esse dump) arquivos diferentes no dump para testar a funcionalidade?
		if (len(latencies) == 100000) {
			intermediaryDump++
			var intermediaryLatencies = latencies
			latencies = []string{}

			wgDump.Add(1)
			go DumpIntermediaryResults(controllerName, nodePort, intermediaryDump, intermediaryLatencies, &wgDump)

		}
	}
	metrics.Close()

	delaying :=  mathrand.Intn(120)

	fmt.Printf("Finishing attack on port: %s. Delaying %d seconds to dump the latencies. \n", strconv.Itoa(nodePort), delaying)

	time.Sleep(time.Duration(delaying) * time.Second)

	fmt.Printf(" Dumping the LAST latencies of service running on port %s. \n", strconv.Itoa(nodePort))

	intermediaryDump++
	//dump(out, "/home/" + controllerName + ".latency")
	for _, value := range latencies {
		out += fmt.Sprintf("%v\n", value)
	}
	dump(out, "/home/" + controllerName + ".dump" + strconv.Itoa(intermediaryDump) + ".latency")

	//fmt.Printf("Finishing attack on port: %s. Dumping the errors... \n", strconv.Itoa(nodePort))

	//dump(out, "/home/" + controllerName + ".latency")
	//for _, value := range errors {
	//        outB += fmt.Sprintf("%v\n", value)
	//}
	//dump(outB, "/home/" + controllerName + ".errors")


	fmt.Printf("Service running on port: %s\n", strconv.Itoa(nodePort))
	fmt.Printf("Mean: %s\n", metrics.Latencies.Mean)
	fmt.Printf("50th percentile: %s\n", metrics.Latencies.P50)
	fmt.Printf("90th percentile: %s\n", metrics.Latencies.P90)
	fmt.Printf("95th percentile: %s\n", metrics.Latencies.P95)
	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
	fmt.Printf("Max percentile: %s\n", metrics.Latencies.Max)
	//fmt.Printf("Total: %s\n", metrics.Latencies.Total)
	fmt.Printf("Requests: %v\n", metrics.Requests)
	fmt.Printf("Success: %f\n", metrics.Success)
	fmt.Printf("Throughtput: %f\n", metrics.Throughput)

	// Mean, 50th percentile, 90th percentile, 95th percentile, 99th percentile, Max percentile, Total
	//out := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s\n", controllerName, metrics.Latencies.Mean,
	//	metrics.Latencies.P50, metrics.Latencies.P90, metrics.Latencies.P95, metrics.Latencies.P99,
	//	metrics.Latencies.Max, metrics.Latencies.Total)

	out = fmt.Sprintf("%s,%f,%f,%f,%f,%f,%f,%v,%f,%f\n", controllerName, metrics.Latencies.Mean.Seconds(),
		metrics.Latencies.P50.Seconds(), metrics.Latencies.P90.Seconds(), metrics.Latencies.P95.Seconds(), metrics.Latencies.P99.Seconds(),
		metrics.Latencies.Max.Seconds(), metrics.Requests, metrics.Success, metrics.Throughput)

	fmt.Printf("Service running on port: %s - %s\n", strconv.Itoa(nodePort), out)

	//dump(out, "/home/service-application.dat")
	dump(out, "/home/serviceapplication.dat")

	// wait the intermediary dump of the latencies before conclude
	wgDump.Wait()

	wg.Done()
}

func DumpIntermediaryResults(controllerName string, nodePort int,  intermediaryDump int, intermediaryLatencies []string, wgDump *sync.WaitGroup) {

	// TODO fazer um sleep aleatorio aqui?
	delaying :=  mathrand.Intn(120)

        fmt.Printf("Delaying %d seconds to dump the latencies (round %d, service on %d). \n", delaying, intermediaryDump, nodePort)

        time.Sleep(time.Duration(delaying) * time.Second)

        fmt.Printf(" Dumping latencies (round %d, service on %d). \n", intermediaryDump, nodePort)

	var intermediaryOut string

	for _, value := range intermediaryLatencies {
		intermediaryOut += fmt.Sprintf("%v\n", value)
	}
	dump(intermediaryOut, "/home/" + controllerName + ".dump" + strconv.Itoa(intermediaryDump) + ".latency" )

	wgDump.Done()
}

func CreateAndManageJobTermination(controllerName string, job *batchv1.Job, expectedRuntime int, wg *sync.WaitGroup) {

	_, err := clientset.BatchV1().Jobs("default").Create(job)
	if err != nil {
		fmt.Printf("[ERROR] Error while creating job information: ", err)
	}

	time.Sleep(time.Duration(expectedRuntime) * time.Second)

	fmt.Println("Killing all jobs after ", expectedRuntime, "seconds")
	//out := fmt.Sprintf("Killing all jobs  after %s", expectedRuntime)
	//dump(out, "/root/broker.log")

	cmd := exec.Command("/usr/bin/kubectl", "delete", "job", controllerName)
	cmd.Run()

	wg.Done()
}


func getDeploymentSpec(controllerRefName string,
	cpuReq string, memReq string, slo string, class string,
	replicaOrParallelism *int32, qosMeasuring string) *appsv1beta2.Deployment {

	memReqFloat, _ := strconv.ParseFloat(memReq, 64)
	memReqKi := memReqFloat * 1000000
	memReqStr := strconv.FormatFloat(memReqKi, 'f', -1, 64)
	memRequest := memReqStr + "Ki"
	fmt.Println(memRequest)
	rl := v1.ResourceList{v1.ResourceName(v1.ResourceMemory): resource.MustParse(memRequest),
		v1.ResourceName(v1.ResourceCPU): resource.MustParse(cpuReq)}

	priorityClass := class + "-priority"
	gracePeriod := int64(0)



	pod := v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "nginx", "controller": controllerRefName}, Annotations: map[string]string{"slo": slo, "controller": controllerRefName, "qos_measuring": qosMeasuring}},
		//ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "golang"}, Annotations: map[string]string{"slo": slo, "controller": controllerRefName, "qos_measuring": qosMeasuring}},
		//ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "couchbase"}, Annotations: map[string]string{"slo": slo, "controller": controllerRefName, "qos_measuring": qosMeasuring}},
		Spec: v1.PodSpec{
			TerminationGracePeriodSeconds: &gracePeriod,
			Containers: []v1.Container{{Name: controllerRefName,
				Image: "viniciusbds/app:v0.0.2",
				//Image: "nginx:1.15",
				//Image: "golang:1.11",
				//Image: "couchbase:6.0.1",
				//Ports: []v1.ContainerPort{{ContainerPort:port}},
				// TODO this command property must be used only for golang and couchbase deploys
				//Command: []string{"/bin/bash", "-ce", "tail -f /dev/null" },
				Resources: v1.ResourceRequirements{
					Limits:   rl,
					Requests: rl,
				}}},
			//},
			//ImagePullPolicy: v1.PullAlways}},

			PriorityClassName: priorityClass},
	}


	deployment := &appsv1beta2.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: controllerRefName},
		Spec:       appsv1beta2.DeploymentSpec{Selector: &metav1.LabelSelector{MatchLabels: pod.Labels}, Replicas: replicaOrParallelism, Template: pod},
	}


	fmt.Print("Image Pull Police: ")

	imagePullPoliceOfDeployment := deployment.Spec.Template.Spec.Containers[0].ImagePullPolicy

	if imagePullPoliceOfDeployment == "" {
		fmt.Println("Policy Default !")
	} else {
		fmt.Println(imagePullPoliceOfDeployment)
	}

	return deployment
}


func getJobSpec(controllerRefName string,
	cpuReq string, memReq string, slo string, class string,
	replicaOrParallelism *int32, completions *int32, qosMeasuring string, taskDuration string) *batchv1.Job {

	memReqFloat, _ := strconv.ParseFloat(memReq, 64)
	memReqKi := memReqFloat * 1000000
	memReqStr := strconv.FormatFloat(memReqKi, 'f', -1, 64)
	memRequest := memReqStr + "Ki"
	fmt.Println(memRequest)
	rl := v1.ResourceList{v1.ResourceName(v1.ResourceMemory): resource.MustParse(memRequest),
		v1.ResourceName(v1.ResourceCPU): resource.MustParse(cpuReq)}

	priorityClass := class + "-priority"
	gracePeriod := int64(0)

	podSpec := v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "nginx"}, Annotations: map[string]string{"slo": slo, "controller": controllerRefName, "qos_measuring": qosMeasuring}},
		//ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "golang"}, Annotations: map[string]string{"slo": slo, "controller": controllerRefName, "qos_measuring": qosMeasuring}},
		//ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "couchbase"}, Annotations: map[string]string{"slo": slo, "controller": controllerRefName, "qos_measuring": qosMeasuring}},
		Spec: v1.PodSpec{
			TerminationGracePeriodSeconds: &gracePeriod,
			Containers: []v1.Container{{Name: controllerRefName,
				Image: "nginx:1.15",
				//Image: "golang:1.11",
				//Image: "couchbase:6.0.1",
				//Command: []string{"sleep", "1800" },
				//Command: []string{"sleep", "25" },
				Command: []string{"sleep", taskDuration },
				Resources: v1.ResourceRequirements{
					Limits:   rl,
					Requests: rl,
				}}},
			//},
			//ImagePullPolicy: v1.PullAlways}},
			PriorityClassName: priorityClass,
			RestartPolicy: "Never"},
	}


	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{Name: controllerRefName, Namespace: "default"},

		Spec: batchv1.JobSpec{
			//Selector: &metav1.LabelSelector{MatchLabels: pod.Labels},
			Template: podSpec,
			Parallelism: replicaOrParallelism,
			Completions: completions,
		},

	}

	job.Spec.Template.Labels = map[string]string{"app": "nginx"}
	//job.Spec.Template.Labels = map[string]string{"app": "golang"}
	//job.Spec.Template.Labels = map[string]string{"app": "couchbase"}

	fmt.Print("Image Pull Police: ")
	imagePullPoliceOfDeployment := job.Spec.Template.Spec.Containers[0].ImagePullPolicy

	if imagePullPoliceOfDeployment == "" {
		fmt.Println("Policy Default !")
	} else {
		fmt.Println(imagePullPoliceOfDeployment)
	}

	return job
}



func getServiceOfDeployment(deployment *appsv1beta2.Deployment, nodePort int32) *v1.Service {
	servicePort := v1.ServicePort{Port: 80, NodePort: nodePort}
	serviceSpec := v1.ServiceSpec{
		Ports: []v1.ServicePort{servicePort},
		Selector: deployment.Spec.Template.Labels,
		Type: v1.ServiceTypeNodePort,
	}

	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: deployment.Name + "-service", Namespace: "default"},
		Spec: serviceSpec}
	return service
}
