#!/bin/bash

# ------------------------------------------------------------------------------
# MINIMUM RUNNING set to 300s
# ------------------------------------------------------------------------------

# ------------------------------------------------------------------------------
# average workload - medium contention - replication 4
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.15 3 30s 300s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-medium-contention-average-replication4-v2.csv 1800 1 1 >> log_experiment1_300s_2022_debug

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.15 2 30s 300s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-medium-contention-average-replication4-v2.csv 1800 1 1 >> log_experiment1_300s_2022_debug

sleep 60

