package database

// Holds the representation of a database row.
type Row struct {
	data map[string]interface{}
}

// Gets the Row data.
func (row *Row) Data() map[string]interface{}{
	return row.data
}