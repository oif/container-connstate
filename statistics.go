package connstate

import (
	"context"
	"fmt"
	"syscall"

	"github.com/vishvananda/netns"
)

func (t *Tracker) GetConnectionStatistics(container Container) (*ConnectionStatistics, error) {
	statistics := ConnectionStatistics{
		Container: container,
	}
	err := t.executeInContainerNetworkNamespace(container, func(nsHandler netns.NsHandle) error {
		statistics.NetNSID = nsHandler.UniqueId()
		for _, family := range t.families {
			var receiver *TCPStatistics
			switch family {
			case syscall.AF_INET:
				receiver = &statistics.IPv4
			case syscall.AF_INET6:
				receiver = &statistics.IPv6
			default:
				return fmt.Errorf("unsupported family %x", family)
			}
			result, err := diagTCPStatistics(family)
			if err != nil {
				return err
			}
			*receiver = *result
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &statistics, nil
}

func (t *Tracker) ListAllConnectionStatistics(ctx context.Context, filter ContainerFilter) ([]ConnectionStatistics, error) {
	containers, err := t.driver.ListContainer(ctx)
	if err != nil {
		return nil, err
	}
	var list []ConnectionStatistics
	for _, container := range containers {
		if filter != nil && !filter(container) {
			// bypass
			continue
		}
		statistics, err := t.GetConnectionStatistics(container)
		if err != nil {
			fmt.Printf("Failed to get container(%s) statistics: %s\n", container.Hostname, err)
			continue
		}
		list = append(list, *statistics)
	}
	return list, nil
}
