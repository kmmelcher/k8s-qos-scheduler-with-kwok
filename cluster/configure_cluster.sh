#!/bin/bash

# configuring the scheduling policy

# kube-scheduler.yaml file arguments
SCHEDULER_IMG=$1
V_LEVEL=$2

# scheduler-config.conf file arguments
SAFETY_MARGIN=$3
MINIMUM_RUNNING=$4
SCHEDULING_POLICY=$5

echo ""
echo ">>> Updating the scheduler-config.conf file"
echo ""
cat conf/scheduler/scheduler-config-base-$SCHEDULING_POLICY.conf | sed "s/SAFETY_MARGIN/$SAFETY_MARGIN/g" | sed "s/MINIMUM_RUNNING/$MINIMUM_RUNNING/g" > scheduler-config.conf

sudo cp scheduler-config.conf /etc/kubernetes/
rm scheduler-config.conf

echo ""
echo ">>> Configuring scheduler image"
echo ""
cat conf/scheduler/kube-scheduler-base.yaml | sed "s/SCHEDULER/bernardoezk\/$SCHEDULER_IMG/g" | sed "s/LEVEL/$V_LEVEL/g" > kube-scheduler.yaml

sudo cp kube-scheduler.yaml /etc/kubernetes/manifests/
rm kube-scheduler.yaml