// +build !linux

package connstate

import (
	"github.com/vishvananda/netns"
)

func findCgroupMountpoint(_ string) (string, error) {
	return "", ErrNotImplementYet
}

func getPidFormCgroupTask(_ string) (int, error) {
	return 0, ErrNotImplementYet
}

func executeInNetns(_, _ netns.NsHandle) (func(), error) {
	return nil, ErrNotImplementYet
}

func diagTCPInfo(_ uint8) ([]TCPState, error) {
	return nil, ErrNotImplementYet
}
