package database

import (
	"database-converter/config"
	"database-converter/errors"
	"database-converter/utils"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strings"
)

type MySQLImpl struct {
	connection *Connection
}

func (db *MySQLImpl) Connect(config config.DatabaseConfig) *errors.ApplicationError {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.User(), config.Password(), config.Host(), config.Port(), config.Database())
	conn, err := sql.Open("mysql", connString)
	if err != nil {
		return errors.WithMessageAndSourceErrorBuilder("Error establishing connection with database.", err).Build()
	}
	if err := conn.Ping(); err != nil {
		return errors.WithMessageAndSourceErrorBuilder("The database is unreachable.", err).Build()
	}
	db.connection = &Connection{config: config, conn: conn}
	return nil
}

func (db *MySQLImpl) Disconnect() *errors.ApplicationError {
	if instance, check := db.connection.conn.(*sql.DB); check {
		if err := instance.Close(); err != nil {
			return errors.WithMessageAndSourceErrorBuilder("Error closing connection.", err).Build()
		}
	}
	return nil
}

func (db *MySQLImpl) Describe(source string) (*Table, *errors.ApplicationError) {
	if connection, valid := db.connection.conn.(*sql.DB); valid {
		query := fmt.Sprintf("DESCRIBE %s;", source)
		if rows, err := connection.Query(query); err != nil {
			return nil, errors.WithMessageAndSourceErrorBuilder(fmt.Sprintf("Can't describe table %s.", source), err).Build()
		} else {
			table := Table{name: source}
			defer rows.Close()
			if columns, err := rows.ColumnTypes(); err != nil {
				return nil, errors.WithMessageAndSourceErrorBuilder(fmt.Sprintf("Can't fetch the table columns."), err).Build()
			} else {
				for rows.Next() {
					values := make([]interface{}, len(columns))
					object := map[string]interface{}{}
					for i, column := range columns {
						object[strings.ToLower(column.Name())] = db.GetInterface(*column)
						values[i] = object[strings.ToLower(column.Name())]
					}
					err = rows.Scan(values...)
					if err != nil {
						return nil, errors.WithMessageAndSourceErrorBuilder(fmt.Sprintf("Can't describe table %s.", source), err).Build()
					}
					colDefinition := &Column{}
					colDefinition.name = utils.ToString(object["field"])
					colDefinition.dataType = utils.ToString(object["type"])
					colDefinition.defaultValue = object["default"]
					colDefinition.extras = utils.ToString(object["extra"])
					colDefinition.null = utils.ToString(object["null"]) != "NO"
					switch strings.ToLower(utils.ToString(object["key"])) {
					case "pri":
						colDefinition.key = PrimaryKey
					}
					table.columns = append(table.columns, colDefinition)
				}
			}
			return &table, nil
		}
	}
	return nil, errors.WithMessageBuilder(fmt.Sprintf("Can't describe %s.", source)).Build()
}

func (db *MySQLImpl) GetInterface(columnType sql.ColumnType) interface{} {
	switch strings.ToLower(columnType.DatabaseTypeName()) {
	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext", "enum":
		return new(sql.NullString)
	case "smallint", "int", "tinyint", "integer", "bigint", "year":
		return new(sql.NullInt64)
	case "bool", "boolean":
		return new(sql.NullBool)
	case "float", "double", "decimal", "dec":
		return new(sql.NullFloat64)
	case "date", "datetime", "timestamp":
		return new(sql.NullString)
	default:
		return reflect.New(columnType.ScanType()).Interface()
	}
}

func (db *MySQLImpl) Count(table Table) (int, *errors.ApplicationError) {
	if connection, valid := db.connection.conn.(*sql.DB); valid {
		var count int
		query := fmt.Sprintf("SELECT COUNT(%s) FROM %s", table.GetPrimaryKey().name, table.Name())
		row := connection.QueryRow(query)
		err := row.Scan(&count)
		if err != nil {
			return 0, errors.WithMessageAndSourceErrorBuilder(fmt.Sprintf("Can't count rows."), err).Build()
		}
		return count, nil
	}
	return 0, nil
}

func (db *MySQLImpl) GetRows(table Table, columns []string, rowChannel chan interface{}) {
	if connection, valid := db.connection.conn.(*sql.DB); valid {
		query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(columns[:], ", "), table.Name())
		if rows, err := connection.Query(query); err != nil {
			rowChannel <- errors.WithMessageAndSourceErrorBuilder(fmt.Sprintf("Error performing the query."), err).Build()
		} else {
			defer rows.Close()
			if columns, err := rows.ColumnTypes(); err != nil {
				rowChannel <- errors.WithMessageAndSourceErrorBuilder(fmt.Sprintf("Can't fetch the table columns."), err).Build()
			} else {
				for rows.Next() {
					values := make([]interface{}, len(columns))
					object := map[string]interface{}{}
					for i, column := range columns {
						object[strings.ToLower(column.Name())] = db.GetInterface(*column)
						values[i] = object[strings.ToLower(column.Name())]
					}
					err = rows.Scan(values...)
					if err != nil {
						rowChannel <- errors.WithMessageAndSourceErrorBuilder(fmt.Sprintf("Can't fetch row."), err).Build()
					} else {
						for _, column := range columns {
							object[strings.ToLower(column.Name())] = utils.GetRawValue(object[strings.ToLower(column.Name())])
						}
						rowChannel <- &Row{data: object}
					}
				}
			}
		}
	}
}

func (db *MySQLImpl) Insert(table Table, row Row) *errors.ApplicationError {
	panic("implement me")
}
