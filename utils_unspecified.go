//go:build !linux
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

func EnterNetNS(_, _ netns.NsHandle) (func(), error) {
	return nil, ErrNotImplementYet
}

func diagTCPInfo(_ uint8) ([]TCPState, error) {
	return nil, ErrNotImplementYet
}

func diagTCPStatistics(_ uint8) (*TCPStatistics, error) {
	return nil, ErrNotImplementYet
}
