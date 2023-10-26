package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"syscall"

	connstate "github.com/oif/container-connstate"

	"github.com/containerd/containerd"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var (
		containerdSocket string
		cgroupRoot       string
	)
	flag.StringVar(&containerdSocket, "containerd", "/run/containerd/containerd.sock", "containerd daemon socket")
	flag.StringVar(&cgroupRoot, "cgroup-root", "/sys/fs/cgroup/", "Fixed Cgroup root path")
	flag.Parse()

	containerdClient, err := containerd.New(containerdSocket)
	panicOnError(err)
	driver, err := connstate.NewContainerdDriver(containerdClient, connstate.WithEnvCollectionFilter(func(s string) bool {
		return strings.HasPrefix(s, "CONTAINER_")
	}), connstate.WithFixedCgroupRoot(cgroupRoot))
	panicOnError(err)
	tracker, err := connstate.NewTracker(driver, []uint8{syscall.AF_INET, syscall.AF_INET6})
	panicOnError(err)
	containerStates, err := tracker.ListAllConnectionState(context.TODO(), nil)
	panicOnError(err)
	for _, container := range containerStates {
		fmt.Printf("%s@%d\n", container.ID, container.PID)
		fmt.Printf("\t Annotations %v\n", container.Annotations)
		fmt.Println("\t Connections")
		for _, connection := range append(container.IPv4, container.IPv6...) {
			fmt.Printf("\t%v %s tx: %d rx: %d\n", connection.Socket.ID,
				connstate.GetReadableState(connection.Socket.State), connection.TXBytes, connection.RXBytes)
		}
	}

	containerStatistics, err := tracker.ListAllConnectionStatistics(context.TODO(), nil)
	panicOnError(err)
	for _, container := range containerStatistics {
		statistics := make(map[string]uint64)
		for stateFlag, count := range container.IPv4.CountByState() {
			statistics[connstate.GetReadableState(stateFlag)] += count
		}
		for stateFlag, count := range container.IPv6.CountByState() {
			statistics[connstate.GetReadableState(stateFlag)] += count
		}
		fmt.Printf("%s@%d -> %v(used %d port(s))\n",
			container.ID, container.PID, statistics,
			container.IPv4.AllocatedPortCount()+container.IPv6.AllocatedPortCount())
	}
}
