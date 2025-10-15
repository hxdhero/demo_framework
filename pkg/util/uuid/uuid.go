package uuid

import (
	"github.com/google/uuid"
)

func UUID4() string {
	return uuid.New().String()
}

// ParseUUID 解析 UUID 字符串
func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// IsValidUUID 验证 UUID 字符串是否有效
func IsValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}
