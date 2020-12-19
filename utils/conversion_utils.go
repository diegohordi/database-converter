package utils

import "database/sql"

/*
Converts a given value into a string.
*/
func ToString(value interface{}) string {
	if value, isString := value.(string); isString {
		return value
	}
	if value, isString := value.(*string); isString {
		return *value
	}
	if value, isString := value.(*sql.NullString); isString {
		return value.String
	}
	return ""
}

func GetRawValue(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		return v
	case *string:
		return *v
	case *sql.NullString:
		return v.String
	case *sql.NullInt64:
		return v.Int64
	case *sql.NullInt32:
		return v.Int32
	case *sql.NullBool:
		return v.Bool
	case *sql.NullFloat64:
		return v.Float64
	default:
		return value
	}
}
