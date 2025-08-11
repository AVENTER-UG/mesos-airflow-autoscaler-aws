package mesosaws

import (
	cfg "github.com/AVENTER-UG/mesos-autoscale/types"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
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

	// Basis-Session
	baseSess, err := session.NewSession()

	if err != nil {
		logrus.WithField("func", "MesosAWSNew").Error("Could not create base session: ", err.Error())
		return nil
	}

	// Temp Credentials with autoupdate
	roleCreds := stscreds.NewCredentials(baseSess, config.AWSAssumeRoleARN)

	e.AWSConfig.Credentials = roleCreds

	e.Session, err = session.NewSession()
	if err != nil {
		logrus.WithField("func", "MesosAWSNew").Error("Could not create session with assumed role: ", err.Error())
		return nil
	}

	// Try-Catch-Block
	errRe := recover()
	if errRe != nil {
		logrus.WithField("func", "MesosAWSNew").Error("Recovered error: ", errRe)
	}

	return e
}


func (e *AWS) FindMatchedInstanceType(mem int64, cpu int64, arch string) string {
	logrus.WithField("func", "aws.FindMarchedInstanceType").Debug()

	for _, p := range e.Config.AWSInstanceAllow {
		if int64((p.MEM)*1024) >= mem && int64(p.CPU) >= cpu {

			logrus.WithField("func", "aws.FindMatchedInstanceType").Trace("Found CPU: ", p.CPU)
			logrus.WithField("func", "aws.FindMatchedInstanceType").Trace("Found MEM: ", p.MEM)
			logrus.WithField("func", "aws.FindMatchedInstanceType").Trace("------------------------------------")

			return p.Name
		}

	}

	logrus.WithField("func", "aws.FindMatchedInstanceType").Warn("Could not found matching instance type")
	logrus.WithField("func", "aws.FindMatchedInstanceType").Warn("Need CPU: ", cpu)
	logrus.WithField("func", "aws.FindMatchedInstanceType").Warn("Need MEM: ", mem)
	logrus.WithField("func", "aws.FindMatchedInstanceType").Warn("Need Architecture: ", arch)
	return e.Config.AWSInstanceFallback
}
