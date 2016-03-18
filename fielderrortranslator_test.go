package errormapper_test

import (
	errormapper "github.com/mbict/go-errormapper"
	validate "github.com/mbict/go-validate"
	. "gopkg.in/check.v1"
)

type FieldErrorTranslatorSuite struct{}

var _ = Suite(&FieldErrorTranslatorSuite{})

func (s *FieldErrorTranslatorSuite) TestAddTranslation(c *C) {
	et := errormapper.FieldErrorTranslator{}

	et.AddTranslation("A", validate.ErrRequired, "translation")

	c.Assert(et, DeepEquals, errormapper.FieldErrorTranslator{
		"A": errormapper.ErrorTranslator{validate.ErrRequired: "translation"},
	})

	//overwrite
	et.AddTranslation("A", validate.ErrRequired, "overritten")

	c.Assert(et, DeepEquals, errormapper.FieldErrorTranslator{
		"A": errormapper.ErrorTranslator{validate.ErrRequired: "overritten"},
	})

	//add more
	et.AddTranslation("A", validate.ErrMin, "min")

	c.Assert(et, DeepEquals, errormapper.FieldErrorTranslator{
		"A": errormapper.ErrorTranslator{
			validate.ErrRequired: "overritten",
			validate.ErrMin:      "min",
		},
	})

	//add other field more
	et.AddTranslation("B", validate.ErrMin, "min on b")
	c.Assert(et, DeepEquals, errormapper.FieldErrorTranslator{
		"A": errormapper.ErrorTranslator{
			validate.ErrRequired: "overritten",
			validate.ErrMin:      "min",
		},
		"B": errormapper.ErrorTranslator{
			validate.ErrMin: "min on b",
		},
	})
}

func (s *FieldErrorTranslatorSuite) TestSetFieldDefaultTranslation(c *C) {
	et := errormapper.FieldErrorTranslator{}

	et.SetFieldDefaultTranslation("A", "default field translation")

	c.Assert(et, DeepEquals, errormapper.FieldErrorTranslator{
		"A": errormapper.ErrorTranslator{nil: "default field translation"},
	})

	//overwrite
	et.SetFieldDefaultTranslation("A", "overwritten")

	c.Assert(et, DeepEquals, errormapper.FieldErrorTranslator{
		"A": errormapper.ErrorTranslator{nil: "overwritten"},
	})
}

func (s *FieldErrorTranslatorSuite) TestSetFallbackTranslation(c *C) {
	et := errormapper.FieldErrorTranslator{}

	et.SetFallbackTranslation(validate.ErrRequired, "default field err translation")

	c.Assert(et, DeepEquals, errormapper.FieldErrorTranslator{
		"": errormapper.ErrorTranslator{validate.ErrRequired: "default field err translation"},
	})

	//overwrite
	et.SetFallbackTranslation(validate.ErrRequired, "overwritten")

	c.Assert(et, DeepEquals, errormapper.FieldErrorTranslator{
		"": errormapper.ErrorTranslator{validate.ErrRequired: "overwritten"},
	})
}

func (s *FieldErrorTranslatorSuite) TestSetFallbackDefaultTranslation(c *C) {
	et := errormapper.FieldErrorTranslator{}

	et.SetFallbackDefaultTranslation("the absolute default")

	c.Assert(et, DeepEquals, errormapper.FieldErrorTranslator{
		"": errormapper.ErrorTranslator{nil: "the absolute default"},
	})

	//overwrite
	et.SetFallbackDefaultTranslation("overwritten")

	c.Assert(et, DeepEquals, errormapper.FieldErrorTranslator{
		"": errormapper.ErrorTranslator{nil: "overwritten"},
	})
}

