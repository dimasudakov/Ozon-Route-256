package database

type Table struct {
	Name    string            `json:"name"`
	Columns map[string]string `json:"columns"`
}
