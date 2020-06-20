package popper

import (
	"testing"

	"gopkg.in/check.v1"
)

type PopperSuite struct{}

var _ = check.Suite(new(PopperSuite))

func Test(t *testing.T) { check.TestingT(t) }

func (s *PopperSuite) SetUpSuite(c *check.C) {
}

func (s *PopperSuite) TearDownSuite(c *check.C) {
}

func (s *PopperSuite) SetUpTest(c *check.C) {
}

func (s *PopperSuite) TearDownTest(c *check.C) {
}

func (s *PopperSuite) TestNew(c *check.C) {
	c.Assert(true, check.Equals, true)
}
