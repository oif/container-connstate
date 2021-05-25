package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/oif/container-connstate"

	"github.com/containerd/containerd"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var containerdSocket string
	flag.StringVar(&containerdSocket, "containerd", "/run/containerd/containerd.sock", "containerd daemon socket")
	flag.Parse()

	containerdClient, err := containerd.New(containerdSocket)
	panicOnError(err)
	driver, err := connstate.NewContainerdDriver(containerdClient, connstate.WithEnvCollectionFilter(func(s string) bool {
		return strings.HasPrefix(s, "CONTAINER_")
	}))
	panicOnError(err)
	tracker, err := connstate.NewTracker(driver)
	panicOnError(err)
	containers, err := tracker.ListAllConnectionState()
	panicOnError(err)
	for _, container := range containers {
		statistics := make(map[string]uint64)
		v4Statistics := container.IPv4.StatisticsByState()
		v6Statistics := container.IPv6.StatisticsByState()
		for stateFlag, count := range v4Statistics {
			statistics[connstate.GetReadableState(stateFlag)] += count
		}
		for stateFlag, count := range v6Statistics {
			statistics[connstate.GetReadableState(stateFlag)] += count
		}
		fmt.Printf("%s@%d -> %v\n", container.ID, container.PID, statistics)
		fmt.Printf("\t Annotations %v\n", container.Annotations)
		fmt.Println("\t Connections")
		for _, connection := range append(container.IPv4, container.IPv6...) {
			fmt.Printf("\t%v %s\n", connection.Socket.ID, connstate.GetReadableState(connection.Socket.State))
		}
	}
}
