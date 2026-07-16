package email

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type EmailService struct {
	client *ses.Client
	from   string
}

func NewEmailService(
	region string,
	accessKey string,
	secretKey string,
	from string,
) (*EmailService, error) {

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				accessKey,
				secretKey,
				"",
			),
		),
	)

	if err != nil {
		return nil, err
	}

	client := ses.NewFromConfig(cfg)

	return &EmailService{
		client: client,
		from:   from,
	}, nil
}

func (e *EmailService) Send(
	to string,
	subject string,
	body string,
) error {

	_, err := e.client.SendEmail(
		context.Background(),
		&ses.SendEmailInput{
			Source: aws.String(e.from),

			Destination: &types.Destination{
				ToAddresses: []string{to},
			},

			Message: &types.Message{
				Subject: &types.Content{
					Data: aws.String(subject),
				},

				Body: &types.Body{
					Html: &types.Content{
						Data: aws.String(body),
					},
				},
			},
		},
	)

	return err
}
