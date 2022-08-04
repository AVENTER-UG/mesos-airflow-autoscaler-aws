package mesosaws

import (
	cfg "github.com/AVENTER-UG/mesos-autoscale/types"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// AWS struct about the AWS functions
type AWS struct {
	Config    *cfg.Config
	AWSConfig *aws.Config
	Session   *session.Session
	SVC       *ec2.EC2
}

func New(config *cfg.Config) *AWS {
	e := &AWS{
		Config: config,
		AWSConfig: &aws.Config{
			Region: aws.String(config.AWSRegion),
		},
	}
	var err error
	e.Session, err = session.NewSession(e.AWSConfig)

	if err != nil {
		logrus.WithField("func", "MesosAWSNew").Error("Could not create session: ", err.Error())
	}

	return e
}
