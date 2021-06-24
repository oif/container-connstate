package connstate

import (
	"context"
	"fmt"
	"syscall"

	"github.com/vishvananda/netns"
)

func (t *Tracker) GetConnectionState(container Container) (*ConnectionState, error) {
	state := ConnectionState{
		Container: container,
	}
	err := t.executeInContainerNetworkNamespace(container, func(nsHandler netns.NsHandle) error {
		state.NetNSID = nsHandler.UniqueId()
		for _, family := range t.families {
			var receiver *TCPStates
			switch family {
			case syscall.AF_INET:
				receiver = &state.IPv4
			case syscall.AF_INET6:
				receiver = &state.IPv6
			default:
				return fmt.Errorf("unsupported family %x", family)
			}
			result, err := diagTCPInfo(family)
			if err != nil {
				return err
			}
			*receiver = result
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func (t *Tracker) ListAllConnectionState(ctx context.Context, filter ContainerFilter) ([]ConnectionState, error) {
	containers, err := t.driver.ListContainer(ctx)
	if err != nil {
		return nil, err
	}
	var list []ConnectionState
	for _, container := range containers {
		if filter != nil && !filter(container) {
			// bypass
			continue
		}
		state, err := t.GetConnectionState(container)
		if err != nil {
			fmt.Printf("Failed to get container(%s) state: %s\n", container.Hostname, err)
			continue
		}
		list = append(list, *state)
	}
	return list, nil
}
