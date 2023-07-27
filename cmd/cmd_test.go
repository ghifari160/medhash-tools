package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CmdSuite struct {
	suite.Suite

	PayloadSize int64
}

func TestCmd(t *testing.T) {
	s := new(CmdSuite)

	if testing.Short() {
		s.PayloadSize = 1024
	} else {
		s.PayloadSize = 1 * 1024 * 1024 * 1024
	}

	suite.Run(t, s)
}
