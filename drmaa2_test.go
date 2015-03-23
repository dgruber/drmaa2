package drmaa2_test

import (
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
	if machine, err := ms.GetAllMachines(nil); err != nil {
		t.Errorf("Error during GetAllMachines(nil): %s", err)
		return
	} else {
		amount := len(machine)
		if amount < 1 {
			t.Errorf("Error: No machine returned in GetAllMachines(nil)")
		}
		// get a single machine
		names := make([]string, 0)
		names = append(names, machine[0].Name)
		if machine2, err := ms.GetAllMachines(names); err != nil {
			t.Errorf("Error in GetAllMachines(string): %s", err)
		} else {
			if len(machine2) != 1 {
				t.Error("Filter for machines in GetAllMachines([]string) seems not to work")
				return
			}
		}
	}
	return
}

// TODO add more :)
