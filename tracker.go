package connstate

import (
	"context"
	"fmt"
	"syscall"

	"github.com/vishvananda/netns"
)

type Tracker struct {
	driver                  ContainerDriver
	currentNetworkNamespace netns.NsHandle
}

func NewTracker(driver ContainerDriver) (*Tracker, error) {
	t := &Tracker{
		driver: driver,
	}
	// Get current runtime network namespace for restore after enter other ns
	currentNetNS, err := netns.Get()
	if err != nil {
		return nil, err
	}
	t.currentNetworkNamespace = currentNetNS
	return t, nil
}

func (t *Tracker) GetConnectionState(container Container) (*ConnectionState, error) {
	if container.PID == 0 {
		return nil, ErrInvalidPID
	}
	// Get netns for executing
	nsHandler, err := netns.GetFromPid(container.PID)
	if err != nil {
		return nil, err
	}
	// enter target netns
	restore, err := executeInNetns(nsHandler, t.currentNetworkNamespace)
	if err != nil {
		return nil, err
	}
	// back to runtime ns to avoid affect other operation
	defer restore()
	state := ConnectionState{
		Container: container,
		NetNSID:   nsHandler.UniqueId(),
	}
	// Load v4
	state.IPv4, err = diagTCPInfo(syscall.AF_INET)
	if err != nil {
		return nil, err
	}
	// Load v6
	state.IPv6, err = diagTCPInfo(syscall.AF_INET6)
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func (t *Tracker) ListAllConnectionState() ([]ConnectionState, error) {
	containers, err := t.driver.ListContainer(context.TODO())
	if err != nil {
		return nil, err
	}
	var list []ConnectionState
	for _, container := range containers {
		state, err := t.GetConnectionState(container)
		if err != nil {
			fmt.Printf("Failed to get container(%s) state: %s\n", container.Hostname, err)
			continue
		}
		list = append(list, *state)
	}
	return list, nil
}
