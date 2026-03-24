package domain

import (
	"bytes"
	"context"
	"strconv"
	"time"
	"videoService/config"
	"videoService/internal/dao"
	"videoService/internal/data"
	"videoService/internal/repo"
	"videoService/internal/utils"
	"west2-video/common/errs"
	"west2-video/common/logs"
	"west2-video/common/minio"
	"west2-video/gateway/biz/model"

	"go.uber.org/zap"

	"github.com/go-redis/redis/v8"
)

func NewVideoDomain() *VideoDomain {
	return &VideoDomain{
		videoRepo: dao.NewVideoDao(),
		cache:     dao.Rc,
	}
}

type VideoDomain struct {
	videoRepo repo.VideoRepo
	cache     repo.Cache
}

func (d *VideoDomain) VerifySeed(ctx context.Context, cur int64) *errs.BError {
	t := time.Now().UnixMilli()
	if t < cur {
		return model.SeedError
	}
	return nil
}

func (d *VideoDomain) FindVideosAfterTime(ctx context.Context, latestTime int64) ([]*data.Video, *errs.BError) {
	items, err := d.videoRepo.FindVideosAfterTime(ctx, latestTime)
	if err != nil {
		logs.LG.Error("VideoDomain.FindVideosAfterTime.FindVideosAfterTime error", zap.Error(err))
		return nil, model.DBError
	}
	if items == nil {
		return make([]*data.Video, 0), nil
	}
	return items, nil
}

func (d *VideoDomain) IncrementVisitCount(ctx context.Context, videoID uint64) *errs.BError {
	// 调用 dao.Rc.IncrVisitCount(videoID)
	// 只更新 Redis，不更新数据库
	_, err := d.cache.IncrVisitCount(ctx, videoID)
	if err != nil {
		logs.LG.Error("VideoDomain.IncrementVisitCount error", zap.Error(err))
		return model.RedisError
	}
	return nil
}
func (d *VideoDomain) GetHotRanking(ctx context.Context, pageSize, pageNum int) ([]*data.Video, int64, *errs.BError) {
	// 1. 从 Redis ZSET 查询（ZRevRange）
	start := (pageNum - 1) * pageSize
	stop := start + pageSize
	ids, err := d.cache.ZRevRange(ctx, dao.VideoRankingKey, int64(start), int64(stop))
	if err != nil {
		logs.LG.Error("VideoDomain.GetHotRanking error", zap.Error(err))
		return nil, 0, model.RedisError
	}
	// 2. 如果 Redis 为空，从数据库查询并写入 Redis
	if ids == nil || len(ids) < pageSize {
		items, err := d.videoRepo.FindHotRankingVideos(ctx, pageSize, pageNum)
		if err != nil {
			logs.LG.Error("VideoDomain.FindHotRankingVideos error", zap.Error(err))
			return nil, 0, model.DBError
		}
		return items, int64(len(items)), nil
	}
	// 3. 根据 videoIDs 查询视频详情
	idsInt := make([]uint64, len(ids))
	for i, id := range ids {
		idsInt[i], _ = strconv.ParseUint(id, 10, 64)
	}
	items, err := d.videoRepo.FindVideosByIDsDesc(ctx, idsInt)
	if err != nil {
		logs.LG.Error("VideoDomain.FindVideosByIDs error", zap.Error(err))
		return nil, 0, model.DBError
	}
	return items, int64(len(items)), nil
}

