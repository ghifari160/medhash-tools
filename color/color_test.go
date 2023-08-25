package color_test

import (
	"testing"

	"github.com/ghifari160/medhash-tools/color"
	"github.com/stretchr/testify/suite"
)

const payloadColor = color.Blue + "This " +
	color.Green + "is " +
	color.Magenta + "a " +
	color.Gray + "test" +
	color.Reset + "\n"

type ColorSuite struct {
	suite.Suite
}

func TestColor(t *testing.T) {
	s := new(ColorSuite)

	suite.Run(t, s)
}

func (s *ColorSuite) TestClean() {
	clean := color.Clean([]byte(payloadColor))
	s.NotContains(clean, color.EscStr)
}

func (s *ColorSuite) TestCleanString() {
	clean := color.CleanString(payloadColor)
	s.NotContains(clean, color.EscStr)
}
