package handler

import (
	"strconv"
	pbcommon "west2-video/api/common/v1"
	pbinteraction "west2-video/api/interaction/v1"
	v1 "west2-video/api/social/v1"
	pbvideo "west2-video/api/video/v1"
)

type HTTPResponseWithData struct {
	Base *pbcommon.BaseResponse `json:"base"`
	Data interface{}            `json:"data"`
}

type HTTPBaseResponse struct {
	Base *pbcommon.BaseResponse `json:"base"`
}

type VideoResponse struct {
	ID           string      `json:"id"`
	UserID       string      `json:"user_id"`
	Author       interface{} `json:"author,omitempty"` // 如果有 author 信息
	VideoURL     string      `json:"video_url"`
	CoverURL     string      `json:"cover_url"`
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	LikeCount    int64       `json:"like_count"`
	CommentCount int64       `json:"comment_count"`
	VisitCount   int64       `json:"visit_count"`
	CreatedAt    string      `json:"created_at"`
	UpdatedAt    string      `json:"updated_at"`
	DeletedAt    string      `json:"deleted_at"`
}

func convertVideoToResponse(v *pbvideo.Video) *VideoResponse {
	if v == nil {
		return nil
	}
	return &VideoResponse{
		ID:           v.GetId(),
		UserID:       v.GetUserId(),
		Author:       v.GetAuthor(),
		VideoURL:     v.GetVideoUrl(),
		CoverURL:     v.GetCoverUrl(),
		Title:        v.GetTitle(),
		Description:  v.GetDescription(),
		LikeCount:    v.GetLikeCount(),    // 即使为 0 也会输出
		CommentCount: v.GetCommentCount(), // 即使为 0 也会输出
		VisitCount:   v.GetVisitCount(),   // 即使为 0 也会输出
		CreatedAt:    v.GetCreatedAt(),
		UpdatedAt:    v.GetUpdatedAt(),
		DeletedAt:    v.GetDeletedAt(),
	}
}

func convertVideoListToResponse(videos []*pbvideo.Video) []*VideoResponse {
	if len(videos) == 0 {
		return make([]*VideoResponse, 0)
	}
	result := make([]*VideoResponse, 0, len(videos))
	for _, v := range videos {
		result = append(result, convertVideoToResponse(v))
	}
	return result
}

type videoListData struct {
	Items []*VideoResponse `json:"items"`
}

type videoListWithTotalData struct {
	Items []*VideoResponse `json:"items"`
	Total int64            `json:"total"`
}
type CommentResponse struct {
	ID         string `json:"id"`
	VideoID    string `json:"video_id"`
	UserID     string `json:"user_id"`
	ParentID   string `json:"parent_id"`
	LikeCount  int64  `json:"like_count"`
	ChildCount int64  `json:"child_count"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	DeletedAt  string `json:"deleted_at"`
}

func convertCommentToResponse(c *pbinteraction.Comment) *CommentResponse {
	if c == nil {
		return nil
	}

	return &CommentResponse{
		ID:         int64ToString(c.GetId()),
		VideoID:    int64ToString(c.GetVideoId()),
		UserID:     int64ToString(c.GetUserId()),
		ParentID:   int64ToString(c.GetParentId()),
		LikeCount:  c.GetLikeCount(),
		ChildCount: c.GetChildCount(),
		Content:    c.GetContent(),
		CreatedAt:  c.GetCreatedAt(),
		UpdatedAt:  c.GetUpdatedAt(),
		DeletedAt:  c.GetDeletedAt(),
	}
}

func convertCommentListToResponse(comments []*pbinteraction.Comment) []*CommentResponse {
	if len(comments) == 0 {
		return make([]*CommentResponse, 0)
	}
	result := make([]*CommentResponse, 0, len(comments))
	for _, c := range comments {
		result = append(result, convertCommentToResponse(c))
	}
	return result
}

// 辅助函数：将 int64 转换为 string
func int64ToString(v int64) string {
	return strconv.FormatInt(v, 10)
}

type commentListData struct {
	Items []*CommentResponse `json:"items"`
}

type socialListData struct {
	Items []*UserResponse `json:"items"`
	Total int64           `json:"total"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

func convertUserToResponse(u *v1.SocialUser) *UserResponse {
	if u == nil {
		return nil
	}
	return &UserResponse{
		ID:        u.GetId(),
		Username:  u.GetUsername(),
		AvatarURL: u.GetAvatarUrl(),
	}
}

func convertUserListToResponse(users []*v1.SocialUser) []*UserResponse {
	if len(users) == 0 {
		return make([]*UserResponse, 0)
	}
	result := make([]*UserResponse, 0, len(users))
	for _, u := range users {
		result = append(result, convertUserToResponse(u))
	}
	return result
}
