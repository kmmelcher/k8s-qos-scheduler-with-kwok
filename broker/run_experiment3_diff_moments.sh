#!/bin/bash

# ------------------------------------------------------------------------------
# MINIMUM RUNNING set to 30s
# ------------------------------------------------------------------------------

# ------------------------------------------------------------------------------
# exp3 - moment 2
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.15 1 30s 30s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp3-moment2.csv 3600 1 1 >> log_experiment3_30s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.15 2 30s 30s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp3-moment2.csv 3600 1 1 >> log_experiment3_30s

sleep 60

# ------------------------------------------------------------------------------
# exp3 - moment 3
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.15 1 30s 30s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp3-moment3.csv 3900 1 1 >> log_experiment3_30s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.15 2 30s 30s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp3-moment3.csv 3900 1 1 >> log_experiment3_30s

sleep 60

# ------------------------------------------------------------------------------
# exp3 - moment 4
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.15 1 30s 30s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp3-moment4.csv 4500 1 1 >> log_experiment3_30s

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.15 2 30s 30s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/workload-exp3-moment4.csv 4500 1 1 >> log_experiment3_30s

sleep 60
