package errormapper

import validate "github.com/mbict/go-validate"

type Translation map[error]string

type ErrorMapper map[string]Translation

func (m ErrorMapper) AddMessage(field string, err error, message string) ErrorMapper {
	if _, ok := m[field]; !ok {
		m[field] = make(Translation)
	}

	m[field][err] = message
	return m
}

func (m ErrorMapper) Translate(errorMap validate.ErrorMap) map[string]string {

	first := false

	result := make(map[string]string)
	for field, errs := range errorMap {
		errMapper, ok := m[field]
		if !ok {
			continue
		}

		for _, err := range errs {

			//direct translation
			translation, ok := errMapper[err]
			if !ok {
				translation, ok = errMapper[nil]
				if !ok {
					//find next
					continue
				}
			}

			_, ok = result[field]
			if ok {
				result[field] = result[field] + ", " + translation
			} else {
				result[field] = translation
			}

			if first {
				break
			}
		}
	}
	return result
}

func (m ErrorMapper) TranslateFirst(errorMap validate.ErrorMap) map[string]string {
	result := make(map[string]string)
	for field, errs := range errorMap {
		errMapper, ok := m[field]
		if !ok {
			continue
		}

		if translation, ok := errMapper[errs[0]]; ok {
			//direct translation
			result[field] = translation
		} else if translation, ok := errMapper[nil]; ok {
			//default translation (nil error)
			result[field] = translation
		}
	}
	return result
}
