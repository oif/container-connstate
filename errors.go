package connstate

import "errors"

var (
	ErrInvalidPID               = errors.New("invalid PID")
	ErrFailedToGetPIDFromCgroup = errors.New("failed to get PID from cgroup")
	ErrNotImplementYet          = errors.New("not implement yet")
)
