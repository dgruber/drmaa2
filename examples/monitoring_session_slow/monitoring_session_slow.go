package main

import (
	"fmt"
	"github.com/dgruber/drmaa2"
	"time"
)

func main() {
	// A Session Manager is always required
	var sm drmaa2.SessionManager

	// Open Monitoring Session and check for an error
	ms, err := sm.OpenMonitoringSession("")
	if err != nil {
		fmt.Printf("Failed when creating monitoring session: %s\n", err)
	}
	// We need to close the Monitoring Session at exit
	defer ms.CloseMonitoringSession()

	d, _ := time.ParseDuration("1s")
	for {
		if jobs, err := ms.GetAllJobs(nil); err != nil {
			fmt.Printf("Error during GetAllJobs() call: %s\n", err)
		} else {
			fmt.Printf("All jobs : %s\n", jobs)
			for _, j := range jobs {
				fmt.Println(j.GetState())
				ji, _ := j.GetJobInfo()
				fmt.Printf("job info: %+v\n", ji)
			}
		}
		time.Sleep(d)
	}
}
