package testcommon

import "fmt"

// Options configure a testing case.
type Options struct {
	config map[string]any
}

func (opt Options) IsBool(key string) bool {
	_, ok := opt.config[key].(bool)
	return ok
}

func (opt Options) IsStr(key string) bool {
	_, ok := opt.config[key].(string)
	return ok
}

func (opt Options) IsStrSlice(key string) bool {
	_, ok := opt.config[key].([]string)
	return ok
}

func (opt Options) Bool(key string) bool {
	if v, ok := opt.config[key].(bool); ok {
		return v
	} else {
		return false
	}
}

func (opt Options) Str(key string) string {
	if v, ok := opt.config[key].(string); ok {
		return v
	} else if opt.config[key] == nil {
		return ""
	} else {
		return fmt.Sprintf("%v", opt.config[key])
	}
}

func (opt Options) StrSlice(key string) []string {
	if v, ok := opt.config[key].([]string); ok {
		return v
	} else {
		return []string{}
	}
}

func (opt Options) Raw(key string) any {
	return opt.config[key]
}

// NewOptions creates a new singular test option.
func NewOptions(key string, value any) Options {
	opts := Options{
		config: make(map[string]any),
	}
	opts.config[key] = value
	return opts
}

// MergeOptions merges a set of Options into one Options containing all the configurations.
// Note that no validation is done.
// Duplicate key-value pairs are overwritten by the one added last.
func MergeOptions(opts ...Options) Options {
	config := make(map[string]any)
	for _, opt := range opts {
		for k, v := range opt.config {
			config[k] = v
		}
	}
	return Options{
		config: config,
	}
}
