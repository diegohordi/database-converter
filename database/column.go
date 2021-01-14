package database

// Holds a table column representation
type Column struct {
	name string
	dataType string
	null bool
	key KeyType
	defaultValue interface{}
	extras string
}

// Gets the Column name.
func (col *Column) Name() string {
	return col.name
}

// Gets the Column data type.
func (col *Column) DataType() string {
	return col.dataType
}

// Determines if the Column can be null.
func (col *Column) Null() bool {
	return col.null
}

// Gets the Column key type.
func (col *Column) KeyType() KeyType {
	return col.key
}

// Gets the Column default value.
func (col *Column) DefaultValue() interface{} {
	return col.defaultValue
}

// Gets the Column extra information.
func (col *Column) GetExtras() string {
	return col.extras
}