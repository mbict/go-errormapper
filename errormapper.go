package errormapper

import (
	validate "github.com/mbict/go-validate"
)

type ErrorTranslator map[error]string

type FieldErrorTranslator map[string]ErrorTranslator

func (et ErrorTranslator) AddTranslation(err error, message string) ErrorTranslator {
	et[err] = message
	return et
}

func (et ErrorTranslator) SetDefaultTranslation(message string) ErrorTranslator {
	return et.AddTranslation(nil, message)
}

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

func (et ErrorTranslator) Translate(errs validate.Errors, fallback ...ErrorTranslator) (string, bool) {
	return et.translateErrors(errs, false, fallback)
}

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

func (ft FieldErrorTranslator) AddTranslation(field string, err error, message string) FieldErrorTranslator {
	if _, ok := ft[field]; !ok {
		ft[field] = make(ErrorTranslator)
	}

	ft[field].AddTranslation(err, message)
	return ft
}

func (ft FieldErrorTranslator) SetFieldDefaultTranslation(field string, message string) FieldErrorTranslator {
	return ft.AddTranslation(field, nil, message)
}

func (ft FieldErrorTranslator) SetFallbackTranslation(err error, message string) FieldErrorTranslator {
	return ft.AddTranslation("", err, message)
}

func (ft FieldErrorTranslator) SetFallbackDefaultTranslation(message string) FieldErrorTranslator {
	return ft.AddTranslation("", nil, message)
}

func (ft FieldErrorTranslator) Translate(errorMap validate.ErrorMap, fallback ...ErrorTranslator) (map[string]string, bool) {
	return ft.translateErrorMap(errorMap, false, fallback)
}

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
