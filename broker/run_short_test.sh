#!/bin/bash

# ------------------------------------------------------------------------------
# MINIMUM RUNNING set to 30s
# ------------------------------------------------------------------------------

# ------------------------------------------------------------------------------
# average workload
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 1 30s 30s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/independent-sanity-test-workload.csv 600 1 1 >> log_sanity_workload_v0.2.5.11_30

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 3 30s 30s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/independent-sanity-test-workload.csv 600 1 1 >> log_sanity_workload_v0.2.5.11_30

sleep 60

