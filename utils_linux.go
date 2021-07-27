package connstate

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

func diagTCPInfo(family uint8) ([]TCPState, error) {
	// Request TCP diag
	res, err := netlink.SocketDiagTCPInfo(family)
	if err != nil {
		return nil, err
	}
	var list []TCPState
	for _, i := range res {
		if i.InetDiagMsg != nil {
			list = append(list, TCPState{
				Socket: *i.InetDiagMsg,
			})
		}
	}
	return list, nil
}

func diagTCPStatistics(family uint8) (*TCPStatistics, error) {
	// Request TCP diag
	res, err := netlink.SocketDiagTCP(family)
	if err != nil {
		return nil, err
	}
	statistics := NewTCPStatistics()
	for _, socket := range res {
		statistics.Record(socket)
	}
	return statistics, nil
}

// executeInNetns sets execution of the code following this call to the
// network namespace newNs, then moves the thread back to curNs if open,
// otherwise to the current netns at the time the function was invoked
// In case of success, the caller is expected to execute the returned function
// at the end of the code that needs to be executed in the network namespace.
// Example:
// func jobAt(...) error {
//      d, err := executeInNetns(...)
//      if err != nil { return err}
//      defer d()
//      < code which needs to be executed in specific netns>
//  }
// Shadow from netns package
func executeInNetns(newNs, curNs netns.NsHandle) (func(), error) {
	var (
		err       error
		moveBack  func(netns.NsHandle) error
		closeNs   func() error
		unlockThd func()
	)
	restore := func() {
		// order matters
		if moveBack != nil {
			moveBack(curNs)
		}
		if closeNs != nil {
			closeNs()
		}
		if unlockThd != nil {
			unlockThd()
		}
	}
	if newNs.IsOpen() {
		runtime.LockOSThread()
		unlockThd = runtime.UnlockOSThread
		if !curNs.IsOpen() {
			if curNs, err = netns.Get(); err != nil {
				restore()
				return nil, fmt.Errorf("could not get current namespace while creating netlink socket: %v", err)
			}
			closeNs = curNs.Close
		}
		if err := netns.Set(newNs); err != nil {
			restore()
			return nil, fmt.Errorf("failed to set into network namespace %d while creating netlink socket: %v", newNs, err)
		}
		moveBack = netns.Set
	}
	return restore, nil
}

// Copy from netns
func findCgroupMountpoint(cgroupType string) (string, error) {
	output, err := ioutil.ReadFile("/proc/mounts")
	if err != nil {
		return "", err
	}

	// /proc/mounts has 6 fields per line, one mount per line, e.g.
	// cgroup /sys/fs/cgroup/devices cgroup rw,relatime,devices 0 0
	for _, line := range strings.Split(string(output), "\n") {
		parts := strings.Split(line, " ")
		if len(parts) == 6 && parts[2] == "cgroup" {
			for _, opt := range strings.Split(parts[3], ",") {
				if opt == cgroupType {
					return parts[1], nil
				}
			}
		}
	}

	return "", fmt.Errorf("cgroup mountpoint not found for %s", cgroupType)
}

func getPidFormCgroupTask(filename string) (int, error) {
	var PID int
	file, err := os.Open(filename)
	if err != nil {
		return PID, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		task := scanner.Text()
		PID, err = strconv.Atoi(task)
		if err != nil {
			return PID, fmt.Errorf("invalid pid '%s': %s", task, err)
		}
		// Read the first PID
		break
	}
	if err = scanner.Err(); err != nil {
		return PID, err
	}
	if PID == 0 {
		return PID, ErrFailedToGetPIDFromCgroup
	}
	return PID, nil
}