// 定时同步
func (d *VideoDomain) SyncVisitCountToDatabase(ctx context.Context) *errs.BError {
	// 1. 获取所有 video:visit:* 键
	keys, err := d.cache.GetAllVisitCountKeys(ctx)
	if err != nil {
		logs.LG.Error("VideoDomain.GetAllVisitCountKeys error", zap.Error(err))
		return model.RedisError
	}
	// 2. 批量读取访问量（使用 Pipeline）
	ids, _ := utils.ExtractVideoIDsFromKeys(keys)
	visits, err := d.cache.GetVisitCountsBatch(ctx, ids)
	if err != nil {
		logs.LG.Error("VideoDomain.GetVisitCountsBatch error", zap.Error(err))
		return model.RedisError
	}
	// 3. 批量更新数据库，并获取更新后的总访问量
	totalVisitCounts, err := d.videoRepo.BatchUpdateVisitCount(ctx, visits)
	if err != nil {
		logs.LG.Error("VideoDomain.BatchUpdateVisitCount error", zap.Error(err))
		return model.RedisError
	}

	// 4. 批量更新 Redis ZSET（使用数据库的总访问量，而不是增量）
	// 将 totalVisitCounts map 转换为 []*redis.Z 格式
	members := make([]*redis.Z, 0, len(totalVisitCounts))
	for videoID, totalVisitCount := range totalVisitCounts {
		members = append(members, &redis.Z{
			Score:  float64(totalVisitCount),        // 使用数据库的总访问量作为分数
			Member: strconv.FormatUint(videoID, 10), // 视频ID作为成员
		})
	}

	if len(members) > 0 {
		err = d.cache.ZAddBatch(ctx, dao.VideoRankingKey, members)
		if err != nil {
			logs.LG.Error("VideoDomain.ZAddBatch error", zap.Error(err))
			return model.RedisError
		}
	}

	// 5. 清空计数器（可选）
	err = d.cache.DelVisitCountBatch(ctx, ids)
	if err != nil {
		logs.LG.Error("VideoDomain.DelVisitCountBatch error", zap.Error(err))
		// 清空失败不影响主流程，只记录日志
	}

	return nil
}
func (d *VideoDomain) ParseMinioUrl(ctx context.Context, bytes []byte, idStr string) (string, *errs.BError) {
	filename := idStr + "" + time.Now().Format("2006_01_02_150405") + ".mp4"
	contentType := minio.GetContentType(bytes)
	minioCli, err := minio.NewMinioClient(config.C.MinIoConfig.Endpoint, config.C.MinIoConfig.AccessKey, config.C.MinIoConfig.SecretKey, config.C.MinIoConfig.UseSSL)
	if err != nil {
		logs.LG.Error("ParseMinioUrl error", zap.Error(err))
		return "", model.MinioError
	}
	_, err = minioCli.Put(ctx, config.C.MinIoConfig.BucketName, filename, bytes, int64(len(bytes)), contentType)
	if err != nil {
		logs.LG.Error("minio.Minio.Put error", zap.Error(err))
		return "", model.MinioError
	}
	fileUrl := minio.ServerUrl + config.C.MinIoConfig.BucketName + "/" + filename
	return fileUrl, nil
}

// VerifyMp4 检查提供的字节切片是否代表一个有效的 MP4 文件。
// 它通过检查文件签名（'ftyp' 盒子）来实现。
func (d *VideoDomain) VerifyMp4(ctx context.Context, data []byte) *errs.BError {
	// 一个 MP4 文件的 ftyp box 签名通常在文件开头。
	// 我们至少需要 12 个字节来读取盒子的大小、类型和主品牌。
	if len(data) < 12 {
		return model.CantLikeMp4
	}

	// 'ftyp' 盒子的类型标识符应该从第 4 个字节开始。
	// 第 4 到第 8 个字节应该是 "ftyp"。
	if !bytes.Equal(data[4:8], []byte("ftyp")) {
		return model.MP4NoEffect
	}

	// 为了更准确地验证，我们还可以检查主品牌。
	// 主品牌位于第 8 个字节之后。
	// 常见的 MP4 兼容主品牌列表。
	majorBrands := [][]byte{
		[]byte("isom"), // ISO Base Media File Format
		[]byte("iso2"),
		[]byte("avc1"), // Advanced Video Coding
		[]byte("mp41"),
		[]byte("mp42"),
		// QuickTime 也是 MP4 的基础，通常是兼容的。
		[]byte("qt  "), // 注意末尾的空格
	}

	majorBrand := data[8:12]
	isCompatible := false
	for _, brand := range majorBrands {
		if bytes.Equal(majorBrand, brand) {
			isCompatible = true
			break
		}
	}

	if !isCompatible {
		// 这是一个更严格的检查。如果你想更宽松，只接受任何包含 'ftyp' 的文件，
		// 可以注释掉或移除这个检查。
		return model.MP4NoSafe
	}

	// 如果所有检查都通过，我们认为这是一个有效的 MP4 文件。
	return nil
}

