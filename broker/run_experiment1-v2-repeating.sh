#!/bin/bash

# ------------------------------------------------------------------------------
# MINIMUM RUNNING set to 30s
# ------------------------------------------------------------------------------

# ------------------------------------------------------------------------------
# concurrent workload - medium contention - replication 4
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.14 1 30s 30s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-medium-contention-concurrent-replication4-v2.csv 1800 1 1 >> log_experiment1_30s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.14 2 30s 30s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-medium-contention-concurrent-replication4-v2.csv 1800 1 1 >> log_experiment1_30s

sleep 60

# ------------------------------------------------------------------------------
# concurrent workload - high contention - replication 4
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.14 1 30s 30s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-high-contention-concurrent-replication4-v2.csv 1800 1 1 >> log_experiment1_30s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.14 2 30s 30s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-high-contention-concurrent-replication4-v2.csv 1800 1 1 >> log_experiment1_30s

sleep 60

# ------------------------------------------------------------------------------
# concurrent workload - medium contention - replication 2
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.14 1 30s 30s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-medium-contention-concurrent-replication2-v2.csv 1800 1 1 >> log_experiment1_30s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.14 2 30s 30s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-medium-contention-concurrent-replication2-v2.csv 1800 1 1 >> log_experiment1_30s

sleep 60

# ------------------------------------------------------------------------------
# concurrent workload - medium contention - replication 8
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.14 1 30s 30s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-medium-contention-concurrent-replication8-v2.csv 1800 1 1 >> log_experiment1_30s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.14 2 30s 30s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-medium-contention-concurrent-replication8-v2.csv 1800 1 1 >> log_experiment1_30s

sleep 60

# ------------------------------------------------------------------------------
# concurrent workload - high contention - replication 2
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.14 1 30s 30s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-high-contention-concurrent-replication2-v2.csv 1800 1 1 >> log_experiment1_30s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.14 2 30s 30s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-high-contention-concurrent-replication2-v2.csv 1800 1 1 >> log_experiment1_30s

sleep 60

# ------------------------------------------------------------------------------
# concurrent workload - high contention - replication 8
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.14 1 30s 30s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-high-contention-concurrent-replication8-v2.csv 1800 1 1 >> log_experiment1_30s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.14 2 30s 30s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-high-contention-concurrent-replication8-v2.csv 1800 1 1 >> log_experiment1_30s

sleep 60
