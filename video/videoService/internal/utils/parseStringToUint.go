package utils

import (
	"strconv"
	"strings"
)

func ExtractVideoIDsFromKeys(keys []string) ([]uint64, error) {
	ids := make([]uint64, 0, len(keys))
	prefix := "video:visit:"

	for _, key := range keys {
		// 移除前缀 "video:visit:"
		idStr := strings.TrimPrefix(key, prefix)

		// 转换为 uint64
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			// 如果转换失败，跳过这个 key（或记录错误）
			continue
		}
		ids = append(ids, id)
	}

	return ids, nil
}
