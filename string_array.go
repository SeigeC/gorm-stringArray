package stringArray

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	
	"gorm.io/gorm/clause"
)

type StringArray []string

func (arr *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	err := json.Unmarshal(bytes, &arr)
	return err
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (arr StringArray) Value() (driver.Value, error) {
	if len(arr) == 0 {
		return nil, nil
	}
	return json.Marshal(arr)
}

type funcType string

const (
	hasKey funcType = "hasKey"
)

type expression struct {
	column   string
	keys     []string
	key      string
	funcType funcType
}

func NewExpression(column string) *expression {
	return &expression{column: column}
}

func (e *expression) HasKey(key string) *expression {
	e.key = key
	e.funcType = hasKey
	return e
}

func (e *expression) Build(builder clause.Builder) {
	if stmt, ok := builder.(*gorm.Statement); ok {
		switch stmt.Dialector.Name() {
		case "mysql", "sqlite":
			switch e.funcType{
			case hasKey:
				builder.WriteString(fmt.Sprintf(`JSON_CONTAINS(%s->'$[*]', '"%s"', '$')`, stmt.Quote(e.column), e.key))
			}
		}
	}
}
