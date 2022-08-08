# AWS Autoscaler for Apache Airflow Mesos Provider

This AWS autoscaler will start EC2 instances if Airflow does not get matched offers from mesos.

## Configuration

| ENV | Default | Description |
| --- | --- | --- |
| AIRFLOW_MESOS_SCHEDULER | 127.0.0.1:11000 | IP Address and port of the Apache Airflow Mesos provider |
| LOGLEVEL | debug | Loglevel (info, warn, debug) |
| WAIT_TIME | 2m | The time in minutes the autoscaler have to wait until it will create a mesos instance in AWS |
| REDIS_SERVER | 127.0.0.1:6480 | Redis server and port |
| REDIS_PASSWORD | | Redis DB password |
| REDIS_DB | 2 | Redis DB |
| REDIS_PREFIX | asg | Prefix for every Redis key |
| API_USERNAME | user | Username to authenticate against the Apache Airflow Mesos provider |
| API_PASSWORD | password | Password to authenticate against the Apache Airflow Mesos provider |
| AWS_SECRET | | AWS Secret |
| AWS_REGION | eu-central-1 | AWS Region |
| AWS_WAIT_TIME | 5m | The time the autoscaler have to wait until it check if the EC2 intance can be terminated. |
| AWS_LAUNCH_TEMPLATE_ID | | The AWS Launche Template ID | 
| AWS_INSTANCE_16GB | t2.xlarge | The scale out instance type for memory limit lt 16GB (also the default instance if no memory limit was set) |
| AWS_INSTANCE_32GB | t3.2xlarge | The scale out instance type for memory limit gt 32GB |
| AWS_INSTANCE_64GB | r5.2xlarge | The scale out instance type for memory limit gt 64GB |


## Documentation

AWS - https://docs.aws.amazon.com/code-samples/latest/catalog/go-ec2-create_instance.go.html
