package database

type Table struct {
	name string
	columns []*Column
}

func (table *Table) GetColumns() []*Column {
	return table.columns
}

func (table *Table) GetName() string {
	return table.name
}

func (table *Table) GetPrimaryKey() *Column {
	for _, column := range table.GetColumns() {
		if column.key == PrimaryKey {
			return column
		}
	}
	return nil
}