[![Build Status](https://travis-ci.org/mbict/go-errortranslator.png?branch=master)](https://travis-ci.org/mbict/go-errortranslator)
[![GoDoc](https://godoc.org/github.com/mbict/go-errortranslator?status.png)](http://godoc.org/github.com/mbict/go-errortranslator)
[![GoCover](http://gocover.io/_badge/github.com/mbict/go-errortranslator)](http://gocover.io/github.com/mbict/go-errortranslator)
[![GoReportCard](http://goreportcard.com/badge/mbict/go-errortranslator)](http://goreportcard.com/report/mbict/go-errortranslator)

Error Translator
================

ErrorTranslator is a package for translating errors that are provided by the go-validate package into human readable messages.
The human readable errors can be directly used in html forms or reporting validation errors in Rest apis. 


Examples
========
Simples example

#### Translation map dynamic creation
```go
translator := errortranslator.New()
translator.SetFallbackTranslation(validate.ErrRequired, "This is a required field")
translator.SetFallbackDefaultTranslation("There is a unknown error")
translator.AddTranslation("A", validate.ErrRequired, "A field is required")
translator.SetFieldDefaultTranslation("A", "A field has a error")

translatedMap, allTranslated := translator.Translate( validate.ErrorMap{
    "A": validate.Errors{ validate.ErrRequired, validate.ErrMin, },
    "B.1": validate.Errors{ validate.ErrRequired },
    "B.2": validate.Errors{ validate.ErrMin },
})


fmt.Println(allTranslated, translatedMap) 
```

#### Translation map by a plain map
```go
translator := errortranslator.FieldErrorTranslator{
    "A": errortranslator.ErrorTranslator{
        validate.ErrRequired: "A field is required",
        nil:                  "A field has a error",
    },
    "": errortranslator.ErrorTranslator{
        validate.ErrRequired: "This is a required field",
        nil:                  "There is a unknown error",
    },
}

translatedMap, allTranslated := translator.Translate(validate.ErrorMap{
    "A":   validate.Errors{validate.ErrRequired, validate.ErrMin},
    "B.1": validate.Errors{validate.ErrRequired},
    "B.2": validate.Errors{validate.ErrMin},
})

fmt.Println(allTranslated, translatedMap) 
```