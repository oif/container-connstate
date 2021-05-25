package main

import (
	"fmt"

	"github.com/oif/container-connstate"

	"github.com/containerd/containerd"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	containerdClient, err := containerd.New("/run/containerd/containerd.sock")
	panicOnError(err)
	driver, err := connstate.NewContainerdDriver(containerdClient)
	panicOnError(err)
	tracker, err := connstate.NewTracker(driver)
	panicOnError(err)
	states, err := tracker.ListAllConnectionState()
	panicOnError(err)
	for _, state := range states {
		statistics := make(map[string]uint64)
		v4Statistics := state.IPv4.StatisticsByState()
		v6Statistics := state.IPv6.StatisticsByState()
		for stateFlag, count := range v4Statistics {
			statistics[connstate.GetReadableState(stateFlag)] += count
		}
		for stateFlag, count := range v6Statistics {
			statistics[connstate.GetReadableState(stateFlag)] += count
		}
		fmt.Printf("%s@%d -> %v\n", state.ID, state.PID, statistics)
		for _, connection := range append(state.IPv4, state.IPv6...) {
			fmt.Printf("\t%v %s\n", connection.Socket.ID, connstate.GetReadableState(connection.Socket.State))
		}
	}
}
