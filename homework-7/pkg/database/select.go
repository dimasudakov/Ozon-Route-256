package database

import (
	"fmt"
	"reflect"
	"slices"
)

type SelectBuilder struct {
	db         *Database
	table      string
	values     []string
	conditions map[string]interface{}
}

func Select(db *Database) SelectBuilder {
	return SelectBuilder{
		db: db,
	}
}

func (sb SelectBuilder) Values(vals []string) SelectBuilder {
	sb.values = append(sb.values, vals...)
	return sb
}

func (sb SelectBuilder) From(from string) SelectBuilder {
	sb.table = from
	return sb
}

func (sb SelectBuilder) Where(conditions map[string]interface{}) SelectBuilder {
	sb.conditions = conditions
	return sb
}

func (sb SelectBuilder) Execute() ([]map[string]interface{}, error) {
	if !sb.db.ContainsTable(sb.table) {
		return nil, fmt.Errorf("table with name: %s doesn't exist", sb.table)
	}
	if err := sb.validateSelectValues(); err != nil {
		return nil, err
	}
	result, err := sb.selectRows()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (sb SelectBuilder) validateSelectValues() error {
	for _, valName := range sb.values {
		if _, exist := sb.db.tables[sb.table].Columns[valName]; !exist {
			return fmt.Errorf("columns with name: %s doesn't exist in table: %s", valName, sb.table)
		}
	}
	return nil
}

func (sb SelectBuilder) selectRows() ([]map[string]interface{}, error) {
	var data map[string][]interface{}
	if err := sb.db.readTablesData(&data); err != nil {
		return nil, err
	}
	if _, ok := data[sb.table]; !ok {
		return nil, nil
	}

	result := make([]map[string]interface{}, 0)
	for _, tableRow := range data[sb.table] {
		row := tableRow.(map[string]interface{})

		if !sb.checkRowConditions(row) {
			continue
		}

		rowSelectedColumns := make(map[string]interface{})
		for columnName, value := range row {
			if slices.Contains(sb.values, columnName) {
				rowSelectedColumns[columnName] = value
			}
		}

		result = append(result, rowSelectedColumns)
	}

	return result, nil
}

func (sb SelectBuilder) checkRowConditions(row map[string]interface{}) bool {
	for columnName, expectedValue := range sb.conditions {
		rowColumnValue, ok := row[columnName]
		if !ok {
			return false
		}

		kind1 := reflect.TypeOf(rowColumnValue).Kind()
		kind2 := reflect.TypeOf(expectedValue).Kind()
		if kind1 == reflect.Int {
			if kind2 == reflect.Float64 {
				rowValueInt := rowColumnValue.(int)
				if float64(rowValueInt) != expectedValue {
					return false
				}
			} else {
				if expectedValue != rowColumnValue {
					return false
				}
			}
		} else if kind1 == reflect.Float64 {
			if kind2 == reflect.Int {
				expValueInt := expectedValue.(int)
				if float64(expValueInt) != rowColumnValue {
					return false
				}
			} else {
				if expectedValue != rowColumnValue {
					return false
				}
			}
		} else {
			if rowColumnValue != expectedValue {
				return false
			}
		}
	}

	return true
}
