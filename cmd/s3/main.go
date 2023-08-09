package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	AWS_ACCESS_KEY_ID     = "AKIA2WX$NLLSMS4MCG7D"
	AWS_SECRET_ACCESS_KEY = "H2/jMZGYYbImP8oZljUfH0ClkJkLopXjkt2cePwe"
)

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-north-1"),
		Credentials: credentials.NewStaticCredentials(AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, ""),
	})

	if err != nil {
		fmt.Println("Error creating session:", err)
		return
	}

	s3Svc := s3.New(sess)
	result, err := s3Svc.ListBuckets(nil)
	if err != nil {
		fmt.Println("Error listing buckets:", err)
		return
	}

	fmt.Println("Buckets:")
	for _, bucket := range result.Buckets {
		fmt.Println(*bucket.Name)
	}
}
