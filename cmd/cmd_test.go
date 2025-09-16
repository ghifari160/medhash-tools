package cmd_test

import (
	"github.com/ghifari160/medhash-tools/medhash"
	"github.com/ghifari160/medhash-tools/testcommon"
)

// withInvalidate invalidates the payload hash for testing.
func withInvalidate(invalidate bool) testcommon.Options {
	return testcommon.NewOptions("invalidate", invalidate)
}

func withForce(force bool) testcommon.Options {
	return testcommon.NewOptions("force", force)
}

// withFiles specifies the Files argument to a command for testing.
func withFiles(files []string) testcommon.Options {
	return testcommon.NewOptions("files", files)
}

func withGenConfig(config medhash.Config) testcommon.Options {
	return testcommon.NewOptions("config_gen", config)
}

func genConfig(options testcommon.Options) medhash.Config {
	if v, ok := options.Raw("config_gen").(medhash.Config); ok {
		return v
	} else {
		return medhash.Config{}
	}
}

func withChkConfig(config medhash.Config) testcommon.Options {
	return testcommon.NewOptions("config_chk", config)
}

func chkConfig(options testcommon.Options) medhash.Config {
	if v, ok := options.Raw("config_chk").(medhash.Config); ok {
		return v
	} else {
		return medhash.Config{}
	}
}
