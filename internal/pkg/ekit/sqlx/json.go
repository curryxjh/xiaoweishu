package sqlx

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JsonColumn 代表存储字段的 JSON类型
// 主要用于没有提供默认 JSON 类型的数据库
// T可以是结构体，也可以是切片或者 map
// 理论上一切可以被 JSON 库处理的类型都能被用作 T
// 不建议使用指针作为 T 的类型
// 如果 T 是指针, 那么在 Val 为 nil 的情况下，一定要把 Valid 设置为 false
type JsonColumn[T any] struct {
	Val   T
	Valid bool
}

func (j JsonColumn[T]) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil
	}
	res, err := json.Marshal(j.Val)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (j *JsonColumn[T]) Scan(src any) error {
	var bs []byte
	switch val := src.(type) {
	case nil:
		return nil
	case []byte:
		bs = val
	case string:
		bs = []byte(val)
	default:
		return errors.New("invalid type")
	}
	if err := json.Unmarshal(bs, &j.Val); err != nil {
		return err
	}
	j.Valid = true
	return nil
}
