#!/bin/bash

# ------------------------------------------------------------------------------
# MINIMUM RUNNING set to 30s
# ------------------------------------------------------------------------------

# ------------------------------------------------------------------------------
# average workload
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.10 1 30s 30s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/average-sanity-test-workload.csv 1800 1 1 >> log_sanity_workload_v0.2.5.10_30

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.10 2 30s 30s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/average-sanity-test-workload.csv 1800 1 1 >> log_sanity_workload_v0.2.5.10_30

sleep 60

# ------------------------------------------------------------------------------
# concurrent workload
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.10 1 30s 30s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/concurrent-sanity-test-workload.csv 1800 1 1 >> log_sanity_workload_v0.2.5.10_30

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.10 2 30s 30s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/concurrent-sanity-test-workload.csv 1800 1 1 >> log_sanity_workload_v0.2.5.10_30

sleep 60

# ------------------------------------------------------------------------------
# independent workload
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.10 1 30s 30s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/independent-sanity-test-workload.csv 1800 1 1 >> log_sanity_workload_v0.2.5.10_30

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.10 2 30s 30s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/independent-sanity-test-workload.csv 1800 1 1 >> log_sanity_workload_v0.2.5.10_30

sleep 60

# ------------------------------------------------------------------------------
# aggregate workload
# ------------------------------------------------------------------------------

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.10 1 30s 30s priority
cd ../broker

sleep 60

bash run_experiment.sh workloads/aggregate-sanity-test-workload.csv 1800 1 1 >> log_sanity_workload_v0.2.5.10_30

sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.10 2 30s 30s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/aggregate-sanity-test-workload.csv 1800 1 1 >> log_sanity_workload_v0.2.5.10_30

sleep 60

# ------------------------------------------------------------------------------
# MINIMUM RUNNING set to 180s
# ------------------------------------------------------------------------------

# ------------------------------------------------------------------------------
# average workload
# ------------------------------------------------------------------------------

#cd ../cluster
#bash configure_cluster.sh qos-driven-scheduler:v0.2.5.10 1 30s 180s priority
#cd ../broker

#sleep 60

#bash run_experiment.sh workloads/average-sanity-test-workload.csv 3600 1 1 >> log_sanity_workload_v0.2.5.10_180

#sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.10 3 30s 180s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/average-sanity-test-workload.csv 3600 1 1 >> log_sanity_workload_v0.2.5.10_180

sleep 60

# ------------------------------------------------------------------------------
# concurrent workload
# ------------------------------------------------------------------------------

#cd ../cluster
#bash configure_cluster.sh qos-driven-scheduler:v0.2.5.10 1 30s 180s priority
#cd ../broker

#sleep 60

#bash run_experiment.sh workloads/concurrent-sanity-test-workload.csv 3600 1 1 >> log_sanity_workload_v0.2.5.10_180

#sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.10 3 30s 180s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/concurrent-sanity-test-workload.csv 3600 1 1 >> log_sanity_workload_v0.2.5.10_180

sleep 60

# ------------------------------------------------------------------------------
# independent workload
# ------------------------------------------------------------------------------

#cd ../cluster
#bash configure_cluster.sh qos-driven-scheduler:v0.2.5.10 1 30s 180s priority
#cd ../broker

#sleep 60

#bash run_experiment.sh workloads/independent-sanity-test-workload.csv 3600 1 1 >> log_sanity_workload_v0.2.5.10_180

#sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.10 3 30s 180s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/independent-sanity-test-workload.csv 3600 1 1 >> log_sanity_workload_v0.2.5.10_180

sleep 60

# ------------------------------------------------------------------------------
# aggregate workload
# ------------------------------------------------------------------------------

#cd ../cluster
#bash configure_cluster.sh qos-driven-scheduler:v0.2.5.10 1 30s 180s priority
#cd ../broker

#sleep 60

#bash run_experiment.sh workloads/aggregate-sanity-test-workload.csv 3600 1 1 >> log_sanity_workload_v0.2.5.10_180

#sleep 60

cd ../cluster
bash configure_cluster.sh qos-driven-scheduler:v0.2.5.10 3 30s 180s qos-driven
cd ../broker

sleep 60

bash run_experiment.sh workloads/aggregate-sanity-test-workload.csv 3600 1 1 >> log_sanity_workload_v0.2.5.10_180

sleep 60
