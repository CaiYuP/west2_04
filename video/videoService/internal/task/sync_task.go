package task

import (
	"context"
	"go.uber.org/zap"
	"time"
	"videoService/internal/domain"
	"west2-video/common/logs"
)

//启动访问量同步任务
func StartVisitCountSyncTask(syncInterval time.Duration) {
	videoDomain := domain.NewVideoDomain()

	// 启动时立即执行一次
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	err := videoDomain.SyncVisitCountToDatabase(ctx)
	cancel()
	if err != nil {
		logs.LG.Error("VideoDomain.SyncVisitCountToDatabase error", zap.Error(err))
	}

	// 然后启动定时任务
	ticker := time.NewTicker(syncInterval)
	go func() {
		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			err := videoDomain.SyncVisitCountToDatabase(ctx)
			cancel()
			if err != nil {
				logs.LG.Error("VideoDomain.SyncVisitCountToDatabase error", zap.Error(err))
			}
		}
	}()
}
