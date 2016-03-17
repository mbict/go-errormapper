package errormapper

import validate "github.com/mbict/go-validate"

type ErrorTranslations map[error]string

type FieldTranslations map[string]ErrorTranslations

func (et ErrorTranslations) AddTranslation(err error, message string) ErrorTranslations {
	et[err] = message
	return et
}

func (et ErrorTranslations) AddDefaultTranslation(message string) ErrorTranslations {
	return et.AddTranslation(nil, message)
}

func (et ErrorTranslations) TranslateError(err error, fallback ...ErrorTranslations) (string, bool) {
	translation, ok := et[err]
	if !ok {
		translation, ok = et[nil]
		if !ok {
			//fallback to default
			if len(fallback) >= 1 {
				return fallback[0].TranslateError(err, fallback[1:]...)
			}
		}
		return "", false
	}
	return translation, true
}

func (et ErrorTranslations) Translate(errs validate.Errors, fallback ...ErrorTranslations) string {
	return et.translateErrors(errs, false, fallback)
}

func (et ErrorTranslations) TranslateFirst(errs validate.Errors, fallback ...ErrorTranslations) string {
	return et.translateErrors(errs, true, fallback)
}

func (et ErrorTranslations) translateErrors(errs validate.Errors, firstOnly bool, fallback []ErrorTranslations) string {
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
			return result
		}
	}
	return result
}

func (ft FieldTranslations) AddTranslation(field string, err error, message string) FieldTranslations {
	if _, ok := ft[field]; !ok {
		ft[field] = make(ErrorTranslations)
	}

	ft[field].AddTranslation(err, message)
	return ft
}

func (ft FieldTranslations) AddDefaultTranslation(field string, message string) FieldTranslations {
	return ft.AddTranslation(field, nil, message)
}

func (ft FieldTranslations) Translate(errorMap validate.ErrorMap, fallback ...ErrorTranslations) map[string]string {
	return ft.translateErrorMap(errorMap, false, fallback)
}

func (ft FieldTranslations) TranslateFirst(errorMap validate.ErrorMap, fallback ...ErrorTranslations) map[string]string {
	return ft.translateErrorMap(errorMap, true, fallback)
}

func (ft FieldTranslations) translateErrorMap(errorMap validate.ErrorMap, firstOnly bool, fallback []ErrorTranslations) map[string]string {
	result := make(map[string]string)
	for field, errs := range errorMap {
		errTrans, ok := ft[field]
		if !ok {
			continue
		}

		var message string
		if firstOnly == true {
			message = errTrans.TranslateFirst(errs, fallback...)
		} else {
			message = errTrans.Translate(errs, fallback...)
		}

		if message != "" {
			result[field] = message
		}
	}
	return result
}
