package utils

/*
Converts a given value into a string.
*/
func ToString(value interface{}) string {
	if value, isString := value.(string); isString {
		return value
	}
	return ""
}
