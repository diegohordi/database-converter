package database

// Holds the representation of a database table.
type Table struct {
	name string
	columns []*Column
}

// Gets the Table columns.
func (table *Table) Columns() []*Column {
	return table.columns
}

// Gets the Table name.
func (table *Table) Name() string {
	return table.name
}

// Gets the Table's primary key. Returns nil if no primary key has been found.
func (table *Table) GetPrimaryKey() *Column {
	for _, column := range table.Columns() {
		if column.key == PrimaryKey {
			return column
		}
	}
	return nil
}