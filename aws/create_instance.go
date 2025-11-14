package mesosaws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
)

// CreateInstance - Create a AWS instance
func (e *AWS) CreateInstance(instanceType string, gpu bool) *ec2.Reservation {
	// Create EC2 service client
	e.SVC = ec2.New(e.Session)

	launchID := e.Config.AWSLaunchTemplateID
	if gpu && len(e.Config.AWSLaunchTemplateGPUID) > 0 {
		launchID = e.Config.AWSLaunchTemplateGPUID
	}

	// Specify the details of the instance that you want to create.
	runResult, err := e.SVC.RunInstances(&ec2.RunInstancesInput{
		// An Amazon Linux AMI ID for t2.micro instances in the us-west-2 region
		InstanceType: aws.String(instanceType),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		LaunchTemplate: &ec2.LaunchTemplateSpecification{
			LaunchTemplateId: aws.String(launchID),
		},
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Groups"),
						Value: aws.String("Worker"),
					},
					{
						Key:   aws.String("Inventory"),
						Value: aws.String("autoscaler"),
					},
					{
						Key:   aws.String("InstanceType"),
						Value: aws.String(instanceType),
					},
				},
			},
		},
	})

	if err != nil {
		// create fallback instance if the AWS capacity is not enough
		if strings.Contains(err.Error(), "InsufficientInstanceCapacity") {
			logrus.WithField("func", "mesosaws.CreateInstance").Info("Insufficient instance capacity. Try to create fallback instance.")
			return e.CreateInstance(e.Config.AWSInstanceFallbackGPU, gpu)
		}
		logrus.WithField("func", "mesosaws.CreateInstance").Error("Could not create instance: ", err.Error())
		return &ec2.Reservation{}
	}

	logrus.WithField("func", "mesosaws.CreateInstance").Info("Created Instance: ", *runResult.Instances[0].InstanceId)
	return runResult
}
