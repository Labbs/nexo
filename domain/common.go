package domain

import (
	"database/sql/driver"
	"encoding/json"
)

// JSONB is a map of strings to interfaces
type JSONB map[string]any

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value any) error {
	return json.Unmarshal([]byte(value.(string)), j)
}

// JSONBArray is a slice that can be stored as JSONB in PostgreSQL
type JSONBArray []any

// Value implements the driver.Valuer interface
func (j JSONBArray) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

// Scan implements the sql.Scanner interface
func (j *JSONBArray) Scan(value any) error {
	if value == nil {
		*j = nil
		return nil
	}
	return json.Unmarshal([]byte(value.(string)), j)
}
