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

type TCPStatistics struct {
	allocatedPorts map[uint16]bool
	states         map[uint8]uint64
}

func NewTCPStatistics() *TCPStatistics {
	return &TCPStatistics{
		allocatedPorts: make(map[uint16]bool),
		states:         make(map[uint8]uint64),
	}
}

func (s *TCPStatistics) Record(socket *netlink.Socket) {
	if socket == nil {
		return
	}
	s.allocatedPorts[socket.ID.SourcePort] = true
	s.states[socket.State]++
}

func (s TCPStatistics) AllocatedPortCount() int {
	counter := 0
	for _, allocated := range s.allocatedPorts {
		if allocated {
			counter++
		}
	}
	return counter
}

func (s TCPStatistics) CountByState() map[uint8]uint64 {
	shadow := make(map[uint8]uint64)
	for k, v := range s.states {
		shadow[k] = v
	}
	return shadow
}

type TCPStates []TCPState

type TCPState struct {
	Socket netlink.Socket
}

// Return false if wanna bypass
type ContainerFilter func(Container) (pass bool)
