package mesosaws

import (
	cfg "github.com/AVENTER-UG/mesos-autoscale/types"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
)

func (e *AWS) TerminateInstance(instance cfg.DagTask) {
	// Create EC2 service client
	e.SVC = ec2.New(e.Session)

	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{instance.EC2.Instances[0].InstanceId},
	}

	// terminate the given instance
	runResult, err := e.SVC.TerminateInstances(input)

	if err != nil {
		logrus.WithField("func", "TerminateInstance").Error("Could not terminate instance: ", err.Error())
		return
	}

	logrus.WithField("func", "TerminateInstance").Info(*runResult)
}
