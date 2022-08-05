package mesosaws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
)

func (e *AWS) CreateInstance(instanceType string) *ec2.Reservation {
	// Create EC2 service client
	e.SVC = ec2.New(e.Session)

	// Specify the details of the instance that you want to create.
	runResult, err := e.SVC.RunInstances(&ec2.RunInstancesInput{
		// An Amazon Linux AMI ID for t2.micro instances in the us-west-2 region
		InstanceType: aws.String(instanceType),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		LaunchTemplate: &ec2.LaunchTemplateSpecification{
			LaunchTemplateId: aws.String(e.Config.AWSLaunchTemplateID),
		},
	})

	if err != nil {
		logrus.WithField("func", "CreateInstance").Error("Could not create instance: ", err.Error())
		return &ec2.Reservation{}
	}

	logrus.WithField("func", "CreateInstance").Info("Created Instance: ", *runResult.Instances[0].InstanceId)
	return runResult

}
