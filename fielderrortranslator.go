package errortranslator

import (
	validate "github.com/mbict/go-validate"
)

// FieldErrorTranslator is a error translator that translates `errormaps` provided by the validator package.
type FieldErrorTranslator map[string]ErrorTranslator

// New creates a new Field error translator
func New() FieldErrorTranslator {
	return FieldErrorTranslator{}
}

// AddTranslation adds a new translation for error mapped to a field name.
// If a translation is already present is will be overwritten by the new translation
// The function returns a reference to the FieldErrorTranslator and is therefor very useful for chaining AddTranslation functions
// Equivalent to this function is: fielderrortranslator[field] = ErrorTranslator{ err: message, }
func (ft FieldErrorTranslator) AddTranslation(field string, err error, message string) FieldErrorTranslator {
	if _, ok := ft[field]; !ok {
		ft[field] = make(ErrorTranslator)
	}

	ft[field].AddTranslation(err, message)
	return ft
}

// SetDefaultTranslation sets a default translation for a field to use when no matching errors in the map is found.
// Equivalent to this function is: fielderrortranslator[field] = ErrorTranslator{ nil: message, }
func (ft FieldErrorTranslator) SetDefaultTranslation(field string, message string) FieldErrorTranslator {
	return ft.AddTranslation(field, nil, message)
}

// SetFallbackTranslation sets a default translation for a error when no field is found or error is found that matches
// error. This is useful to provide default translation for specific error messages
// Equivalent to this function is: fielderrortranslator[""] = ErrorTranslator{ err: message, }
func (ft FieldErrorTranslator) SetFallbackTranslation(err error, message string) FieldErrorTranslator {
	return ft.AddTranslation("", err, message)
}

// SetFallbackDefaultTranslation sets a default translation for when all other matching attempts fail.
// This is the last resort to return a translation. For example message `unknown occurred` or `there was a error`
// Equivalent to this function is: fielderrortranslator[""] = ErrorTranslator{ nil: message, }
func (ft FieldErrorTranslator) SetFallbackDefaultTranslation(message string) FieldErrorTranslator {
	return ft.AddTranslation("", nil, message)
}

// Translate will translate a map (validate.ErrorMap) with errors (validate.Errors) into a human readable
// message per field/map key.
// If any of the provided error fields fail to find a translation, the function will return the map with the translated
// errors and the second will be false indicated that we have a incomplete translation
func (ft FieldErrorTranslator) Translate(errorMap validate.ErrorMap, fallback ...ErrorTranslator) (map[string]string, bool) {
	return ft.translateErrorMap(errorMap, false, fallback)
}

// TranslateFirst works the same as Translate but will stop after the first positive match is found per field entry.
func (ft FieldErrorTranslator) TranslateFirst(errorMap validate.ErrorMap, fallback ...ErrorTranslator) (map[string]string, bool) {
	return ft.translateErrorMap(errorMap, true, fallback)
}

func (ft FieldErrorTranslator) translateErrorMap(errorMap validate.ErrorMap, firstOnly bool, fallback []ErrorTranslator) (map[string]string, bool) {

	//add default field translations as the last fallback
	translations, hasDefault := ft[""]
	if hasDefault {
		fallback = append(fallback, translations)
	}

	result := make(map[string]string)
	allTranslated := true
	for field, errs := range errorMap {
		errTrans, ok := ft[field]
		if !ok {
			if len(fallback) == 0 {
				allTranslated = false
				continue
			}
			errTrans = fallback[0]
		}

		var message string
		if firstOnly == true {
			message, ok = errTrans.TranslateFirst(errs, fallback...)
		} else {
			message, ok = errTrans.Translate(errs, fallback...)
		}

		allTranslated = allTranslated && ok
		if ok {
			result[field] = message
		}
	}
	return result, allTranslated
}
