#!/bin/bash

# ------------------------------------------------------------------------------
# aggregate workload - very high contention - replication 4 - 300s
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.15 1 30s 300s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp2-very-high-contention-aggregate-replication4-v2.csv 1800 1 1 >> log_experiment2_300s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.15 2 30s 300s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp2-very-high-contention-aggregate-replication4-v2.csv 1800 1 1 >> log_experiment2_300s

sleep 60


# ------------------------------------------------------------------------------
# aggregate workload - medium contention - replication 4 - 30s
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.15 1 30s 30s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp2-medium-contention-aggregate-replication4-v2.csv 1800 1 1 >> log_experiment2_30s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.15 2 30s 30s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp2-medium-contention-aggregate-replication4-v2.csv 1800 1 1 >> log_experiment2_30s

sleep 60

# ------------------------------------------------------------------------------
# aggregate workload - higher contention - replication 4 - 30s
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.15 3 30s 30s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp2-higher-contention-aggregate-replication4-v2.csv 1800 1 1 >> log_experiment2_30s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.15 2 30s 30s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp2-higher-contention-aggregate-replication4-v2.csv 1800 1 1 >> log_experiment2_30s

sleep 60
