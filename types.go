package connstate

import (
	"context"

	"github.com/vishvananda/netlink"
)

type Container struct {
	ID       string
	PID      int
	Hostname string
	// Stores extension data
	Annotations map[string]string
}

type ContainerDriver interface {
	ListContainer(context.Context) ([]Container, error)
}

type ConnectionState struct {
	Container
	NetNSID string
	IPv4    TCPStates
	IPv6    TCPStates
}

type ConnectionStatistics struct {
	Container
	NetNSID string
	IPv4    TCPStatistics
	IPv6    TCPStatistics
}

type TCPStatistics map[uint8]uint64

type TCPStates []TCPState

func (s TCPStates) StatisticsByState() TCPStatistics {
	mapping := make(TCPStatistics)
	for _, state := range s {
		mapping[state.Socket.State]++
	}
	return mapping
}

type TCPState struct {
	Socket netlink.Socket
}
