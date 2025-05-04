package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSQueue struct {
	Name     string
	FullName string
}

type SQSConfig struct {
	WaitTime   int64
	BaseURL    string
	Queues     []SQSQueue
	Session    *session.Session
	Client     *sqs.SQS
	NumWorkers int
}

type AWSConfig struct {
	AccessKey string
	SecretKey string
	Region    string
	SQSConfig SQSConfig
}

const defaultSQSWaitTime = int64(5)
const defaultSQSNumWorkers = 5

func LoadAWSConfig(projectCfg *Config) *AWSConfig {
	return &AWSConfig{
		AccessKey: fetchEnv("AWS_ACCESS_KEY_ID"),
		SecretKey: fetchEnv("AWS_SECRET_ACCESS_KEY"),
		Region:    getEnv("AWS_REGION", "us-east-1"),
		SQSConfig: SQSConfig{
			WaitTime:   getEnvAsInt64("AWS_SQS_WAIT_TIME", defaultSQSWaitTime),
			BaseURL:    getEnv("AWS_BASE_URL", ""),
			Queues:     projectCfg.loadQueues(),
			NumWorkers: getEnvAsInt("AWS_SQS_NUM_WORKERS", defaultSQSNumWorkers),
		},
	}
}

func (appConfig *Config) loadQueues() []SQSQueue {
	loadedSQSQueues := make([]SQSQueue, 0, len(sqsQueueNames()))

	for _, queueName := range sqsQueueNames() {
		sqsQueue := SQSQueue{
			Name:     queueName,
			FullName: appConfig.AppName + "_" + queueName,
		}

		loadedSQSQueues = append(loadedSQSQueues, sqsQueue)
	}

	return loadedSQSQueues
}

func InstantiateSQSClient(awsConfig *AWSConfig) {
	awsSession := session.Must(
		session.NewSession(&aws.Config{
			Endpoint: aws.String(awsConfig.SQSConfig.BaseURL),
			Region:   aws.String(awsConfig.Region)},
		),
	)

	awsConfig.SQSConfig.Session = awsSession
	awsConfig.SQSConfig.Client = sqs.New(awsConfig.SQSConfig.Session)
}

func sqsQueueNames() []string {
	return []string{"payment_events"}
}
