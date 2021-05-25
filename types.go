package connstate

import (
	"context"

	"github.com/vishvananda/netlink"
)

type Container struct {
	ID  string
	PID int
	// Stores extension data
	Annotations map[string]string
}

type ContainerDriver interface {
	ListContainer(context.Context) ([]Container, error)
}

type ConnectionState struct {
	Container
	IPv4 TCPStates
	IPv6 TCPStates
}

type TCPStates []TCPState

func (s TCPStates) StatisticsByState() map[uint8]uint64 {
	mapping := make(map[uint8]uint64)
	for _, state := range s {
		mapping[state.Socket.State]++
	}
	return mapping
}

type TCPState struct {
	Socket netlink.Socket
}
