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

	// The Univa Grid Engine DRMAA2 implementation is event based
	// meaning that GetAllJobs() can be called as often as required,
	// it does not trigger addional communication to the qmaster
	// process.
	if jobs, err := ms.GetAllJobs(nil); err != nil {
		fmt.Printf("Error during GetAllJobs() call: %s\n", err)
	} else {
		for _, j := range jobs {
			t1 := time.Now()
			if i, err := j.GetJobInfo(); err != nil {
				fmt.Printf("Error during GetJobInfo() call: %s\n", err)
			} else {
				fmt.Printf("Job id [%s] Job Name [%s].  Call took [%s]\n", i.Id, i.JobName, time.Since(t1).String())
			}
		}

	}
}
