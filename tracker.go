package connstate

import (
	"github.com/vishvananda/netns"
)

type Tracker struct {
	driver                  ContainerDriver
	currentNetworkNamespace netns.NsHandle
	families                []uint8
}

func NewTracker(driver ContainerDriver, families []uint8) (*Tracker, error) {
	t := &Tracker{
		driver:   driver,
		families: families,
	}
	// Get current runtime network namespace for restore after enter other ns
	currentNetNS, err := netns.Get()
	if err != nil {
		return nil, err
	}
	t.currentNetworkNamespace = currentNetNS
	return t, nil
}

func (t *Tracker) executeInContainerNetworkNamespace(container Container, execution func(netns.NsHandle) error) error {
	if container.PID == 0 {
		return ErrInvalidPID
	}
	// Get netns for executing
	nsHandler, err := netns.GetFromPid(container.PID)
	if err != nil {
		return err
	}
	defer func() {
		_ = nsHandler.Close()
	}()
	// enter target netns
	restore, err := EnterNetNS(nsHandler, t.currentNetworkNamespace)
	if err != nil {
		return err
	}
	// back to runtime ns to avoid affect other operation
	defer restore()
	return execution(nsHandler)
}
