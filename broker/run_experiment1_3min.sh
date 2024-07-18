#!/bin/bash

# ------------------------------------------------------------------------------
# MINIMUM RUNNING set to 180s
# ------------------------------------------------------------------------------

# ------------------------------------------------------------------------------
# average workload - medium contention - replication 4
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 1 30s 180s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-medium-contention-average-replication4.csv 3600 1 1 >> log_experiment1_180s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 2 30s 180s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-medium-contention-average-replication4.csv 3600 1 1 >> log_experiment1_180s

sleep 60

# ------------------------------------------------------------------------------
# concurrent workload - medium contention - replication 4
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 1 30s 180s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-medium-contention-concurrent-replication4.csv 3600 1 1 >> log_experiment1_180s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 2 30s 180s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-medium-contention-concurrent-replication4.csv 3600 1 1 >> log_experiment1_180s

sleep 60

# ------------------------------------------------------------------------------
# independent workload - medium contention - replication 4
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 1 30s 180s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-medium-contention-independent-replication4.csv 3600 1 1 >> log_experiment1_180s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 2 30s 180s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-medium-contention-independent-replication4.csv 3600 1 1 >> log_experiment1_180s

sleep 60


# ------------------------------------------------------------------------------
# average workload - high contention - replication 4
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 1 30s 180s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-high-contention-average-replication4.csv 3600 1 1 >> log_experiment1_180s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 2 30s 180s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-high-contention-average-replication4.csv 3600 1 1 >> log_experiment1_180s

sleep 60

# ------------------------------------------------------------------------------
# concurrent workload - high contention - replication 4
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 1 30s 180s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-high-contention-concurrent-replication4.csv 3600 1 1 >> log_experiment1_180s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 2 30s 180s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-high-contention-concurrent-replication4.csv 3600 1 1 >> log_experiment1_180s

sleep 60

# ------------------------------------------------------------------------------
# independent workload - high contention - replication 4
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 1 30s 180s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-high-contention-independent-replication4.csv 3600 1 1 >> log_experiment1_180s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 2 30s 180s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-high-contention-independent-replication4.csv 3600 1 1 >> log_experiment1_180s

sleep 60

# ------------------------------------------------------------------------------
# average workload - no contention - replication 4
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 1 30s 180s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-no-contention-average-replication4.csv 3600 1 1 >> log_experiment1_180s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 2 30s 180s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-no-contention-average-replication4.csv 3600 1 1 >> log_experiment1_180s

sleep 60

# ------------------------------------------------------------------------------
# concurrent workload - no contention - replication 4
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 1 30s 180s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-no-contention-concurrent-replication4.csv 3600 1 1 >> log_experiment1_180s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 2 30s 180s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-no-contention-concurrent-replication4.csv 3600 1 1 >> log_experiment1_180s

sleep 60

# ------------------------------------------------------------------------------
# independent workload - no contention - replication 4
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 1 30s 180s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-no-contention-independent-replication4.csv 3600 1 1 >> log_experiment1_180s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 2 30s 180s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-no-contention-independent-replication4.csv 3600 1 1 >> log_experiment1_180s

sleep 60

# ------------------------------------------------------------------------------
# concurrent workload - medium contention - replication 2
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 1 30s 180s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-medium-contention-concurrent-replication2.csv 3600 1 1 >> log_experiment1_180s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 2 30s 180s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-medium-contention-concurrent-replication2.csv 3600 1 1 >> log_experiment1_180s

sleep 60

# ------------------------------------------------------------------------------
# concurrent workload - medium contention - replication 8
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 1 30s 180s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-medium-contention-concurrent-replication8.csv 3600 1 1 >> log_experiment1_180s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 2 30s 180s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-medium-contention-concurrent-replication8.csv 3600 1 1 >> log_experiment1_180s

sleep 60

# ------------------------------------------------------------------------------
# concurrent workload - high contention - replication 2
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 1 30s 180s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-high-contention-concurrent-replication2.csv 3600 1 1 >> log_experiment1_180s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 2 30s 180s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-high-contention-concurrent-replication2.csv 3600 1 1 >> log_experiment1_180s

sleep 60

# ------------------------------------------------------------------------------
# concurrent workload - high contention - replication 8
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 1 30s 180s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-high-contention-concurrent-replication8.csv 3600 1 1 >> log_experiment1_180s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.12 2 30s 180s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp1-high-contention-concurrent-replication8.csv 3600 1 1 >> log_experiment1_180s

sleep 60
