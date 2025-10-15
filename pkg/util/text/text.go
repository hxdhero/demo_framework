package text

import (
	"strings"
	"unicode"
)

/**
文本相关操作
*/

// ToSnakeCase 转换为 snake_case
func ToSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