func (s *FieldErrorTranslatorSuite) TestTranslate(c *C) {
	et := errormapper.FieldErrorTranslator{
		"A": errormapper.ErrorTranslator{
			validate.ErrRequired: "a required translate",
		},
		"B": errormapper.ErrorTranslator{
			validate.ErrRequired: "b required translate",
			validate.ErrMin:      "b min translate",
			validate.ErrMax:      "b max translate",
		},
	}

	tests := []struct {
		Description string
		Errors      validate.ErrorMap
		ExpectedOk  bool
		Expected    map[string]string
	}{
		{
			Description: "default",
			Errors: validate.ErrorMap{
				"A": validate.Errors{validate.ErrRequired, validate.ErrMin},
			},
			ExpectedOk: true,
			Expected: map[string]string{
				"A": "a required translate",
			},
		}, {
			Description: "2 fields",
			Errors: validate.ErrorMap{
				"A": validate.Errors{validate.ErrRequired, validate.ErrMin},
				"B": validate.Errors{validate.ErrMax, validate.ErrMin, validate.ErrMax},
			},
			ExpectedOk: true,
			Expected: map[string]string{
				"A": "a required translate",
				"B": "b max translate, b min translate, b max translate",
			},
		}, {
			Description: "1 field to translate no suitable translations found",
			Errors: validate.ErrorMap{
				"A": validate.Errors{validate.ErrMin},
			},
			ExpectedOk: false,
			Expected:   map[string]string{},
		}, {
			Description: "fails 1 field to translate no suitable translations found",
			Errors: validate.ErrorMap{
				"A": validate.Errors{validate.ErrMin},
			},
			ExpectedOk: false,
			Expected:   map[string]string{},
		}, {
			Description: "partial fails 1 field to translate no suitable translations found",
			Errors: validate.ErrorMap{
				"A": validate.Errors{validate.ErrMin},
				"B": validate.Errors{validate.ErrRequired},
			},
			ExpectedOk: false,
			Expected: map[string]string{
				"B": "b required translate",
			},
		}, {
			Description: "fails field has no translation set at all",
			Errors: validate.ErrorMap{
				"C": validate.Errors{validate.ErrMin},
			},
			ExpectedOk: false,
			Expected:   map[string]string{},
		},
	}

	for _, test := range tests {
		translated, ok := et.Translate(test.Errors)

		c.Assert(ok, Equals, test.ExpectedOk, Commentf(test.Description))
		c.Assert(translated, DeepEquals, test.Expected, Commentf(test.Description))
	}
}

func (s *FieldErrorTranslatorSuite) TestTranslateWithDefault(c *C) {
	//a translation with a default (nil) translation always succeeds translation
	et := errormapper.FieldErrorTranslator{
		"A": errormapper.ErrorTranslator{
			validate.ErrRequired: "a required translate",
		},
		"B": errormapper.ErrorTranslator{
			validate.ErrRequired: "b min translate",
			validate.ErrMin:      "b min translate",
			validate.ErrMax:      "b max translate",
		},
		"": errormapper.ErrorTranslator{
			validate.ErrMin: "nil min default translate",
		},
	}

	tests := []struct {
		Description string
		Errors      validate.ErrorMap
		ExpectedOk  bool
		Expected    map[string]string
	}{
		{
			Description: "fallback min error on A to field default translations",
			Errors: validate.ErrorMap{
				"A": validate.Errors{validate.ErrRequired, validate.ErrMin},
			},
			ExpectedOk: true,
			Expected: map[string]string{
				"A": "a required translate, nil min default translate",
			},
		}, {
			Description: "not set error field should fallback to default translations",
			Errors: validate.ErrorMap{
				"C": validate.Errors{validate.ErrMin},
			},
			ExpectedOk: true,
			Expected: map[string]string{
				"C": "nil min default translate",
			},
		}, {
			Description: "fallback no translation found should error out with a failed translation",
			Errors: validate.ErrorMap{
				"A": validate.Errors{validate.ErrRequired},
				"C": validate.Errors{validate.ErrRequired},
			},
			ExpectedOk: false,
			Expected: map[string]string{
				"A": "a required translate",
			},
		},
	}

	for _, test := range tests {
		translated, ok := et.Translate(test.Errors)

		c.Assert(ok, Equals, test.ExpectedOk, Commentf(test.Description))
		c.Assert(translated, DeepEquals, test.Expected, Commentf(test.Description))
	}
}

func (s *FieldErrorTranslatorSuite) TestTranslateFirst(c *C) {
	et := errormapper.FieldErrorTranslator{
		"A": errormapper.ErrorTranslator{
			validate.ErrRequired: "a required translate",
		},
		"B": errormapper.ErrorTranslator{
			validate.ErrRequired: "b required translate",
			validate.ErrMin:      "b min translate",
			validate.ErrMax:      "b max translate",
		},
	}

	tests := []struct {
		Description string
		Errors      validate.ErrorMap
		ExpectedOk  bool
		Expected    map[string]string
	}{
		{
			Description: "default",
			Errors: validate.ErrorMap{
				"A": validate.Errors{validate.ErrRequired, validate.ErrMin},
			},
			ExpectedOk: true,
			Expected: map[string]string{
				"A": "a required translate",
			},
		}, {
			Description: "2 fields",
			Errors: validate.ErrorMap{
				"A": validate.Errors{validate.ErrMin, validate.ErrRequired},
				"B": validate.Errors{validate.ErrMax, validate.ErrMin, validate.ErrMax},
			},
			ExpectedOk: true,
			Expected: map[string]string{
				"A": "a required translate",
				"B": "b max translate",
			},
		}, {
			Description: "1 field to translate no suitable translations found",
			Errors: validate.ErrorMap{
				"A": validate.Errors{validate.ErrMin},
			},
			ExpectedOk: false,
			Expected:   map[string]string{},
		}, {
			Description: "fails 1 field to translate no suitable translations found",
			Errors: validate.ErrorMap{
				"A": validate.Errors{validate.ErrMin},
			},
			ExpectedOk: false,
			Expected:   map[string]string{},
		}, {
			Description: "partial fails 1 field to translate no suitable translations found",
			Errors: validate.ErrorMap{
				"A": validate.Errors{validate.ErrMin},
				"B": validate.Errors{validate.ErrRequired},
			},
			ExpectedOk: false,
			Expected: map[string]string{
				"B": "b required translate",
			},
		}, {
			Description: "fails field has no translation set at all",
			Errors: validate.ErrorMap{
				"C": validate.Errors{validate.ErrMin},
			},
			ExpectedOk: false,
			Expected:   map[string]string{},
		},
	}

	for _, test := range tests {
		translated, ok := et.TranslateFirst(test.Errors)

		c.Assert(ok, Equals, test.ExpectedOk, Commentf(test.Description))
		c.Assert(translated, DeepEquals, test.Expected, Commentf(test.Description))
	}
}

