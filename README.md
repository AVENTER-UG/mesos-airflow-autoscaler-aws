# AWS Autoscaler for Apache Airflow Mesos Provider

[![Issues](https://img.shields.io/static/v1?label=&message=Issues&color=brightgreen)](https://github.com/m3scluster/mesos-airflow-autoscaler-aws/issues)
[![Chat](https://img.shields.io/static/v1?label=&message=Chat&color=brightgreen)](https://matrix.to/#/#mesos:matrix.aventer.biz?via=matrix.aventer.biz)
[![GoDoc](https://godoc.org/github.com/AVENTER-UG/mesos-dns?status.svg)](https://godoc.org/github.com/AVENTER-UG/mesos-airflow-autoscaler-aws)
[![Docker Pulls](https://img.shields.io/docker/pulls/avhost/mesos-airflow-autoscaler-aws)](https://hub.docker.com/repository/docker/avhost/mesos-airflow-autoscaler-aws/)

This AWS autoscaler will start EC2 instances if Airflow does not get matched offers from mesos.

## Issues

To open an issue, please use this place: https://github.com/m3scluster/mesos-airflow-autoscaler-aws/issues

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
| WAIT_TIME_OVERWRITE_INSTANCE | 15m | If the DAG is still after 15m in the queue, the autoscaler will ignore the custom instance type to start a new ec2 instance. |
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
| AWS_INSTANCE_FALLBACK | t3a.2xlarge | Fallback instance type will be used if there are no more ec2 resources in AWS. |
| AWS_INSTANCE_DEFAULT_ARCHITECTURE | x86_64 | Default architecture of ec2 instance. | 
| AWS_INSTANCE_ALLOW | All instances are allowed by default | Only these instances are allowed. Format should be: `[{ 'Name': 't3.large', 'CPU': 2.0, 'MEM': 8.0 }]` |
| AWS_INSTANCE_MAX_TRY | 3 | The max tries to create a instance |
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
