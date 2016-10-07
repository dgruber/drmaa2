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

	// Create a JobInfo as filter with UNSET values. This is required
	// otherwise the filter will not work correctly.
	// ji := drmaa2.CreateJobInfo()
	// We want to list all jobs submitted by that user.
	// ji.JobOwner = "root"

	// The Univa Grid Engine DRMAA2 implementation is event based
	// meaning that GetAllJobs() can be called as often as required,
	// it does not trigger addional communication to the qmaster
	// process.
	d, _ := time.ParseDuration("1s")
	for {
		if jobs, err := ms.GetAllJobs(nil); err != nil {
			fmt.Printf("Error during GetAllJobs() call: %s\n", err)
		} else {
			fmt.Printf("All jobs : %s\n", jobs)
			for _, j := range jobs {
				fmt.Printf("Job ID: %s\n", j.GetId())
				fmt.Println(j.GetState())
				if ji, errJI := j.GetJobInfo(); errJI == nil {
					fmt.Printf("job info extensions: %s\n", ji.ListExtensions())
				}
			}
		}
		break
		time.Sleep(d)
	}
}
