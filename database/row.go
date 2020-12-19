package database

type Row struct {
	data map[string]interface{}
}

func (row *Row) GetData() map[string]interface{}{
	return row.data
}