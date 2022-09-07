package mesosaws

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
)

// TerminateInstance - Terminate a running instance
func (e *AWS) TerminateInstance(instance *string) {
	// Create EC2 service client
	e.SVC = ec2.New(e.Session)

	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{instance},
	}

	// terminate the given instance
	runResult, err := e.SVC.TerminateInstances(input)

	if err != nil {
		logrus.WithField("func", "mesosaws.TerminateInstance").Error("Could not terminate instance: ", err.Error())
		return
	}

	logrus.WithField("func", "mesosaws.TerminateInstance").Infof("Terminate instance %s State: %s", *runResult.TerminatingInstances[0].InstanceId, *runResult.TerminatingInstances[0].CurrentState)
}
