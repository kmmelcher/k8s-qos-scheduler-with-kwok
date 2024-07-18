package qos_driven_scheduler

import (
	"context"
	"fmt"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	frameworkruntime "k8s.io/kubernetes/pkg/scheduler/framework/runtime"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

const (
	Name                                = "QosAware"
	DefaultSlo                          = 0.80
	DefaultAcceptablePreemptionOverhead = 1.00
)

type QosDrivenScheduler struct {
	fh          framework.FrameworkHandle
	Args        QosDrivenSchedulerArgs
	PodInformer cache.SharedIndexInformer
	Controllers map[string]ControllerMetricInfo
	lock        sync.RWMutex
}

var _ framework.QueueSortPlugin = &QosDrivenScheduler{}
var _ framework.ReservePlugin = &QosDrivenScheduler{}
var _ framework.PreBindPlugin = &QosDrivenScheduler{}
var _ framework.PostBindPlugin = &QosDrivenScheduler{}
var _ framework.PostFilterPlugin = &QosDrivenScheduler{}

func (scheduler *QosDrivenScheduler) Name() string {
	return Name
}

func (scheduler *QosDrivenScheduler) GetControllerMetricInfo(pod *core.Pod) ControllerMetricInfo {
	scheduler.lock.RLock()
	defer scheduler.lock.RUnlock()
	cMetricInfo := scheduler.Controllers[ControllerName(pod)]
	cMetricInfo.ReferencePod = pod
	return cMetricInfo
}

// getUpdatedVersion returns the latest pod object in cache
func (scheduler *QosDrivenScheduler) getUpdatedVersion(pod *core.Pod) PodMetricInfo {
	scheduler.lock.RLock()
	defer scheduler.lock.RUnlock()

	cMetricInfo := scheduler.Controllers[ControllerName(pod)]
	pMetricInfo, _ := cMetricInfo.GetPodMetricInfo(pod)

	return pMetricInfo
}

// Atomically updates the PodMetricInfo.
// The mapper f should return a PodMetricInfo based on it's previous state, besides that f also
// returns if the pod is a new running pod or a pod running that is being deleted.
func (scheduler *QosDrivenScheduler) UpdatePodMetricInfo(pod *core.Pod, f func(PodMetricInfo) PodMetricInfo) {
	scheduler.lock.Lock()
	defer scheduler.lock.Unlock()

	if pod.Namespace == "kube-system" {
		klog.V(1).Infof("Pod %s is from kube-system namespace. There is no need to update podMetricInfo.", pod.Name)
		return
	}

	controllerName := ControllerName(pod)
	podName := PodName(pod)

	if scheduler.Controllers == nil {
		scheduler.Controllers = map[string]ControllerMetricInfo{}
	}

	cMetricInfo, found := scheduler.Controllers[controllerName]
	if !found {
		cMetricInfo.SafetyMargin = scheduler.Args.SafetyMargin.Duration
		cMetricInfo.MinimumRunningTime = scheduler.Args.MinimumRunningTime.Duration
		cMetricInfo.QoSMeasuringApproach = ControllerQoSMeasuring(pod)

		// TODO remove code, it is only for debugging
		kubeclient := scheduler.fh.ClientSet()

		namespace := "default"

		var numberOfReplicas int32
		var expectedPodCompletions int32

		if controllerRef := metav1.GetControllerOf(pod); controllerRef == nil {
			klog.V(1).Infof("pod %s has no controllerRef", pod.Name)
			numberOfReplicas = 1
			expectedPodCompletions = 0
		} else {

			if controllerRef.Kind == "Job" {
				klog.V(1).Infof("Controller is a JOB! ControllerRefName of pod %s is %s ", pod.Name, controllerRef.Name)
				controller, err := kubeclient.BatchV1().Jobs(namespace).Get(context.TODO(), controllerRef.Name, metav1.GetOptions{})

				if err != nil {
					klog.V(1).Infof("ERROR while getting info about Job %v", err)
				} else {
					numberOfReplicas = *controller.Spec.Parallelism
					expectedPodCompletions = *controller.Spec.Completions

					klog.V(1).Infof("The pod %s is associated with Job %s and its parallelism is %v", pod.Name, controllerRef.Name, numberOfReplicas)
				}

			} else {
				klog.V(1).Infof("Controller is NOT a JOB! ControllerRefName of pod %s is %s ", pod.Name, controllerRef.Name)
				controller, err := kubeclient.AppsV1().ReplicaSets(namespace).Get(context.TODO(), controllerRef.Name, metav1.GetOptions{})

				if err != nil {
					klog.V(1).Infof("ERROR while getting info about ReplicaSet %v", err)
				} else {
					numberOfReplicas = *controller.Spec.Replicas
					expectedPodCompletions = 0
					klog.V(1).Infof("The pod %s is associated with ReplicaSet %s and its number of replicas is %v", pod.Name, controllerRef.Name, numberOfReplicas)
				}
			}
		}
		cMetricInfo.NumberOfReplicas = numberOfReplicas
		cMetricInfo.ExpectedPodCompletions = expectedPodCompletions
	}

	// initializing replicas data structure if the cMetricInfo is a new one
	if cMetricInfo.replicas == nil {
		for n := int32(0); n < cMetricInfo.NumberOfReplicas; n++ {
			cMetricInfo.replicas = append(cMetricInfo.replicas, &Replica{
				Incarnations:       map[string]PodMetricInfo{},
				CurrentIncarnation: "",
			})
		}
	}

	oldPodMetricInfo, isNewPod := cMetricInfo.GetPodMetricInfo(pod)
	newPodMetricInfo := f(oldPodMetricInfo)

	// Check if this pod is being deleted, if yes, it needs to realease the replicaId of its controller
	wasMarkedForDeletionNow := false

	if (oldPodMetricInfo.LastStatus != nil) && (newPodMetricInfo.LastStatus != nil) {
		wasMarkedForDeletionNow = (oldPodMetricInfo.LastStatus.DeletionTimestamp == nil) && (newPodMetricInfo.LastStatus.DeletionTimestamp != nil)
	}

	notMarkedForDeletion := (newPodMetricInfo.LastStatus == nil) || (newPodMetricInfo.LastStatus.DeletionTimestamp == nil)
	wasDeletedFromCacheNow := oldPodMetricInfo.EndTime.IsZero() && !newPodMetricInfo.EndTime.IsZero()

	// When the pod is marked for deletion (terminating) we should deallocate it's replicaId.
	// But we can miss this state. In this case, we deallocate when the informer informs pod's deletion
	deallocate := !newPodMetricInfo.IsSucceeded && (wasMarkedForDeletionNow || (notMarkedForDeletion && wasDeletedFromCacheNow))

	// Check if this pod is a new succeeded pod. Is yes, it needs to release the replicaId and increment the SucceededPods variable
	if (oldPodMetricInfo.LastStatus == nil || oldPodMetricInfo.LastStatus.Status.Phase != core.PodSucceeded) &&
		((newPodMetricInfo.LastStatus != nil) && (newPodMetricInfo.LastStatus.Status.Phase == core.PodSucceeded)) {
		klog.V(1).Infof("The pod %s is being successfully completed", PodName(pod))
		cMetricInfo.SucceededPods++
		newPodMetricInfo.IsSucceeded = true
		newPodMetricInfo.EndTime = time.Now()
		deallocate = true
	}

	// If the pod was deleted or succeeded, dump it and disassociate it from the replica's currentIncarnation
	if deallocate {
		cMetricInfoClone := cMetricInfo.Clone()
		cMetricInfo.DeallocateReplicaId(pod)

		cMetricInfoClone.ReferencePod = pod
		go dumpMetrics(cMetricInfoClone, newPodMetricInfo, scheduler.lock.RLocker())
	}

	// This is a new pod and it needs to be associated with a replicaId of the controller
	if isNewPod {
		replicaId := cMetricInfo.AllocateReplicaId(pod)
		newPodMetricInfo.ReplicaId = replicaId
		klog.V(1).Infof("the pod %s was new, allocated replicaId = %d", PodName(pod), replicaId)
	}

	previouslyRunningPods := cMetricInfo.NumberOfRunningPods()
	replicaId, err := cMetricInfo.GetPodReplicaId(pod)
	if err != nil {
		klog.V(1).Infof("if the pod was new, it should have been allocated (see variable isNewPod)")
		panic(err)
	}

	cMetricInfo.replicas[replicaId].Incarnations[podName] = newPodMetricInfo
	scheduler.Controllers[controllerName] = cMetricInfo
	klog.V(1).Infof("The number of pods running of %s was %d and now is %d", controllerName, previouslyRunningPods, cMetricInfo.NumberOfRunningPods())
	klog.V(1).Infof("The number of succeeded pods of %s is %d", controllerName, cMetricInfo.SucceededPods)
}

func dumpMetrics(cMetricInfo ControllerMetricInfo, pMetricInfo PodMetricInfo, accessReaderLocker sync.Locker) {
	podName := PodName(pMetricInfo.LastStatus)
	controllerName := ControllerName(pMetricInfo.LastStatus)
	klog.V(1).Infof("Dumping metrics of pod %s (status %s | controller %s)", podName, pMetricInfo.LastStatus.Status.Phase, controllerName)

	success := false
	timeRef := time.Now()
	cMetrics := cMetricInfo.Metrics(timeRef, accessReaderLocker)

	// timestamp, podName, controllerName, replicaId, qosMeasuring, waitingTime, allocationTime, runningTime, terminationStatus, controllerQoS, qosMetric
	entry := fmt.Sprintf("%d,%s,%s,%f,%d,%s,%d,%d,%d,%d,%t,%f,%f,%d,%d,%d,%d\n",
		time.Now().Unix(),
		podName,
		controllerName,
		ControllerSlo(pMetricInfo.LastStatus),
		pMetricInfo.ReplicaId,
		cMetrics.QoSMeasuringApproach,
		pMetricInfo.CreationTime.Unix(),
		pMetricInfo.Waiting(timeRef).Milliseconds(),
		pMetricInfo.Binding(timeRef).Milliseconds(),
		pMetricInfo.Running(timeRef).Milliseconds(),
		pMetricInfo.IsSucceeded,
		cMetrics.QoS(),
		cMetrics.QoSMetric(pMetricInfo.LastStatus),
		cMetrics.WaitingTime.Milliseconds(),
		cMetrics.BindingTime.Milliseconds(),
		cMetrics.EffectiveRunningTime.Milliseconds(),
		cMetrics.DiscardedRunningTime.Milliseconds())

	for !success {
		f, err := os.OpenFile("/metrics.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err == nil {
			_, err = f.WriteString(entry)

			if err == nil {
				success = true
			}
		}

		if !success {
			klog.V(1).Infof("Failed to log pod metrics: %s, retrying...\n", err)
			time.Sleep(time.Second)
		}
	}
}

// TODO maybe we should sort only by controller's SLO and importance if we can't find its metrics
func (scheduler *QosDrivenScheduler) WaitForPodsOnCache(pods ...*core.Pod) {
	for _, pod := range pods {
		cacheMiss := func() bool {
			cMetricInfo := scheduler.GetControllerMetricInfo(pod)
			scheduler.lock.RLock()
			_, notFound := cMetricInfo.GetPodMetricInfo(pod)
			scheduler.lock.RUnlock()
			return notFound
		}

		for cacheMiss() {
			time.Sleep(time.Millisecond * 100)
		}
	}
}

// We calculate precedence between pods checking if the time to violate p1's controller SLO is lower than p2's controller.
// If the pods' controllers has different importances and their time to violate are below the safety margin the precedence is for the highest importance controller's pod.
func (scheduler *QosDrivenScheduler) HigherPrecedence(p1, p2 *core.Pod) bool {
	now := time.Now()

	scheduler.WaitForPodsOnCache(p1, p2)
	cMetricInfo1 := scheduler.GetControllerMetricInfo(p1)
	cMetrics1 := cMetricInfo1.Metrics(now, scheduler.lock.RLocker())
	cMetricInfo2 := scheduler.GetControllerMetricInfo(p2)
	cMetrics2 := cMetricInfo2.Metrics(now, scheduler.lock.RLocker())

	importance1 := ControllerImportance(p1)
	importance2 := ControllerImportance(p2)

	qosMetric1 := cMetrics1.QoSMetric(p1)
	qosMetric2 := cMetrics2.QoSMetric(p2)

	safetyMargin := scheduler.Args.SafetyMargin.Duration.Seconds()

	// Is in resource contention
	if (qosMetric1 < safetyMargin) && (qosMetric2 < safetyMargin) && importance1 != importance2 {
		return importance1 > importance2
	}

	return qosMetric1 < qosMetric2
}

// This is the method required to the QueueSortPlugin interface.
// It is used to sort the pods on the ActiveQueue.
// Our implementation sorts them by increasing precedence, which is derived from TTV and importance.
func (scheduler *QosDrivenScheduler) Less(pInfo1, pInfo2 *framework.QueuedPodInfo) bool {
	return scheduler.HigherPrecedence(pInfo1.Pod, pInfo2.Pod)
}

// This is the method required to the ReservePlugin interface.
// Reverse method is called by the scheduling framework when resources on a node are being reserved for a given Pod.
func (scheduler *QosDrivenScheduler) Reserve(_ context.Context, _ *framework.CycleState, p *core.Pod, nodeName string) *framework.Status {
	now := time.Now()
	klog.V(1).Infof("resources reserved at node %s for pod %s at %s", nodeName, PodName(p), now)
	scheduler.UpdatePodMetricInfo(p, func(old PodMetricInfo) PodMetricInfo {
		old.AllocationStatus = AllocatingState
		return old
	})
	return nil
}

// This is the method required to the ReservePlugin interface.
// Unreserve is called by the scheduling framework when a reserved pod was
// rejected, an error occurred during reservation of subsequent plugins, or
// in a later phase.
func (scheduler *QosDrivenScheduler) Unreserve(_ context.Context, _ *framework.CycleState, p *core.Pod, nodeName string) {
	now := time.Now()
	klog.V(1).Infof("resources are being unreserved at node %s for pod %s at %s", nodeName, PodName(p), now)
	scheduler.UpdatePodMetricInfo(p, func(old PodMetricInfo) PodMetricInfo {
		old.AllocationStatus = ""
		return old
	})
}

// This is the method required to the PreBindPlugin interface.
// PreBind will run after a plugin Permit the pod, we will use it in combination
// with PostBind to calculate the binding time
func (scheduler *QosDrivenScheduler) PreBind(_ context.Context, state *framework.CycleState, p *core.Pod, _ string) *framework.Status {
	now := time.Now()
	state.Write(BindingStart, CloneableTime{now})
	return nil
}

// This is the method required to the PostBindPlugin interface.
// It will run if the binding ended successfully.
func (scheduler *QosDrivenScheduler) PostBind(_ context.Context, state *framework.CycleState, p *core.Pod, _ string) {
	start, err := state.Read(BindingStart)
	if err != nil {
		klog.Error(err.Error)
	}

	startBinding := start.(CloneableTime).Time
	endBinding := time.Now()

	scheduler.UpdatePodMetricInfo(p, func(old PodMetricInfo) PodMetricInfo {
		old.StartBindingTime = startBinding
		old.StartRunningTime = endBinding
		old.AllocationStatus = AllocatedState
		// the pod starts running at this moment
		return old
	})
}

func (scheduler *QosDrivenScheduler) OnAddPod(obj interface{}) {
	p := obj.(*core.Pod).DeepCopy()
	if p.Namespace == "kube-system" {
		return
	}

	start := p.GetCreationTimestamp().Time
	klog.V(1).Infof("Pods %s being added:\n%s", PodName(p), p.String())

	// CreationTimestamp is marked as +optional in the API, so we put this fallback
	if start.IsZero() {
		start = time.Now()
	}

	scheduler.UpdatePodMetricInfo(p, func(pMetricInfo PodMetricInfo) PodMetricInfo {
		pMetricInfo.ControllerName = ControllerName(p)
		pMetricInfo.LastStatus = p
		pMetricInfo.CreationTime = start

		// pod starts running at this moment
		if pMetricInfo.StartRunningTime.IsZero() && p.Status.Phase == core.PodRunning {
			pMetricInfo.StartRunningTime = time.Now()
		}

		return pMetricInfo
	})
	klog.V(1).Infof("Pods %s added successfully", PodName(p))
}

func (scheduler *QosDrivenScheduler) OnUpdatePod(_, newObj interface{}) {
	p := newObj.(*core.Pod).DeepCopy()
	if p.Namespace == "kube-system" {
		return
	}
	scheduler.UpdatePodMetricInfo(p, func(pMetricInfo PodMetricInfo) PodMetricInfo {
		pMetricInfo.LastStatus = p
		return pMetricInfo
	})
}

func (scheduler *QosDrivenScheduler) OnDeletePod(lastState interface{}) {
	p, ok := lastState.(*core.Pod)
	if !ok {
		klog.V(1).Infof("last state from deleted pod was unknown")
		return
	}
	p = p.DeepCopy()
	if p.Namespace == "kube-system" {
		return
	}
	klog.V(1).Infof("Pod %s deleted:\n%s", PodName(p), p.String())

	scheduler.UpdatePodMetricInfo(p, func(pMetricInfo PodMetricInfo) PodMetricInfo {
		// if pod was succeeded terminate, its endTime has already been set
		if pMetricInfo.IsSucceeded {
			return pMetricInfo
		}

		// pod is being deleted at this moment
		pMetricInfo.EndTime = time.Now()
		return pMetricInfo
	})
}

func (scheduler *QosDrivenScheduler) addEventHandler() *QosDrivenScheduler {
	scheduler.PodInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		// When a new pod gets created
		AddFunc: scheduler.OnAddPod,
		// When a pod gets updated
		UpdateFunc: scheduler.OnUpdatePod,
		// When a pod gets deleted
		DeleteFunc: scheduler.OnDeletePod,
	})
	return scheduler
}

