package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/opsworks"
	"log"
	"os"
)

func main() {

	var stack_id = flag.String("stack_id", "", "A valid Opsworks stack ID")
	var stack = flag.String("stack", "", "A valid Opsworks stack name")
	var layer_id = flag.String("layer_id", "", "A valid Opsworks layer ID")
	var layer = flag.String("layer", "", "A valid Opsworks layer name")
	var region = flag.String("region", "us-east-1", "A valid AWS region")
	var brief = flag.Bool("brief", false, "")

	flag.Parse()

	cfg := aws.NewConfig()
	cfg.WithRegion(*region)

	sess := session.New(cfg)
	ops := opsworks.New(sess)

	params := opsworks.DescribeInstancesInput{}

	if *stack_id != "" {
		params.StackId = stack_id
	} else if *layer_id != "" {
		params.LayerId = layer_id
	} else if *stack != "" {

		stacks_args := opsworks.DescribeStacksInput{}
		stacks_rsp, err := ops.DescribeStacks(&stacks_args)

		if err != nil {
			log.Fatal(err)
		}

		for _, st := range stacks_rsp.Stacks {

			if *st.Name == *stack {
				stack_id = st.StackId
				break
			}
		}

		if *stack_id == "" {
			err := errors.New("Invalid or unknown stack name")
			log.Fatal(err)
		}

		if *layer == "" {
			params.StackId = stack_id
		} else {

			layers_args := opsworks.DescribeLayersInput{
				StackId: stack_id,
			}

			layers_rsp, err := ops.DescribeLayers(&layers_args)

			if err != nil {
				log.Fatal(err)
			}

			for _, lyr := range layers_rsp.Layers {

				if *lyr.Name == *layer {
					layer_id = lyr.LayerId
					break
				}
			}

			if *layer_id == "" {
				err := errors.New("Invalid or unknown layer name")
				log.Fatal(err)
			}

			params.LayerId = layer_id

		}

	} else {
		log.Fatal("Insufficient parameters")
	}

	desc, err := ops.DescribeInstances(&params)

	if err != nil {
		log.Fatal(err)
	}

	if *brief {

		wr := csv.NewWriter(os.Stdout)

		for _, i := range desc.Instances {

			out := []string{*i.Hostname, *i.Status}

			if *i.Status == "online" {
				out = append(out, *i.PublicIp)
				out = append(out, *i.PrivateIp)
			}

			err := wr.Write(out)

			if err != nil {
				log.Fatal(err)
			}

		}

		wr.Flush()

		if wr.Error() != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}

	b, err := json.Marshal(desc)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}
