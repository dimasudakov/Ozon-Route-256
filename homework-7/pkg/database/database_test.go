package database

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestCreateTable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		tables        []Table
		expectedError error
	}{
		{
			name:          "Success",
			tables:        []Table{testTable1},
			expectedError: nil,
		},
		{
			name:          "Fail",
			tables:        []Table{testTable1, testTable1},
			expectedError: errors.New("table with name: accounts already exists"),
		},
		{
			name:          "Success 2",
			tables:        []Table{testTable1, testTable2},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			NewDatabaseFixture()

			fixture.SetUp(t)
			defer fixture.TearDown(t)

			var err error
			for _, table := range tc.tables {
				err = fixture.db.CreateTable(ctx, table.Name, table.Columns)
				if err != nil {
					break
				}
			}

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestReadTables(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		existingTables []Table
	}{
		{
			name:           "Success",
			existingTables: []Table{testTable1, testTable2},
		},
		{
			name:           "NoTables",
			existingTables: []Table{},
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

			err := fixture.db.readTables()
			assert.NoError(t, err)

			assert.Equal(t, len(tc.existingTables), len(fixture.db.tables))
			for _, table := range tc.existingTables {
				assert.True(t, fixture.db.ContainsTable(table.Name))
				assert.Equal(t, table, fixture.db.tables[table.Name])
			}
		})
	}
}

func TestNoDataRaces(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup

	var (
		selectValues1 = []string{"id", "name"}
		selectValues2 = []string{"id", "name", "price", "description"}
	)

	NewDatabaseFixture()
	fixture.SetUp(t)
	defer fixture.TearDown(t)

	err := fixture.db.CreateTable(ctx, testTable1.Name, testTable1.Columns)
	assert.NoError(t, err)
	err = fixture.db.CreateTable(ctx, testTable2.Name, testTable2.Columns)
	assert.NoError(t, err)

	N := 1000
	records := make([]map[string]interface{}, N)
	for i := 0; i < N; i++ {
		if i%2 == 0 {
			records[i] = fixture.GetRecordForTable(testTable1)
		} else {
			records[i] = fixture.GetRecordForTable(testTable2)
		}
	}

	wg.Add(N)
	for i := 0; i < N; i++ {
		go func(i int) {
			defer wg.Done()

			var err error
			if i%2 == 0 {
				err = Insert(fixture.db).Into(testTable1.Name).Values(records[i]).Execute()
			} else {
				err = Insert(fixture.db).Into(testTable2.Name).Values(records[i]).Execute()
			}
			assert.NoError(t, err)

			var result []map[string]interface{}
			if i%2 == 0 {
				result, err = Select(fixture.db).Values(selectValues1).From(testTable1.Name).Where(records[i]).Execute()
			} else {
				result, err = Select(fixture.db).Values(selectValues2).From(testTable2.Name).Where(records[i]).Execute()
			}
			assert.NoError(t, err)
			assert.Equal(t, 1, len(result))
		}(i)
	}

	wg.Wait()
	insertedRecords1, err := Select(fixture.db).Values(selectValues1).From(testTable1.Name).Execute()
	insertedRecords2, err := Select(fixture.db).Values(selectValues2).From(testTable2.Name).Execute()
	assert.Equal(t, N, len(insertedRecords1)+len(insertedRecords2))

	for i := 0; i < N; i++ {
		var result []map[string]interface{}
		var err error
		if i%2 == 0 {
			result, err = Select(fixture.db).Values(selectValues1).From(testTable1.Name).Where(records[i]).Execute()
		} else {
			result, err = Select(fixture.db).Values(selectValues2).From(testTable2.Name).Where(records[i]).Execute()
		}
		assert.NoError(t, err)
		assert.Equal(t, 1, len(result))
		assert.True(t, fixture.IsRecordsEqual(records[i], result[0]))
	}
}
