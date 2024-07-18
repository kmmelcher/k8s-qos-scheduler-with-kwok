#!/bin/bash

# ------------------------------------------------------------------------------
# MINIMUM RUNNING set to 30s
# ------------------------------------------------------------------------------

# ------------------------------------------------------------------------------
# aggregate workload - very high contention - replication 4 - 45 bot controllers
# ------------------------------------------------------------------------------

#cd ../cluster
#bash configure_cluster.sh qos-driven-scheduler:v0.2.5.15 3 30s 30s priority
#cd ../broker

#sleep 60

#bash run_experiment.sh workloads/workload-exp2-very-high-contention-45-bot-aggregate-replication4-v2.csv 1800 1 1 >> log_experiment2_new_workloads
#
#sleep 60

#cd ../cluster
#bash configure_cluster.sh qos-driven-scheduler:v0.2.5.15 2 30s 30s qos-driven
#cd ../broker

#sleep 60

#bash run_experiment.sh workloads/workload-exp2-very-high-contention-45-bot-aggregate-replication4-v2.csv 1800 1 1 >> log_experiment2_new_workloads

#sleep 60


# ------------------------------------------------------------------------------
# aggregate workload - very high contention - replication 4 - 50 bot controllers
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.15 1 30s 30s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp2-very-high-contention-55-bot-aggregate-replication4-v2.csv 1800 1 1 >> log_experiment2_new_workloads

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.15 2 30s 30s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp2-very-high-contention-55-bot-aggregate-replication4-v2.csv 1800 1 1 >> log_experiment2_new_workloads

sleep 60
