package drmaa2_test

import (
	"fmt"
	"github.com/dgruber/drmaa2"
	"testing"
)

// Tests if MonitoringSession can be opened and closed.
// Requires the libdrmaa2.so in $LD_LIBRARY_PATH.
func TestOpenMonitoringSession(t *testing.T) {
	// Simple test for open and closing as MonitoringSession
	var sm drmaa2.SessionManager
	ms, err := sm.OpenMonitoringSession("")
	if err != nil {
		t.Errorf("Couldn't open Monitoring session. %s", err)
		if ms != nil {
			t.Errorf("MonitoringSession needs to be nil in case of error")
		}
		return
	}
	t.Log("OpenMonitoringSession() created a MonitoringSession succesfully")
	if err := ms.CloseMonitoringSession(); err != nil {
		t.Errorf("CloseMonitoringSession() returned error: %s", err)
	}
}

func TestMonitoringSessionGetAllMachines(t *testing.T) {
	var sm drmaa2.SessionManager
	ms, err := sm.OpenMonitoringSession("")
	if err != nil {
		t.Errorf("Couldn't open Monitoring session: %s", err)
		if ms != nil {
			t.Errorf("MonitoringSession needs to be nil in case of error")
		}
		return
	}
	// get all machines
	machine, err := ms.GetAllMachines(nil)

	if err != nil {
		t.Errorf("Error during GetAllMachines(nil): %s", err)
		return
	}
	amount := len(machine)
	if amount < 1 {
		t.Errorf("Error: No machine returned in GetAllMachines(nil)")
	}
	// get a single machine
	var names []string
	names = append(names, machine[0].Name)
	if machine2, err := ms.GetAllMachines(names); err != nil {
		t.Errorf("Error in GetAllMachines(string): %s", err)
	} else {
		if len(machine2) != 1 {
			t.Error("Filter for machines in GetAllMachines([]string) seems not to work")
			return
		}
	}

	return
}

// TestReapJob tests job reaping from internal lists by calling the job's
// Reap() method.
func TestReapJob(t *testing.T) {
	var sm drmaa2.SessionManager
	var err error
	var js *drmaa2.JobSession

	// create or open job session
	if js, err = sm.OpenJobSession("TestReapJob"); err != nil {
		if js, err = sm.CreateJobSession("TestReapJob", ""); err != nil {
			t.Fatal(err)
		}
	}
	defer sm.DestroyJobSession("TestReapJob")

	var jt drmaa2.JobTemplate
	jt.RemoteCommand = "/bin/sleep"
	jt.Args = []string{"1"}

	jt.OutputPath = "/dev/null"
	jt.JoinFiles = true

	job, errRun := js.RunJob(jt)
	if errRun != nil {
		t.Fatal(errRun)
	}

	// wait until sleep is finished
	job.WaitTerminated(drmaa2.InfiniteTime)

	// it finsihed jobs appear in all job lists (GetJobs / monitoring session GetAllJobs)
	jl, errJL := js.GetJobs(nil)
	if errJL != nil {
		t.Fatalf("Error during GetJobs(): %s\n", errJL)
	}

	if len(jl) != 1 {
		t.Logf("Job list must be 1 but it is %d\n", len(jl))
	}

	if errReap := job.Reap(); errReap != nil {
		t.Fatalf("Reaping of job caused an error: %s\n", errReap)
	}
	t.Log("Reaping of job successful")

	jl, errJL = js.GetJobs(nil)

	if len(jl) != 0 {
		t.Fatalf("Job list still contains reaped jobs: %d\n", len(jl))
	}
}

func TestGetJobTemplate(t *testing.T) {
	var jt drmaa2.JobTemplate
	var sm drmaa2.SessionManager
	var js *drmaa2.JobSession
	var err error

	if js, err = sm.OpenJobSession("TestGetJobTemplate"); err != nil {
		if js, err = sm.CreateJobSession("TestGetJobTemplate", ""); err != nil {
			t.Errorf("Failed when creating job session: %s\n", err)
			return
		}
	}
	defer sm.DestroyJobSession("TestGetJobTemplate")

	jt.JobEnvironment = make(map[string]string, 0)
	jt.JobEnvironment["one"] = "1"
	jt.JobEnvironment["two"] = "2"
	jt.JobEnvironment["tree"] = "3"

	jt.RemoteCommand = "/bin/sleep"
	jt.Args = []string{"0"}

	job, _ := js.RunJob(jt)
	template, _ := job.GetJobTemplate()

	env := template.JobEnvironment

	for k, v := range jt.JobEnvironment {
		if env[k] != v {
			t.Errorf("JobEnvironment is not correctly recovered (%s != %s)\n", env[k], v)
		} else {
			fmt.Println("Found environment variable.")
		}
	}
}

func TestOpenCloseSession(t *testing.T) {
	// randomly open and closing a session
	var sm drmaa2.SessionManager

	js, errcjs := sm.CreateJobSession("1", "")
	if errcjs != nil {
		t.Errorf("Error during CreateJobsSession(1): %s", errcjs)
	}
	ms1, errms := sm.OpenMonitoringSession("")
	if errms != nil {
		t.Errorf("Error during OpenMonitoringSession(): %s", errms)
	}
	js2, errcjs2 := sm.CreateJobSession("2", "")
	if errcjs2 != nil {
		t.Errorf("Error during CreateJobsSession(2): %s", errcjs2)
	}
	if errmsc := ms1.CloseMonitoringSession(); errmsc != nil {
		t.Errorf("Error during CloseMonitoringSession(): %s", errmsc)
	}
	js3, errcjs3 := sm.CreateJobSession("3", "")
	if errcjs3 != nil {
		t.Errorf("Error during CreateJobsSession(3): %s", errcjs3)
	}
	ms2, errms2 := sm.OpenMonitoringSession("")
	if errms2 != nil {
		t.Errorf("Error during OpenMonitoringSession(): %s", errms2)
	}
	if errmsc2 := ms2.CloseMonitoringSession(); errmsc2 != nil {
		t.Errorf("Error during CloseMonitoringSession(): %s", errmsc2)
	}
	if errcljs := js.Close(); errcljs != nil {
		t.Errorf("Error during js.Close(): %s", errcljs)
	}
	if errcljs3 := js3.Close(); errcljs3 != nil {
		t.Errorf("Error during js.Close(): %s", errcljs3)
	}
	if errcljs2 := js2.Close(); errcljs2 != nil {
		t.Errorf("Error during js.Close(): %s", errcljs2)
	}
}

// TODO add more :)
