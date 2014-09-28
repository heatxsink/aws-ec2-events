package main

import (
	"github.com/heatxsink/goamz/aws"
	"github.com/heatxsink/goamz/ec2"
	"github.com/heatxsink/go-colour"
	"fmt"
	"os"
	"flag"
)

var (
	aws_access_key_id = ""
	aws_secret_access_key = ""
)

func init() {
	flag.StringVar(&aws_access_key_id, "key", "", "AWS Access Key Id.")
	flag.StringVar(&aws_secret_access_key, "secret", "", "AWS Secret Access Key")
	flag.Parse()
}

func main() {
	if aws_access_key_id != "" && aws_secret_access_key != "" {
		auth := aws.Auth{ SecretKey: aws_secret_access_key, AccessKey: aws_access_key_id }
		elastic_compute_cloud := ec2.New(auth, aws.USEast)
		response, err := elastic_compute_cloud.DescribeInstances(nil, nil)
		if err != nil {
			fmt.Println(err)
		}
		for _, reservation := range response.Reservations {
			for _, instance := range reservation.Instances {
				for _, tag := range instance.Tags {
					if tag.Key == "Name" {
						fmt.Println("Name: ", colour.Colourize(tag.Value, colour.FgYellow))
					}
				}
				//http://docs.aws.amazon.com/AWSEC2/latest/APIReference/ApiReference-query-DescribeInstanceStatus.html
				r, err := elastic_compute_cloud.DescribeInstanceStatus([]string{instance.InstanceId}, false, nil)
				if err != nil {
					fmt.Println(err)
				}
				for _, instance_set := range r.InstanceStatusSet {
					fmt.Println("\tAvailability Zone: ", instance_set.AvailabilityZone)
					fmt.Println("\tInstance State:    ", instance_set.InstanceState.Name)
					fmt.Println("\tSystem Status:     ", instance_set.SystemStatus.Name)
					fmt.Println("\tInstance Status:   ", instance_set.InstanceStatus.Name)
					if len(instance_set.EventsSet) > 0 {
						fmt.Println("\tEvents")
						fmt.Println("\t------")
						for _, event_set := range instance_set.EventsSet {
							fmt.Println("\t\tEvent Code:         ", colour.Colourize(event_set.Code, colour.FgWhite + colour.BgRed))
							fmt.Println("\t\tEvent Description : ", event_set.Description)
							fmt.Println("\t\tNot Before:         ", event_set.NotBefore)
							fmt.Println("\t\tNot After:          ", event_set.NotAfter)
						}
					}
				}
				fmt.Println()
			}
		}
	} else {
		fmt.Println("ERROR: Both key and secret need to be present in cli args.")
		os.Exit(1)
	}
}