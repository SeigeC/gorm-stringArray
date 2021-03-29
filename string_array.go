package stringArray

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type StringArray []string


func (arr *StringArray)Scan(value interface{})error{
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