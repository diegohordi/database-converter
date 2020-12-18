package utils

func ToString(value interface{}) string {
	if value, isString := value.(string); isString {
		return value
	}
	return ""
}
