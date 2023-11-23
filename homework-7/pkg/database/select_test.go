package database

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		existingTables     []Table
		existingTablesData []InsertValues
		tableName          string
		columns            []string
		conditions         map[string]interface{}
		expected           []map[string]interface{}
		expectedError      error
	}{
		{
			name:           "Success",
			existingTables: []Table{testTable1},
			existingTablesData: []InsertValues{
				{
					tableName: testTable1.Name,
					values: []Vals{
						{"id": 1, "name": "Dima"},
					},
				},
			},
			tableName:  testTable1.Name,
			columns:    []string{"id", "name"},
			conditions: nil,
			expected: []map[string]interface{}{
				{"id": 1, "name": "Dima"},
			},
			expectedError: nil,
		},
		{
			name:           "Success 2",
			existingTables: []Table{testTable2},
			existingTablesData: []InsertValues{
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
			tableName:  testTable2.Name,
			columns:    []string{"id", "price"},
			conditions: nil,
			expected: []map[string]interface{}{
				{"id": 1, "price": 10.5},
				{"id": 2, "price": 20.3},
				{"id": 3, "price": 15.7},
				{"id": 4, "price": 30.2},
				{"id": 5, "price": 25.9},
			},
		},
		{
			name:           "Success 3",
			existingTables: []Table{testTable2},
			existingTablesData: []InsertValues{
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
			tableName: testTable2.Name,
			columns:   []string{"id", "price", "description"},
			conditions: map[string]interface{}{
				"price": 15.7,
			},
			expected: []map[string]interface{}{
				{"id": 3, "price": 15.7, "description": "Description 3"},
			},
		},
		{
			name:               "Fail (No column in table)",
			existingTables:     []Table{testTable1},
			existingTablesData: []InsertValues{},
			tableName:          testTable1.Name,
			columns:            []string{"id", "name", "surname"},
			conditions:         nil,
			expected:           nil,
			expectedError:      errors.New("columns with name: surname doesn't exist in table: accounts"),
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
			fixture.InsertSampleTableData(t, tc.existingTablesData)

			result, err := Select(fixture.db).Values(tc.columns).From(tc.tableName).Where(tc.conditions).Execute()

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				fixture.CheckRows(t, tc.expected, result)
			}
		})
	}
}
