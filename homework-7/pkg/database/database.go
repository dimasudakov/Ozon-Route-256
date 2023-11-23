package database

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var (
	instance *Database
	onceDB   sync.Once
)

const (
	TABLES_FILE = "data/tables.json"
	DATA_FILE   = "data/data.json"
)

type Vals map[string]interface{}

type Database struct {
	tables map[string]Table
	sync.RWMutex
}

func GetConnection() (*Database, error) {
	var err error
	onceDB.Do(func() {
		db := &Database{}
		if err = db.readTables(); err != nil {
			onceDB = sync.Once{}
			return
		}
		instance = db
	})
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func (d *Database) CreateTable(ctx context.Context, name string, columns map[string]string) error {
	if _, exists := d.tables[name]; exists {
		return fmt.Errorf("table with name: %s already exists", name)
	}

	newTable := Table{
		Name:    name,
		Columns: columns,
	}

	d.tables[newTable.Name] = newTable

	err := d.saveTables()
	return err
}

func (d *Database) ContainsTable(name string) bool {
	_, exist := d.tables[name]
	return exist
}

func (d *Database) DropDatabase() error {
	d.Lock()
	defer d.Unlock()

	if err := d.clearFile(TABLES_FILE); err != nil {
		return err
	}
	if err := d.clearFile(DATA_FILE); err != nil {
		return err
	}

	d.tables = make(map[string]Table)

	return nil
}

func (d *Database) DropTable(table string) error {
	if !d.ContainsTable(table) {
		return fmt.Errorf("table with name: %s doesn't exist", table)
	}
	d.Lock()
	defer d.Unlock()

	delete(d.tables, table)

	err := d.saveTables()
	return err
}

func (d *Database) readTables() error {
	d.Lock()
	defer d.Unlock()

	data, err := os.ReadFile(TABLES_FILE)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		d.tables = make(map[string]Table)
		return nil
	}

	if err := json.Unmarshal(data, &d.tables); err != nil {
		return err
	}

	return nil
}

func (d *Database) saveTables() error {
	tablesJson, err := json.Marshal(d.tables)
	if err != nil {
		return err
	}
	err = os.WriteFile(TABLES_FILE, tablesJson, 0666)
	return err
}

func (d *Database) insertRecordIntoTable(tableName string, record map[string]interface{}) error {
	d.Lock()
	defer d.Unlock()

	allRecords := make(map[string][]interface{})

	dataJson, err := os.ReadFile(DATA_FILE)
	if err != nil {
		return err
	}
	if len(dataJson) != 0 {
		if err := json.Unmarshal(dataJson, &allRecords); err != nil {
			return err
		}
	}

	allRecords[tableName] = append(allRecords[tableName], record)

	if err := d.updateTableData(tableName, &allRecords); err != nil {
		return err
	}

	return nil
}

func (d *Database) readTablesData(dataContainer *map[string][]interface{}) error {
	d.RLock()
	dataJson, err := os.ReadFile(DATA_FILE)
	defer d.RUnlock()

	if err != nil {
		return err
	}
	if len(dataJson) == 0 {
		*dataContainer = make(map[string][]interface{})
	} else {
		if err := json.Unmarshal(dataJson, dataContainer); err != nil {
			return err
		}
	}

	return nil
}

func (d *Database) updateTableData(table string, dataContainer *map[string][]interface{}) error {
	if !d.ContainsTable(table) {
		return fmt.Errorf("table with name: %s doesn't exist", table)
	}

	dataJson, err := json.Marshal(dataContainer)
	if err != nil {
		return err
	}

	err = os.WriteFile(DATA_FILE, dataJson, 0666)
	return err
}

func (d *Database) clearFile(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	if err := file.Truncate(0); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return nil
}
