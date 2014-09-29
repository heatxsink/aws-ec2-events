package main

import (
	"github.com/heatxsink/goamz/aws"
	"github.com/heatxsink/goamz/ec2"
	"github.com/heatxsink/go-colour"
	"github.com/heatxsink/go-yeam"
	"github.com/heatxsink/go-simpleconfig"
	"fmt"
	"os"
	"flag"
	"time"
)

var (
	aws_access_key_id = ""
	aws_secret_access_key = ""
	default_cache_path = os.Getenv("HOME")
	cache_path = ""
	cache_file_path = ""
	alert_flag = false
	alert_email_address = ""
	imap_username = ""
	imap_password = ""
)

func init() {
	flag.StringVar(&aws_access_key_id, "key", "", "AWS Access Key Id.")
	flag.StringVar(&aws_secret_access_key, "secret", "", "AWS Secret Access Key")
	flag.StringVar(&cache_path, "cache_path", default_cache_path, "Change where you would like the file backed cache to live.")
	flag.BoolVar(&alert_flag, "alert", true, "Enable/disable alerts via email.")
	flag.StringVar(&alert_email_address, "alert_email", "", "Email addres to send alerts when there's been an event.")
	flag.StringVar(&imap_username, "imap_username", "", "imap username")
	flag.StringVar(&imap_password, "imap_password", "", "imap password")

	flag.Parse()
	cache_file_path = fmt.Sprintf("%s/.aws-ec2-events-cache", cache_path)
}

type Ec2Event struct {
	Code string `json:"code"`
	Description string `json:"description"`
	NotBefore time.Time `json:"not-before"`
	NotAfter time.Time `json:"not-after"`
}

type Ec2Instance struct {
	Name string `json:"name"`
	Id string `json:"id"`
	Events []Ec2Event `json:"events"`
}

func main() {
	if aws_access_key_id != "" && aws_secret_access_key != "" {
		auth := aws.Auth{ SecretKey: aws_secret_access_key, AccessKey: aws_access_key_id }
		elastic_compute_cloud := ec2.New(auth, aws.USEast)
		response, err := elastic_compute_cloud.DescribeInstances(nil, nil)
		if err != nil {
			fmt.Println(err)
		}
		var current_instances []Ec2Instance
		for _, reservation := range response.Reservations {
			var ec2_instance Ec2Instance
			for _, instance := range reservation.Instances {
				ec2_instance.Id = instance.InstanceId
				for _, tag := range instance.Tags {
					if tag.Key == "Name" {
						fmt.Println("Name: ", colour.Colourize(tag.Value, colour.FgYellow))
						ec2_instance.Name = tag.Value
					}
				}
				//
				//http://docs.aws.amazon.com/AWSEC2/latest/APIReference/ApiReference-query-DescribeInstanceStatus.html
				//
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
						fmt.Println("\t\tEvents")
						fmt.Println("\t\t------")
						for _, event_set := range instance_set.EventsSet {
							fmt.Println("\t\tEvent Code:         ", colour.Colourize(event_set.Code, colour.FgWhite + colour.BgRed))
							fmt.Println("\t\tEvent Description : ", event_set.Description)
							fmt.Println("\t\tNot Before:         ", event_set.NotBefore)
							fmt.Println("\t\tNot After:          ", event_set.NotAfter)
							var ec2_event Ec2Event
							ec2_event.Code = event_set.Code
							ec2_event.Description = event_set.Description
							ec2_event.NotBefore = event_set.NotBefore
							ec2_event.NotAfter = event_set.NotAfter
							ec2_instance.Events = append(ec2_instance.Events, ec2_event)
						}
					}
				}
				current_instances = append(current_instances, ec2_instance)
				fmt.Println()
			}
		}
		var last_instances []Ec2Instance
		err = simpleconfig.Load(cache_file_path, &last_instances)
		if err != nil {
			// File does not exist
			simpleconfig.Save(cache_file_path, current_instances)
			os.Exit(0)
		}
		if alert_flag {
			for _, l := range last_instances {
				for _, c := range current_instances {
					if l.Id == c.Id {
						if len(l.Events) != len(c.Events) {
							// save current_instances
							simpleconfig.Save(cache_file_path, current_instances)
							// send ghetto alert
							if imap_username != "" && imap_password != "" && alert_email_address != "" {
								e := yaem.New(imap_username, imap_password, yaem.GMAIL_SMTP_HOSTNAME, yaem.GMAIL_SMTP_PORT)
								reciepents := []string{ alert_email_address }
								err := e.SendEmail(reciepents, "New EC2 Event", fmt.Sprintf("The instance %s has new events.", l.Id))
								if err != nil {
									fmt.Println("ERROR: ", err)
								}
							}
						}
					}
				}
			}
		}
	} else {
		fmt.Println("ERROR: Both key and secret need to be present in cli args.")
		flag.PrintDefaults()
		os.Exit(1)
	}
}