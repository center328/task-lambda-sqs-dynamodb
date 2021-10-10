package stream

import "os"

// Config Wrapper for Config methods
var (
	AWSRegion 	= "ap-southeast-2"
	AWSKey		= ""
	AWSSecret 	= ""

	// Poll from this SQS URL
	URL 		= ""
)

func init() {

	awsRegion := os.Getenv("AWS_REGION")
	if len(awsRegion) > 0 {
		AWSRegion = awsRegion
	}

	awsAccessKey := os.Getenv("AWS_ACCESS_KEY")
	if len(awsAccessKey) > 0 {
		AWSKey = awsAccessKey
	}

	awsSecretKey := os.Getenv("AWS_SECRET_KEY")
	if len(awsSecretKey) > 0 {
		AWSSecret = awsSecretKey
	}

	qURL := os.Getenv("QUEUE_URL")
	if len(qURL) > 0 {
		URL = qURL
	}

}