func (s *FieldErrorTranslatorSuite) TestTranslateFirstWithDefault(c *C) {
	//a translation with a default (nil) translation always succeeds translation
	et := errormapper.FieldErrorTranslator{
		"A": errormapper.ErrorTranslator{
			validate.ErrRequired: "a required translate",
		},
		"B": errormapper.ErrorTranslator{
			validate.ErrRequired: "b min translate",
			validate.ErrMin:      "b min translate",
			validate.ErrMax:      "b max translate",
		},
		"": errormapper.ErrorTranslator{
			validate.ErrMin: "nil min default translate",
		},
	}

	tests := []struct {
		Description string
		Errors      validate.ErrorMap
		ExpectedOk  bool
		Expected    map[string]string
	}{
		{
			Description: "fallback min error on A to field default translations",
			Errors: validate.ErrorMap{
				"A": validate.Errors{validate.ErrRequired, validate.ErrMin},
			},
			ExpectedOk: true,
			Expected: map[string]string{
				"A": "a required translate",
			},
		}, {
			Description: "not set error field should fallback to default translations",
			Errors: validate.ErrorMap{
				"C": validate.Errors{validate.ErrMin},
			},
			ExpectedOk: true,
			Expected: map[string]string{
				"C": "nil min default translate",
			},
		}, {
			Description: "fallback no translation found should error out with a failed translation",
			Errors: validate.ErrorMap{
				"A": validate.Errors{validate.ErrRequired},
				"C": validate.Errors{validate.ErrRequired},
			},
			ExpectedOk: false,
			Expected: map[string]string{
				"A": "a required translate",
			},
		},
	}

	for _, test := range tests {
		translated, ok := et.TranslateFirst(test.Errors)

		c.Assert(ok, Equals, test.ExpectedOk, Commentf(test.Description))
		c.Assert(translated, DeepEquals, test.Expected, Commentf(test.Description))
	}
}

func (s *FieldErrorTranslatorSuite) TestTranslateFirstFallback(c *C) {

	fallback := errormapper.ErrorTranslator{
		validate.ErrMin:      "fallback err min",
		validate.ErrRequired: "fallback err required",
	}

	et := errormapper.FieldErrorTranslator{
		"A": errormapper.ErrorTranslator{
			validate.ErrMax: "a max translate",
		},
		"": errormapper.ErrorTranslator{
			validate.ErrRequired: "field map default err min",
		},
	}

	tests := []struct {
		Description string
		Errors      validate.ErrorMap
		Expected    map[string]string
	}{
		{
			Description: "fallback to provided fallaback",
			Errors: validate.ErrorMap{
				"A": validate.Errors{validate.ErrMin},
			},
			Expected: map[string]string{
				"A": "fallback err min",
			},
		}, {
			Description: "fallback used and has precedence over fieldmap default required",
			Errors: validate.ErrorMap{
				"A": validate.Errors{validate.ErrRequired},
			},
			Expected: map[string]string{
				"A": "fallback err required",
			},
		}, {
			Description: "fallback no field tranlation present at all",
			Errors: validate.ErrorMap{
				"B": validate.Errors{validate.ErrMin},
			},
			Expected: map[string]string{
				"B": "fallback err min",
			},
		}, {
			Description: "fallback no field translation present at all and fallback has precedence over fieldmap defaults",
			Errors: validate.ErrorMap{
				"B": validate.Errors{validate.ErrRequired},
			},
			Expected: map[string]string{
				"B": "fallback err required",
			},
		},
	}

	for _, test := range tests {
		translated, ok := et.TranslateFirst(test.Errors, fallback)

		c.Assert(ok, Equals, true, Commentf(test.Description))
		c.Assert(translated, DeepEquals, test.Expected, Commentf(test.Description))
	}
}
