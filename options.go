package connstate

type OptionsFunc func(*Options) error

type Options struct {
	fixedCgroupRoot string
	// Default(nil) for not collecting ENV to avoid performance issue
	_ENVCollectionFilter func(ENV string) bool
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

func WithEnvCollectionFilter(filter func(string) bool) OptionsFunc {
	return func(o *Options) error {
		o._ENVCollectionFilter = filter
		return nil
	}
}
