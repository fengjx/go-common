package utils

import (
	"encoding/json"
	"fmt"
)

func ToString(src interface{}) string {
	if src == nil {
		return ""
	}

	switch v := src.(type) {
	case string:
		return src.(string)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", src)
	case float32, float64:
		bs, _ := json.Marshal(v)
		return string(bs)
	case bool:
		if b, ok := src.(bool); ok && b {
			return "true"
		} else {
			return "false"
		}
	default:
		return fmt.Sprintf("%v", v)
	}
}