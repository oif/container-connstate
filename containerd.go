package connstate

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/typeurl"
	"github.com/opencontainers/runtime-spec/specs-go"
)

type ContainerdDriver struct {
	cgroupResourcePath string
	client             *containerd.Client

	opts *Options
}

func NewContainerdDriver(client *containerd.Client, opts ...OptionsFunc) (*ContainerdDriver, error) {
	d := &ContainerdDriver{
		client: client,
		opts:   NewDefaultOptions(),
	}
	for _, optFn := range opts {
		if err := optFn(d.opts); err != nil {
			return nil, err
		}
	}
	const cgroupType = "memory"
	if d.opts.fixedCgroupRoot != "" {
		d.cgroupResourcePath = filepath.Join(d.opts.fixedCgroupRoot, cgroupType)
	} else {
		var err error
		d.cgroupResourcePath, err = findCgroupMountpoint(cgroupType)
		if err != nil {
			return nil, err
		}
	}
	return d, nil
}

func (d *ContainerdDriver) ListContainer(ctx context.Context) ([]Container, error) {
	containerdNamespaces, err := d.client.NamespaceService().List(ctx)
	if err != nil {
		return nil, err
	}

	var containerList []Container
	for _, namespace := range containerdNamespaces {
		containers, err := d.client.ContainerService().List(namespaces.WithNamespace(ctx, namespace))
		if err != nil {
			return nil, err
		}
		for _, container := range containers {
			receiver, err := typeurl.UnmarshalAny(container.Spec)
			if err != nil {
				return nil, err
			}
			spec, ok := receiver.(*specs.Spec)
			if !ok {
				return nil, fmt.Errorf("container %s spec is not a runtime spec", container.ID)
			}
			containerPID, err := getPidFormCgroupTask(filepath.Join(d.cgroupResourcePath, spec.Linux.CgroupsPath, "tasks"))
			if err != nil {
				return nil, err
			}
			containerList = append(containerList, Container{
				ID:          container.ID,
				PID:         containerPID,
				Annotations: nil,
			})
		}
	}
	return containerList, nil
}
