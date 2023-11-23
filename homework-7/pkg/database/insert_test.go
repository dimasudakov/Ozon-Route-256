package database

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type InsertValues struct {
	tableName string
	values    []Vals
}

func TestInsert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		existingTables []Table
		insertValues   []InsertValues
		expectedError  error
	}{
		{
			name:           "Success",
			existingTables: []Table{testTable1},
			insertValues: []InsertValues{
				{
					tableName: testTable1.Name,
					values: []Vals{
						{"id": 1, "name": "Dima"},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:           "Success 2",
			existingTables: []Table{testTable1, testTable2},
			insertValues: []InsertValues{
				{
					tableName: testTable1.Name,
					values: []Vals{
						{"id": 1, "name": "Dima"},
						{"id": 2, "name": "Alex"},
						{"id": 3, "name": "Egor"},
						{"id": 4, "name": "Artem"},
						{"id": 101, "name": "Maxim"},
					},
				},
				{
					tableName: testTable2.Name,
					values: []Vals{
						{"id": 1, "name": "Item 1", "price": 10.5, "description": "Description 1"},
						{"id": 2, "name": "Item 2", "price": 20.3, "description": "Description 2"},
						{"id": 3, "name": "Item 3", "price": 15.7, "description": "Description 3"},
						{"id": 4, "name": "Item 4", "price": 30.2, "description": "Description 4"},
						{"id": 5, "name": "Item 5", "price": 25.9, "description": "Description 5"},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:           "Fail",
			existingTables: []Table{testTable1},
			insertValues: []InsertValues{
				{
					tableName: testTable1.Name,
					values: []Vals{
						{"id": 1, "name": "Dima"},
						{"id": 1, "name": "Dima"},
					},
				},
			},
			expectedError: errors.New("attempt to duplicate a record in a table"),
		},
		{
			name:           "Fail 2",
			existingTables: []Table{testTable1},
			insertValues: []InsertValues{
				{
					tableName: testTable1.Name,
					values: []Vals{
						{"id": 1, "name": "Dima", "surname": "Sudakov"},
					},
				},
			},
			expectedError: errors.New("column with name: surname doesn't exist in table: accounts"),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			NewDatabaseFixture()
			fixture.SetUp(t)
			defer fixture.TearDown(t)

			fixture.InsertSampleTables(t, tc.existingTables)

			var err error
			for _, tableValues := range tc.insertValues {
				for _, row := range tableValues.values {
					err = Insert(fixture.db).Into(tableValues.tableName).Values(row).Execute()
					if err != nil {
						break
					}
				}
			}

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				fixture.CheckSampleRowsInDatabase(t, tc.insertValues)
			}
		})
	}
}

func TestValidateValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		existingTables []Table
		tableName      string
		values         Vals
		expectedError  []string
	}{
		{
			name:           "Success",
			existingTables: []Table{testTable1},
			tableName:      "accounts",
			values:         Vals{"id": 1, "name": "Dima"},
			expectedError:  nil,
		},
		{
			name:           "Fail_invalid_type_1",
			existingTables: []Table{testTable1},
			tableName:      "accounts",
			values:         Vals{"id": "Dima", "name": 1},
			expectedError: []string{
				"invalid type for column: id, expected: int, received: string",
				"invalid type for column: name, expected: string, received: int",
			},
		},
		{
			name:           "Fail_invalid_type_2",
			existingTables: []Table{testTable1},
			tableName:      "accounts",
			values:         Vals{"id": uint32(1), "name": "Dima"},
			expectedError:  []string{"invalid type for column: id, expected: int, received: uint32"},
		},
		{
			name:           "Fail_invalid_table_name",
			existingTables: []Table{testTable2},
			tableName:      "users",
			values:         Vals{"id": 1, "name": "Dima"},
			expectedError:  []string{"table with name: users doesn't exist"},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			NewDatabaseFixture()
			fixture.SetUp(t)
			defer fixture.TearDown(t)

			for _, table := range tc.existingTables {
				err := fixture.db.CreateTable(ctx, table.Name, table.Columns)
				assert.NoError(t, err)
			}

			query := Insert(fixture.db).Into(tc.tableName).Values(tc.values)
			err := query.validateValues()
			if tc.expectedError != nil {
				assert.Contains(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
