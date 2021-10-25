package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"sync"
)

var (
	AwsConf        *aws.Config
	AwsSession     *session.Session
	onceAwsSession sync.Once
)

func init() {
	onceAwsSession.Do(func() {
		AwsConf = &aws.Config{
			Region: aws.String(Env().AWSRegion),
		}
		session, err := session.NewSession(AwsConf)
		if err != nil {
			log.Fatalln("Error creating AWS session:", err)
		}
		AwsSession = session
	})
}

