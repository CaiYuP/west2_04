package ids

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ParseIds(idsStr string) ([]int64, error) {
	// 获取查询参数
	if idsStr == "" {
		return nil, nil
	}

	// 将逗号分隔的字符串转换为 []int64
	ids, err := parseCommaSeparatedIDs(idsStr)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// 解析逗号分隔的ID字符串
func parseCommaSeparatedIDs(idsStr string) ([]int64, error) {
	var ids []int64

	parts := strings.Split(idsStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("ID格式错误: %s", part)
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return nil, errors.New("未提供有效的任务ID")
	}

	return ids, nil
}
