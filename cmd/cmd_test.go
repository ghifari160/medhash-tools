package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CmdSuite struct {
	suite.Suite
}

func TestCmd(t *testing.T) {
	suite.Run(t, new(CmdSuite))
}
