package sqlx

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JsonColumn 核心目的是在 go 程序与数据库(MySQL PostgreSQL SQLite)之间建立一个自动化的 JSON 转换机制
// 日常开发中, 我们经常需要把复杂的对象存储到数据库的一个字段中
// 手动处理时, 一般情况下我们需要在存储之前调用 json.Marshal(), 在读取之后调用 json.Unmarshal()
// 数据库驱动只认识基本数据类型(string, int, []byte)

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
		return fmt.Errorf("ekit: json column can only scan string or []byte, but got %T", val)
	}
	if err := json.Unmarshal(bs, &j.Val); err != nil {
		return err
	}
	j.Valid = true
	return nil
}
