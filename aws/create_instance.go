package mesosaws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
)

func (e *AWS) CreateInstance(instanceType string) {
	// Create EC2 service client
	svc := ec2.New(e.Session)

	// Specify the details of the instance that you want to create.
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		// An Amazon Linux AMI ID for t2.micro instances in the us-west-2 region
		ImageId:      aws.String(e.Config.AWSImageID),
		InstanceType: aws.String(instanceType),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
	})

	if err != nil {
		logrus.WithField("func", "CreateInstance").Error("Could not create instance: ", err.Error())
		return
	}

	logrus.WithField("func", "CreateInstance").Info("Created Instance: ", *runResult.Instances[0].InstanceId)
}
