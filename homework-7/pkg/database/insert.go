package database

import (
	"errors"
	"fmt"
	"reflect"
)

type InsertBuilder struct {
	db     *Database
	table  string
	values map[string]interface{}
}

func Insert(connection *Database) InsertBuilder {
	return InsertBuilder{
		db: connection,
	}
}

func (ib InsertBuilder) Into(into string) InsertBuilder {
	ib.table = into
	return ib
}

func (ib InsertBuilder) Values(vals map[string]interface{}) InsertBuilder {
	ib.values = vals
	return ib
}

func (ib InsertBuilder) Execute() error {
	if !ib.db.ContainsTable(ib.table) {
		return fmt.Errorf("table with name: %s doesn't exist", ib.table)
	}
	if err := ib.validateValues(); err != nil {
		return err
	}
	if ok, err := ib.ContainRow(ib.values); ok || err != nil {
		if err != nil {
			return err
		} else {
			return fmt.Errorf("attempt to duplicate a record in a table")
		}
	}

	err := ib.db.insertRecordIntoTable(ib.table, ib.values)

	return err
}

func (ib InsertBuilder) validateValues() error {
	if !ib.db.ContainsTable(ib.table) {
		return fmt.Errorf("table with name: %s doesn't exist", ib.table)
	}
	for name, val := range ib.values {
		columnType, exist := ib.db.tables[ib.table].Columns[name]
		if !exist {
			return fmt.Errorf("column with name: %s doesn't exist in table: %s", name, ib.table)
		}
		valType := reflect.TypeOf(val).String()
		if columnType != valType {
			return fmt.Errorf("invalid type for column: %s, expected: %s, received: %s", name, columnType, valType)
		}
	}
	return nil
}

func (ib InsertBuilder) ContainRow(row map[string]interface{}) (bool, error) {
	var tables map[string][]interface{}
	if err := ib.db.readTablesData(&tables); err != nil {
		return false, err
	}

	for _, tableRow := range tables[ib.table] {
		var tableRowMap map[string]interface{}
		var ok bool
		tableRowMap, ok = tableRow.(map[string]interface{})
		if !ok {
			return false, errors.New("invalid row type")
		}

		checkResult := true
		if len(row) != len(tableRowMap) {
			checkResult = false
		}

		for key, value1 := range row {
			value2, ok := tableRowMap[key]
			if !ok {
				checkResult = false
			}
			kind := reflect.TypeOf(value1).Kind()
			if kind == reflect.Int {
				intVal1 := value1.(int)
				if float64(intVal1) != value2 {
					checkResult = false
					break
				}
			} else if kind == reflect.Float64 {
				if value1 != value2 {
					checkResult = false
					break
				}
			} else {
				if value1 != value2 {
					checkResult = false
					break
				}
			}
		}

		if checkResult {
			return true, nil
		}
	}

	return false, nil
}
