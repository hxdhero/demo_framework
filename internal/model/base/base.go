package base

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"lls_api/common"
	"time"
)

// NullTime 是一个 time.Time 的包装类型
// - 零值（IsZero）↔ 数据库 NULL
// - 非零值 ↔ 正常时间
type NullTime struct {
	time.Time
}

// Value 实现 driver.Valuer
// 如果是零值 → 返回 nil → 存为 NULL
func (nt NullTime) Value() (driver.Value, error) {
	if nt.IsZero() {
		return nil, nil
	}
	return nt.Time, nil
}

// Scan 实现 sql.Scanner
// 如果数据库是 NULL → 设置为零值
func (nt *NullTime) Scan(value interface{}) error {
	if value == nil {
		nt.Time = time.Time{} // 零值
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		nt.Time = v
	case []byte:
		t, err := time.Parse("2006-01-02 15:04:05.999999", string(v))
		if err != nil {
			return err
		}
		nt.Time = t
	case string:
		t, err := time.Parse("2006-01-02 15:04:05.999999", v)
		if err != nil {
			return err
		}
		nt.Time = t
	default:
		return fmt.Errorf("cannot scan %T into NullTime", value)
	}
	return nil
}

type ID = common.ID

func NullIDFromID(id ID) NullID {
	return NullID{NullInt32: sql.NullInt32{
		Int32: int32(id),
		Valid: true,
	}}
}

func NullIDFromInt(i int) NullID {
	return NullID{NullInt32: sql.NullInt32{
		Int32: int32(i),
		Valid: true,
	}}
}

// 自定义 NullID 类型
type NullID struct {
	sql.NullInt32
}

func (n *NullID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		return nil
	}
	n.Valid = true
	return json.Unmarshal(data, &n.Int32)
}

func (n NullID) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Int32)
}
