# AWS Autoscaler for Apache Airflow Mesos Provider

This AWS autoscaler will start EC2 instances if Airflow does not get matched offers from mesos.

## Requirements

- Airflow min 2.1.x
- Apache Mesos min 1.6.x
- Airflow Mesos Provider min 2.0
- AWS


## Configuration

| ENV | Default | Description |
| --- | --- | --- |
| AIRFLOW_MESOS_SCHEDULER | 127.0.0.1:11000 | IP Address and port of the Apache Airflow Mesos provider |
| LOGLEVEL | debug | Loglevel (info, warn, debug) |
| WAIT_TIME | 2m | The time in minutes the autoscaler have to wait until it will create a mesos instance in AWS |
| REDIS_SERVER | 127.0.0.1:6480 | Redis server and port |
| REDIS_PASSWORD | | Redis DB password |
| REDIS_DB | 2 | Redis DB Number |
| REDIS_PREFIX | asg | Prefix for every Redis key |
| API_USERNAME | user | Username to authenticate against the Apache Airflow Mesos provider |
| API_PASSWORD | password | Password to authenticate against the Apache Airflow Mesos provider |
| AWS_SECRET | | AWS Secret |
| AWS_REGION | eu-central-1 | AWS Region |
| AWS_WAIT_TIME | 10m | The time the autoscaler have to wait until it check if the EC2 intance can be terminated. |
| AWS_LAUNCH_TEMPLATE_ID | | The AWS Launche Template ID | 
| AWS_INSTANCE_16GB | t2.xlarge | The scale out instance type for memory limit lt 16GB (also the default instance if no memory limit was set) |
| AWS_INSTANCE_32GB | t3.2xlarge | The scale out instance type for memory limit gt 32GB |
| AWS_INSTANCE_64GB | r5.2xlarge | The scale out instance type for memory limit gt 64GB |
| MESOS_AGENT_SSL | false | Enable SSL for the communication to the Mesos agent |
| MESOS_AGENT_USERNAME | | Username of the Mesos Agent |
| MESOS_AGENT_PASSWORD | | Password of the Mesos Agent |
| MESOS_AGENT_PORT | 5051 | Port of the Mesos Agent |
| MESOS_AGENT_TIMEOUT | 10m | Mesos agent timeout |
| SSL | false | Enable SSL for the communication to the Airflow Scheduler API |
| DAG_TTL | 6h | Set the TTL for DAG'keys in Redis. | 

## Airflow Executor Config

These autoscaler support the following executor config parameters:

```bash
  executor_config={
    "mem_limit": 2048,
    "instance_type": "t2.xlarge"
 }
```

- With the help if the `mem_limit` parameter, the autoscaler will determin which AWS instance is the right one. 
- `instance_type` will overwrite the mem_limit and just crate a instance of that type.

## Documentation

AWS - https://docs.aws.amazon.com/code-samples/latest/catalog/go-ec2-create_instance.go.html