func NewQosScheduler(fh framework.FrameworkHandle) *QosDrivenScheduler {

	return (&QosDrivenScheduler{
		PodInformer: fh.SharedInformerFactory().Core().V1().Pods().Informer(),
		fh:          fh,
		Controllers: map[string]ControllerMetricInfo{},
	}).addEventHandler()
}

type QosDrivenSchedulerArgs struct {
	// Controllers with a time-to-violate below this configured SafetyMargin
	// is treated as close to violate it's SLO by our scheduler.
	SafetyMargin metav1.Duration `json:"safetyMargin"`
	// AcceptablePreemptionOverhead is a global configuration to provide a default
	// acceptable preemption overhead to controllers without explicit one.
	AcceptablePreemptionOverhead float64 `json:"acceptablePreemptionOverhead,omitempty"`
	// Pods will run MinimumRunningTime until it can be preempted by another pod with same importance.
	MinimumRunningTime metav1.Duration `json:"minimumRunningTime"`
}

func format(d time.Duration) string {
	if d > time.Second {
		return d.Truncate(time.Second).String()
	}

	return d.Truncate(time.Millisecond).String()
}

func (scheduler *QosDrivenScheduler) list(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "%20s  %9s  %19s  %15s  %15s  %15s  %20s  %20s  %20s  %20s\n",
		"CONTROLLER",
		"QOS/SLO",
		"PREEMPTION OVERHEAD",
		"QOS APPROACH",
		"POD PHASE",
		"QOS METRIC",
		"EFF RUNNING",
		"DIS RUNNING",
		"WAITING",
		"BINDING")

	scheduler.lock.RLock()
	defer scheduler.lock.RUnlock()

	now := time.Now()
	cNames := sortedKeys(scheduler.Controllers)

	for _, cName := range cNames {
		controller := scheduler.Controllers[cName]

		for replicaId, replica := range controller.replicas {
			pod := replica.Incarnations[replica.CurrentIncarnation].LastStatus

			// listing only pods that are active in the system
			if pod != nil {
				// FIXME: scheduler.GetControllerMetricInfo also locks and this lock may cause deadlock when recursive locking
				cMetricInfo := scheduler.GetControllerMetricInfo(pod)
				metrics := cMetricInfo.Metrics(now, noopLocker)
				fmt.Fprintf(w, "%20s  %4.2f/%4.2f  %9.2f/%-9.2f  %15s  %15s  %15.4f  %20s  %20s  %20s  %20s\n",
					cName+"-"+strconv.Itoa(replicaId),
					metrics.QoS(),
					ControllerSlo(pod),
					metrics.PreemptionOverhead(),
					scheduler.AcceptablePreemptionOverhead(pod),
					cMetricInfo.QoSMeasuringApproach,
					pod.Status.Phase,
					metrics.QoSMetric(pod),
					format(metrics.EffectiveRunningTime),
					format(metrics.DiscardedRunningTime),
					format(metrics.WaitingTime),
					format(metrics.BindingTime))
			}
		}
	}
}

func sortedKeys(m map[string]ControllerMetricInfo) []string {
	keys := make([]string, len(m))
	idx := 0
	for key := range m {
		keys[idx] = key
		idx++
	}
	sort.Strings(keys)
	return keys
}

func (scheduler *QosDrivenScheduler) debugApi() {
	http.HandleFunc("/", scheduler.list)
	klog.Fatal(http.ListenAndServe(":10000", nil))
}

// The PluginFactory initializes a new plugin and returns it.
func New() func(args runtime.Object, f framework.FrameworkHandle) (framework.Plugin, error) {
	return func(args runtime.Object, f framework.FrameworkHandle) (framework.Plugin, error) {
		scheduler := NewQosScheduler(f)

		scheduler.Args.AcceptablePreemptionOverhead = DefaultAcceptablePreemptionOverhead
		if err := frameworkruntime.DecodeInto(args, &scheduler.Args); err != nil {
			return nil, err
		}

		klog.V(1).Infof("Plugin started with args %+v", scheduler.Args)
		go scheduler.debugApi()
		return scheduler, nil
	}
}
