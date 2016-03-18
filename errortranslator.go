package errortranslator

import (
	validate "github.com/mbict/go-validate"
)

// ErrorTranslator is a map that stores a human readable translations for errors.
type ErrorTranslator map[error]string

// AddTranslation adds a new translation to the map. It is stored in the map by the error object as key
// If a translation is already present is will be overwritten by the new translation
// The function returns a reference to the ErrorTranslator and is therefor very useful for chaining AddTranslation functions
// Equivalent to this function is: errortranslator[err] = message
func (et ErrorTranslator) AddTranslation(err error, message string) ErrorTranslator {
	et[err] = message
	return et
}

// SetDefaultTranslation sets a default translation to use when no matching errors in the map is found.
// Equivalent to this function is: errortranslator[nil] = message
func (et ErrorTranslator) SetDefaultTranslation(message string) ErrorTranslator {
	return et.AddTranslation(nil, message)
}

// TranslateError tries to translate the error into a human readable message.
// It will try to lookup the error in its map and returns the value as the message/translation,
// A fallback is used (if provided) when no match is found in the current map.
// When no match is found in the map or the fallback the default translation is returned (if set)
// When no match can be made at all the function will return a empty string and false as the succes flag
func (et ErrorTranslator) TranslateError(err error, fallback ...ErrorTranslator) (string, bool) {
	translation, ok := et[err]
	if !ok {
		translation, ok = et[nil]
		if !ok {
			//fallback to default
			if len(fallback) >= 1 {
				return fallback[0].TranslateError(err, fallback[1:]...)
			}
			return "", false
		}
	}
	return translation, true
}

// Translate will translate a slice of errors into a single human readable string.
// The validate.Errors is used from the validation package and is a slice with errors
func (et ErrorTranslator) Translate(errs validate.Errors, fallback ...ErrorTranslator) (string, bool) {
	return et.translateErrors(errs, false, fallback)
}

// TranslateFirst will only translate the first translatable error found in the map.
func (et ErrorTranslator) TranslateFirst(errs validate.Errors, fallback ...ErrorTranslator) (string, bool) {
	return et.translateErrors(errs, true, fallback)
}

func (et ErrorTranslator) translateErrors(errs validate.Errors, firstOnly bool, fallback []ErrorTranslator) (string, bool) {
	result := ""
	for _, err := range errs {
		translation, ok := et.TranslateError(err, fallback...)
		if !ok {
			continue
		}

		if result == "" {
			result = translation
		} else {
			result = result + ", " + translation
		}

		if firstOnly {
			//first message is enough head over to the next field
			return result, true
		}
	}
	return result, !(result == "")
}
