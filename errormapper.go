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

func (et ErrorTranslations) Translate(errs validate.Errors) string {
	return et.translateErrors(errs, false)
}

func (et ErrorTranslations) TranslateFirst(errs validate.Errors) string {
	return et.translateErrors(errs, true)
}

func (et ErrorTranslations) translateErrors(errs validate.Errors, firstOnly bool) string {
	result := ""
	for _, err := range errs {

		//direct translation
		translation, ok := et[err]
		if !ok {
			translation, ok = et[nil]
			if !ok {
				//find next
				continue
			}
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

func (ft FieldTranslations) Translate(errorMap validate.ErrorMap) map[string]string {
	return ft.translateErrorMap(errorMap, false)
}

func (ft FieldTranslations) TranslateFirst(errorMap validate.ErrorMap) map[string]string {
	return ft.translateErrorMap(errorMap, true)
}

func (ft FieldTranslations) translateErrorMap(errorMap validate.ErrorMap, firstOnly bool) map[string]string {
	result := make(map[string]string)
	for field, errs := range errorMap {
		errTrans, ok := ft[field]
		if !ok {
			continue
		}

		var message string
		if firstOnly == true {
			message = errTrans.TranslateFirst(errs)
		} else {
			message = errTrans.Translate(errs)
		}

		if message != "" {
			result[field] = message
		}
	}
	return result
}