func (d *VideoDomain) CreateVideo(ctx context.Context, url string, id int64, des, title string) *errs.BError {
	v := &data.Video{
		UserID:      uint64(id),
		VideoURL:    url,
		Description: des,
		Title:       title,
		CreatedAt:   time.Now(),
	}
	err := d.videoRepo.CreateVideo(ctx, v)
	if err != nil {
		logs.LG.Error("VideoDomain.CreateVideo error", zap.Error(err))
		return model.DBError
	}
	return nil
}

func (d *VideoDomain) FindVideosByUserId(ctx context.Context, id int64, size int32, page int32) ([]*data.Video, int64, *errs.BError) {
	videos, total, err := d.videoRepo.FindVideosByUserId(ctx, id, size, page)
	if err != nil {
		logs.LG.Error("VideoDomain.FindVideosByUserId error", zap.Error(err))
		return nil, 0, model.DBError
	}
	return videos, total, nil
}

func (d *VideoDomain) FindVideosById(ctx context.Context, id int64) (*data.Video, *errs.BError) {
	video, err := d.videoRepo.FindVideosById(ctx, id)
	if err != nil {
		logs.LG.Error("VideoDomain.FindVideosByUserId error", zap.Error(err))
		return nil, model.DBError
	}
	return video, nil
}

func (d *VideoDomain) VerifyTime(ctx context.Context, fromDate int64, toDate int64) (time.Time, time.Time, *errs.BError) {
	err := d.VerifySeed(ctx, fromDate)
	if err != nil {
		logs.LG.Error("VideoDomain.VerifySeed error", zap.Error(err))
		return time.Time{}, time.Time{}, err
	}
	err = d.VerifySeed(ctx, toDate)
	if err != nil {
		logs.LG.Error("VideoDomain.VerifySeed error", zap.Error(err))
		return time.Time{}, time.Time{}, err
	}
	if toDate > 0 && fromDate < toDate {
		return time.Time{}, time.Time{}, model.FromTimeMustBigger
	}
	return time.UnixMilli(fromDate), time.UnixMilli(toDate), nil
}

func (d *VideoDomain) FindVideosByTimeAndUserName(ctx context.Context, ft time.Time, tt time.Time, userId int64, size int32, num int32, keyword string) ([]*data.Video, int64, *errs.BError) {
	videos := make([]*data.Video, 0)
	var err error
	var total int64
	if keyword == "" {
		videos, total, err = d.videoRepo.FindVideosByTimeAndUserName(ctx, ft, tt, userId, size, num)
		if err != nil {
			logs.LG.Error("VideoDomain.FindVideosByTimeAndUserName error", zap.Error(err))
			return nil, 0, model.DBError
		}
	} else {
		videos, total, err = d.videoRepo.FindVideosByTimeAndUserNameWithKeyWord(ctx, ft, tt, userId, size, num, keyword)
		if err != nil {
			logs.LG.Error("VideoDomain.FindVideosByTimeAndUserNameWithKeyWord error", zap.Error(err))
			return nil, 0, model.DBError
		}
	}

	return videos, total, nil
}

func (d *VideoDomain) FindVideosByIds(ctx context.Context, ids []int64) ([]*data.Video, *errs.BError) {
	videos, err := d.videoRepo.FindVideosByIds(ctx, ids)
	if err != nil {
		logs.LG.Error("VideoDomain.FindVideosByIds error", zap.Error(err))
		return nil, model.DBError
	}
	return videos, nil
}

func (d *VideoDomain) IncrLikeCount(ctx context.Context, ids int64, isLike bool) *errs.BError {
	var err error
	if isLike {
		err = d.videoRepo.IncrLikeCount(ctx, ids)
	} else {
		err = d.videoRepo.DecrLikeCount(ctx, ids)
	}
	if err != nil {
		logs.LG.Error("VideoDomain.IncrLikeCount error", zap.Error(err))
		return model.DBError
	}
	return nil
}
