package errormapper_test

import (
	errormapper "github.com/mbict/go-errormapper"
	validate "github.com/mbict/go-validate"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) {
	TestingT(t)
}

type ErrorMapperSuite struct {
	errMap errormapper.FieldTranslations
}

var _ = Suite(&ErrorMapperSuite{})

func (s *ErrorMapperSuite) SetUpTest(c *C) {

	s.errMap = make(errormapper.FieldTranslations)
	s.errMap.AddDefaultTranslation("A", "A Nil")
	s.errMap.AddTranslation("B", validate.ErrRequired, "B Required")
	s.errMap.AddTranslation("B", validate.ErrMax, "B Max")
	s.errMap.AddDefaultTranslation("B", "B Nil")
	s.errMap.AddTranslation("C", validate.ErrRequired, "C Required")
	s.errMap.AddTranslation("C", validate.ErrMax, "C Max")
}

func (s *ErrorMapperSuite) TestTranslationFirstFallbackToDefault(c *C) {
	em := validate.ErrorMap{
		"A": {validate.ErrRequired, validate.ErrMin},
		"B": {validate.ErrMin},
		"C": {validate.ErrMin},
	}

	translated := s.errMap.TranslateFirst(em)

	c.Assert(translated, HasLen, 2)
	c.Assert(translated["A"], Equals, "A Nil")
	c.Assert(translated["B"], Equals, "B Nil")
}

func (s *ErrorMapperSuite) TestTranslateFirst(c *C) {
	em := validate.ErrorMap{
		"B": {validate.ErrRequired},
		"C": {validate.ErrMax},
	}

	translated := s.errMap.TranslateFirst(em)

	c.Assert(translated, HasLen, 2)
	c.Assert(translated["B"], Equals, "B Required")
	c.Assert(translated["C"], Equals, "C Max")
}
