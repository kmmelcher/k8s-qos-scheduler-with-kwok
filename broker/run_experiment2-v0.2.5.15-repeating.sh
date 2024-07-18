#!/bin/bash

# ------------------------------------------------------------------------------
# MINIMUM RUNNING set to 30s
# ------------------------------------------------------------------------------

# ------------------------------------------------------------------------------
# aggregate workload - medium contention - replication 4
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.15 1 30s 30s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp2-medium-contention-aggregate-replication4.csv 1800 1 1 >> log_experiment2_30s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.15 2 30s 30s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp2-medium-contention-aggregate-replication4.csv 1800 1 1 >> log_experiment2_30s

sleep 60