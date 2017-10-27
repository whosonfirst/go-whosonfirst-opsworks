package main

import (
	"flag"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/opsworks"
	"log"
)

func main() {

	var stack = flag.String("stack", "", "A valid Opsworks stack")
	var layer = flag.String("layer", "", "A valid Opsworks layer")
	var region = flag.String("region", "us-east-1", "A valid AWS region")

	flag.Parse()

	cfg := aws.NewConfig()
	cfg.WithRegion(*region)

	sess := session.New(cfg)
	ops := opsworks.New(sess)

	params := opsworks.DescribeInstancesInput{
		LayerId: layer,
		StackId: stack,
	}

	desc, err := ops.DescribeInstances(&params)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(desc)
}
