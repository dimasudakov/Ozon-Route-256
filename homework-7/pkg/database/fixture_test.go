package database

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"
)

var (
	db      *Database
	fixture *DatabaseFixture
	once    sync.Once
)

var (
	ctx        = context.Background()
	testTable1 = Table{
		Name: "accounts",
		Columns: map[string]string{
			"id":   "int",
			"name": "string",
		},
	}
	testTable2 = Table{
		Name: "subscriptions",
		Columns: map[string]string{
			"id":          "int",
			"name":        "string",
			"price":       "float64",
			"description": "string",
		},
	}
)

type DatabaseFixture struct {
	db *Database
	sync.Mutex
	rnd *rand.Rand
}

func NewDatabaseFixture() {
	once.Do(func() {
		var err error
		db, err = GetConnection()
		if err != nil {
			panic("can't connect to database")
		}
		source := rand.NewSource(time.Now().UnixNano())
		rnd := rand.New(source)
		fixture = &DatabaseFixture{
			db:  db,
			rnd: rnd,
		}
	})
}

func (f *DatabaseFixture) SetUp(t *testing.T) {
	f.Lock()
	if err := f.db.DropDatabase(); err != nil {
		panic("can't set up database")
	}
}

func (f *DatabaseFixture) TearDown(t *testing.T) {
	if err := f.db.DropDatabase(); err != nil {
		panic("can't set up database")
	}
	f.Unlock()
}

func (f *DatabaseFixture) InsertSampleTables(t *testing.T, tables []Table) {
	for _, table := range tables {
		err := fixture.db.CreateTable(ctx, table.Name, table.Columns)
		assert.NoError(t, err)
	}
}

func (f *DatabaseFixture) InsertSampleTableData(t *testing.T, records []InsertValues) {
	for _, tableValues := range records {
		for _, row := range tableValues.values {
			err := Insert(fixture.db).Into(tableValues.tableName).Values(row).Execute()
			assert.NoError(t, err)
		}
	}
}

func (f *DatabaseFixture) CheckSampleRowsInDatabase(t *testing.T, sampleRows []InsertValues) {
	var tables map[string][]interface{}
	err := f.db.readTablesData(&tables)
	assert.NoError(t, err)

	for _, tableRows := range sampleRows {
		for _, row := range tableRows.values {
			err := Insert(f.db).Into(tableRows.tableName).Values(row).Execute()
			assert.EqualError(t, err, "attempt to duplicate a record in a table")
		}
	}
}

func (f *DatabaseFixture) CheckRows(t *testing.T, expected []map[string]interface{}, actual []map[string]interface{}) {
	assert.Equal(t, len(expected), len(actual))
	for _, expRow := range expected {
		var checkResult bool
		for _, actRow := range actual {
			if f.IsRecordsEqual(expRow, actRow) {
				checkResult = true
				break
			}
		}
		assert.True(t, checkResult, fmt.Sprintf("no expected row: %v in result", expRow))
	}
}

func (f *DatabaseFixture) IsRecordsEqual(expRow, actRow map[string]interface{}) bool {
	if len(expRow) != len(actRow) {
		return false
	}

	for key, value1 := range expRow {
		value2, ok := actRow[key]
		if !ok {
			return false
		}
		kind1 := reflect.TypeOf(value1).Kind()
		kind2 := reflect.TypeOf(value2).Kind()
		if kind1 == reflect.Int {
			if kind2 == reflect.Float64 {
				intVal1 := value1.(int)
				if float64(intVal1) != value2 {
					return false
				}
			} else {
				if value1 != value2 {
					return false
				}
			}
		} else if kind1 == reflect.Float64 {
			if kind2 == reflect.Int {
				intVal2 := value2.(int)
				if float64(intVal2) != value1 {
					return false
				}
			} else {
				if value1 != value2 {
					return false
				}
			}
		} else {
			if value1 != value2 {
				return false
			}
		}
	}
	return true
}

func (f *DatabaseFixture) GetRecordForTable(table Table) map[string]interface{} {
	record := make(map[string]interface{})

	for columnName, columnType := range table.Columns {
		switch columnType {
		case "int":
			record[columnName] = f.rnd.Intn(100000)
		case "float64":
			record[columnName] = f.rnd.Float64() * 100000
		case "string":
			record[columnName] = f.generateRandomString(10)
		}
	}

	return record
}

func (f *DatabaseFixture) generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[f.rnd.Intn(len(charset))]
	}
	return string(result)
}
