package handler

import (
	pbcommon "west2-video/api/common/v1"
	v1 "west2-video/api/video/v1"
)

// HTTPResponse 用于 HTTP 网关层统一返回 JSON 结构
// 所有接口统一返回 {base, data} 格式
type HTTPResponse struct {
	Base *pbcommon.BaseResponse `json:"base"`
	Data interface{}            `json:"data"`
}

type videoListResp struct {
	Items []*v1.Video `json:"items"`
}
