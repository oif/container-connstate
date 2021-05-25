package connstate

type OptionsFunc func(*Options) error

type Options struct {
	fixedCgroupRoot string
}

func NewDefaultOptions() *Options {
	return &Options{
		fixedCgroupRoot: "",
	}
}

func WithFixedCgroupRoot(path string) OptionsFunc {
	return func(o *Options) error {
		o.fixedCgroupRoot = path
		return nil
	}
}
