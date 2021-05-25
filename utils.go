package connstate

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

func GetReadableState(state uint8) string {
	switch state {
	case netlink.TCP_ESTABLISHED:
		return "ESTABLISHED"
	case netlink.TCP_SYN_SENT:
		return "SYN_SENT"
	case netlink.TCP_SYN_RECV:
		return "SYN_RECV"
	case netlink.TCP_FIN_WAIT1:
		return "FIN_WAIT1"
	case netlink.TCP_FIN_WAIT2:
		return "FIN_WAIT2"
	case netlink.TCP_TIME_WAIT:
		return "TIME_WAIT"
	case netlink.TCP_CLOSE:
		return "CLOSE"
	case netlink.TCP_CLOSE_WAIT:
		return "CLOSE_WAIT"
	case netlink.TCP_LAST_ACK:
		return "LAST_ACK"
	case netlink.TCP_LISTEN:
		return "LISTEN"
	case netlink.TCP_CLOSING:
		return "CLOSING"
	case netlink.TCP_NEW_SYN_REC:
		return "NEW_SYN_REC"
	case netlink.TCP_MAX_STATES:
		return "MAX_STATES"
	default:
		return fmt.Sprintf("unknown state %x", state)
	}
}
