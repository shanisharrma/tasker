package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awsConfig "github.com/shanisharrma/tasker/internal/shared/config"
)

type AWS struct {
	S3 *S3Client
}

func NewAWS(awsConfig *awsConfig.Config) (*AWS, error) {

	configOptions := []func(*config.LoadOptions) error{
		config.WithRegion(awsConfig.AWS.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			awsConfig.AWS.AccessKeyID,
			awsConfig.AWS.SecretAccessKey,
			"",
		)),
	}

	// Add custom endpoint if provided (for S3-compatible services like Sevalla or cloudflare R2)
	if awsConfig.AWS.EndpointURL != "" {
		configOptions = append(configOptions, config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string,
				options ...interface{},
			) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           awsConfig.AWS.EndpointURL,
					SigningRegion: awsConfig.AWS.Region,
				}, nil
			}),
		))
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), configOptions...)
	if err != nil {
		return nil, err
	}

	return &AWS{
		S3: NewS3Client(cfg),
	}, nil
}
