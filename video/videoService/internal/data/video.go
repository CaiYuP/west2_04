package data

import (
	"strconv"
	"time"
	videoPb "west2-video/api/video/v1"
	"west2-video/common/tms"
)

// Video 表模型，对应表 `videos`
// 字段根据最新 API/数据库设计：
//   - user_id          视频作者
//   - video_url        视频文件链接
//   - cover_url        封面图
//   - title            标题
//   - description      描述
//   - visit_count      访问量
//   - like_count       点赞数
//   - comment_count    评论数
//   - created_at / updated_at / deleted_at 时间戳
//
// 注意：由于 gorm.Model 自带部分字段，这里手动声明以完全控制列名和默认值。

type Video struct {
	ID           uint64    `gorm:"column:id;primaryKey;autoIncrement;comment:视频ID" json:"id"`
	UserID       uint64    `gorm:"column:user_id;not null;comment:视频作者" json:"user_id"`
	VideoURL     string    `gorm:"column:video_url;type:varchar(512);not null;comment:视频链接" json:"video_url"`
	CoverURL     string    `gorm:"column:cover_url;type:varchar(512);not null;comment:封面链接" json:"cover_url"`
	Title        string    `gorm:"column:title;type:varchar(128);not null;comment:标题" json:"title"`
	Description  string    `gorm:"column:description;type:varchar(512);not null;comment:描述" json:"description"`
	VisitCount   uint64    `gorm:"column:visit_count;default:0;comment:访问量" json:"visit_count"`
	LikeCount    uint64    `gorm:"column:like_count;default:0;comment:点赞数" json:"like_count"`
	CommentCount uint64    `gorm:"column:comment_count;default:0;comment:评论数" json:"comment_count"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`
	DeletedAt    time.Time `gorm:"column:deleted_at;type:timestamp;comment:删除时间" json:"deleted_at"`
}

func (Video) TableName() string {
	return "videos"
}

func CopierVideos(videos []*Video) []*videoPb.Video {
	items := make([]*videoPb.Video, len(videos))
	for i, video := range videos {
		items[i] = &videoPb.Video{
			Id:           strconv.FormatUint(video.ID, 10),
			UserId:       strconv.FormatUint(video.UserID, 10),
			Title:        video.Title,
			Description:  video.Description,
			VisitCount:   int64(video.VisitCount),
			VideoUrl:     video.VideoURL,
			CoverUrl:     video.CoverURL,
			LikeCount:    int64(video.LikeCount),
			CommentCount: int64(video.CommentCount),
			CreatedAt:    tms.Format(video.CreatedAt),
			UpdatedAt:    tms.Format(video.UpdatedAt),
			DeletedAt:    tms.Format(video.DeletedAt),
		}
	}
	return items
}

func CopierVideo(video Video) *videoPb.Video {
	items := &videoPb.Video{
		Id:           strconv.FormatUint(video.ID, 10),
		UserId:       strconv.FormatUint(video.UserID, 10),
		Title:        video.Title,
		Description:  video.Description,
		VisitCount:   int64(video.VisitCount),
		VideoUrl:     video.VideoURL,
		CoverUrl:     video.CoverURL,
		CommentCount: int64(video.CommentCount),
		CreatedAt:    tms.Format(video.CreatedAt),
		UpdatedAt:    tms.Format(video.UpdatedAt),
		DeletedAt:    tms.Format(video.DeletedAt),
	}
	return items
}
