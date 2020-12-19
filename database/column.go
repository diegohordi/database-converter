package database

type KeyType string

const (
	PrimaryKey KeyType = "PRIMARY"
)

type Column struct {
	name string
	dataType string
	null bool
	key KeyType
	defaultValue interface{}
	extras string
}

func (col *Column) GetName() string {
	return col.name
}

func (col *Column) GetDataType() string {
	return col.dataType
}

func (col *Column) IsNull() bool {
	return col.null
}

func (col *Column) GetKey() KeyType {
	return col.key
}

func (col *Column) GetDefaultValue() interface{} {
	return col.defaultValue
}

func (col *Column) GetExtras() string {
	return col.extras
}