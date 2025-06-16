package gormx

import (
	"database/sql/driver"
	"encoding/json"
)

type GormList []string

func (g *GormList) Scan(src any) error {
	return json.Unmarshal(src.([]byte), &g)
}

func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}