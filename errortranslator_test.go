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

type ErrorTranslatorSuite struct{}

var _ = Suite(&ErrorTranslatorSuite{})

func (s *ErrorTranslatorSuite) TestTranslateError(c *C) {
	et := errormapper.ErrorTranslator{
		validate.ErrRequired: "required translate",
	}

	trans, ok := et.TranslateError(validate.ErrRequired)
	c.Assert(trans, Equals, "required translate")
	c.Assert(ok, Equals, true)

	trans, ok = et.TranslateError(validate.ErrMin)
	c.Assert(trans, Equals, "")
	c.Assert(ok, Equals, false)
}

func (s *ErrorTranslatorSuite) TestTranslateErrorWitDefault(c *C) {
	et := errormapper.ErrorTranslator{
		validate.ErrRequired: "required translate",
		nil:                  "nil default translate",
	}

	trans, ok := et.TranslateError(validate.ErrRequired)
	c.Assert(trans, Equals, "required translate")
	c.Assert(ok, Equals, true)

	trans, ok = et.TranslateError(validate.ErrMin)
	c.Assert(trans, Equals, "nil default translate")
	c.Assert(ok, Equals, true)
}

func (s *ErrorTranslatorSuite) TestTranslate(c *C) {
	et := errormapper.ErrorTranslator{
		validate.ErrRequired: "required translate",
		validate.ErrMin:      "min translate",
	}

	tests := []struct {
		Description string
		Errors      validate.Errors
		ExpectedOk  bool
		Expected    string
	}{
		{
			Description: "default",
			Errors:      validate.Errors{validate.ErrRequired, validate.ErrMin},
			ExpectedOk:  true,
			Expected:    "required translate, min translate",
		}, {
			Description: "order test",
			Errors:      validate.Errors{validate.ErrMin, validate.ErrRequired},
			ExpectedOk:  true,
			Expected:    "min translate, required translate",
		}, {
			Description: "only known errors are translated",
			Errors:      validate.Errors{validate.ErrMax, validate.ErrRequired},
			ExpectedOk:  true,
			Expected:    "required translate",
		}, {
			Description: "No translation at all",
			Errors:      validate.Errors{validate.ErrMax},
			ExpectedOk:  false,
			Expected:    "",
		},
	}

	for _, test := range tests {
		translated, ok := et.Translate(test.Errors)

		c.Assert(ok, Equals, test.ExpectedOk, Commentf(test.Description))
		c.Assert(translated, Equals, test.Expected, Commentf(test.Description))
	}
}

func (s *ErrorTranslatorSuite) TestTranslateWithDefault(c *C) {
	//a translation with a default (nil) translation always succeeds translation
	et := errormapper.ErrorTranslator{
		validate.ErrRequired: "required translate",
		validate.ErrMin:      "min translate",
		nil:                  "nil default translate",
	}

	tests := []struct {
		Description string
		Errors      validate.Errors
		Expected    string
	}{
		{
			Description: "default",
			Errors:      validate.Errors{validate.ErrRequired, validate.ErrMin},
			Expected:    "required translate, min translate",
		}, {
			Description: "order test, only first error",
			Errors:      validate.Errors{validate.ErrMin, validate.ErrRequired},
			Expected:    "min translate, required translate",
		}, {
			Description: "fallback to nil translation on unkown error",
			Errors:      validate.Errors{validate.ErrMax, validate.ErrRequired},
			Expected:    "nil default translate, required translate",
		},
	}

	for _, test := range tests {
		translated, ok := et.Translate(test.Errors)

		c.Assert(translated, Equals, test.Expected, Commentf(test.Description))
		c.Assert(ok, Equals, true, Commentf(test.Description))
	}
}

func (s *ErrorTranslatorSuite) TestTranslateFirst(c *C) {
	et := errormapper.ErrorTranslator{
		validate.ErrRequired: "required translate",
		validate.ErrMin:      "min translate",
	}

	tests := []struct {
		Description string
		Errors      validate.Errors
		ExpectedOk  bool
		Expected    string
	}{
		{
			Description: "first error",
			Errors:      validate.Errors{validate.ErrRequired, validate.ErrMin},
			ExpectedOk:  true,
			Expected:    "required translate",
		}, {
			Description: "order test, only first error",
			Errors:      validate.Errors{validate.ErrMin, validate.ErrRequired},
			ExpectedOk:  true,
			Expected:    "min translate",
		}, {
			Description: "only known errors are translated",
			Errors:      validate.Errors{validate.ErrMax, validate.ErrRequired},
			ExpectedOk:  true,
			Expected:    "required translate",
		}, {
			Description: "No translation at all",
			Errors:      validate.Errors{validate.ErrMax},
			ExpectedOk:  false,
			Expected:    "",
		},
	}

	for _, test := range tests {
		translated, ok := et.TranslateFirst(test.Errors)

		c.Assert(ok, Equals, test.ExpectedOk, Commentf(test.Description))
		c.Assert(translated, Equals, test.Expected, Commentf(test.Description))
	}
}

func (s *ErrorTranslatorSuite) TestTranslateFirstWithDefault(c *C) {
	//a translation with a default (nil) translation always succeeds translation
	et := errormapper.ErrorTranslator{
		validate.ErrRequired: "required translate",
		validate.ErrMin:      "min translate",
		nil:                  "nil default translate",
	}

	tests := []struct {
		Description string
		Errors      validate.Errors
		Expected    string
	}{
		{
			Description: "first error",
			Errors:      validate.Errors{validate.ErrRequired, validate.ErrMin},
			Expected:    "required translate",
		}, {
			Description: "order test, only first error",
			Errors:      validate.Errors{validate.ErrMin, validate.ErrRequired},
			Expected:    "min translate",
		}, {
			Description: "fallback to nil translation on unkown error",
			Errors:      validate.Errors{validate.ErrMax, validate.ErrRequired},
			Expected:    "nil default translate",
		},
	}

	for _, test := range tests {
		translated, ok := et.TranslateFirst(test.Errors)

		c.Assert(translated, Equals, test.Expected, Commentf(test.Description))
		c.Assert(ok, Equals, true, Commentf(test.Description))
	}
}
