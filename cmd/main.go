package main

import (
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
	"os"
	scheduler "qos-driven-scheduler/pkg/qos-driven-scheduler"
)

func main() {
	command := app.NewSchedulerCommand(
		app.WithPlugin(scheduler.Name, scheduler.New()))

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
