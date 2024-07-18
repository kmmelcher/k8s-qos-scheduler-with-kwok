#!/bin/bash

# ------------------------------------------------------------------------------
# MINIMUM RUNNING set to 30s
# ------------------------------------------------------------------------------

# ------------------------------------------------------------------------------
# validting kubernetes new version using custom scheduler
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.18 1 30s 30s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-validation-k8s-version-2workers.csv 1800 1 1 >> log_validation2023

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.18 2 30s 30s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-validation-k8s-version-2workers.csv 1800 1 1 >> log_validation2023
