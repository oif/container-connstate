package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"syscall"

	"github.com/oif/container-connstate"

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
	containerStates, err := tracker.ListAllConnectionState(context.TODO())
	panicOnError(err)
	for _, container := range containerStates {
		fmt.Printf("%s@%d\n", container.ID, container.PID)
		fmt.Printf("\t Annotations %v\n", container.Annotations)
		fmt.Println("\t Connections")
		for _, connection := range append(container.IPv4, container.IPv6...) {
			fmt.Printf("\t%v %s\n", connection.Socket.ID, connstate.GetReadableState(connection.Socket.State))
		}
	}

	containerStatistics, err := tracker.ListAllConnectionStatistics(context.TODO())
	panicOnError(err)
	for _, container := range containerStatistics {
		statistics := make(map[string]uint64)
		for stateFlag, count := range container.IPv4 {
			statistics[connstate.GetReadableState(stateFlag)] += count
		}
		for stateFlag, count := range container.IPv6 {
			statistics[connstate.GetReadableState(stateFlag)] += count
		}
		fmt.Printf("%s@%d -> %v\n", container.ID, container.PID, statistics)
	}
}
