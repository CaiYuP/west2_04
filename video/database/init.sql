-- ============================================
-- 视频网站数据库初始化脚本
-- ============================================

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS `west2_video` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `west2_video`;

-- ============================================
-- 1. 用户表 (users)
-- ============================================
CREATE TABLE IF NOT EXISTS `users` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `username` VARCHAR(32) NOT NULL COMMENT '用户名（唯一）',
    `password` VARCHAR(255) NOT NULL COMMENT '密码（BCrypt 加密）',
    `email` VARCHAR(128) DEFAULT NULL COMMENT '邮箱',
    `nickname` VARCHAR(64) DEFAULT NULL COMMENT '昵称',
    `avatar_url` VARCHAR(512) DEFAULT NULL COMMENT '头像URL',
    `description` VARCHAR(512) DEFAULT NULL COMMENT '个人简介',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ============================================
-- 2. 视频表 (videos)
-- ============================================
CREATE TABLE IF NOT EXISTS `videos` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '视频ID',
    `author_id` BIGINT UNSIGNED NOT NULL COMMENT '作者ID',
    `play_url` VARCHAR(512) NOT NULL COMMENT '视频播放地址',
    `cover_url` VARCHAR(512) DEFAULT NULL COMMENT '封面图URL',
    `title` VARCHAR(128) NOT NULL COMMENT '视频标题',
    `description` VARCHAR(512) DEFAULT NULL COMMENT '视频描述',
    `favorite_count` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '点赞数（冗余字段，用于性能优化）',
    `comment_count` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '评论数（冗余字段）',
    `visit_count` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '点击/播放次数（用于热门排行榜）',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_author_id` (`author_id`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_visit_count` (`visit_count` DESC) COMMENT '用于热门排行榜',
    KEY `idx_title_desc` (`title`, `description`) COMMENT '用于搜索',
    CONSTRAINT `fk_videos_author` FOREIGN KEY (`author_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='视频表';

-- ============================================
-- 3. 点赞表 (likes)
-- ============================================
CREATE TABLE IF NOT EXISTS `likes` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '点赞ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `video_id` BIGINT UNSIGNED NOT NULL COMMENT '视频ID',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_video` (`user_id`, `video_id`) COMMENT '防止重复点赞',
    KEY `idx_user_id` (`user_id`),
    KEY `idx_video_id` (`video_id`),
    CONSTRAINT `fk_likes_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_likes_video` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='点赞表';

-- ============================================
-- 4. 评论表 (comments)
-- ============================================
CREATE TABLE IF NOT EXISTS `comments` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '评论ID',
    `video_id` BIGINT UNSIGNED NOT NULL COMMENT '视频ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `content` TEXT NOT NULL COMMENT '评论内容',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_video_id` (`video_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_created_at` (`created_at`),
    CONSTRAINT `fk_comments_video` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_comments_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论表';

-- ============================================
-- 5. 关注表 (follows)
-- ============================================
CREATE TABLE IF NOT EXISTS `follows` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '关注ID',
    `follower_id` BIGINT UNSIGNED NOT NULL COMMENT '关注者ID（粉丝）',
    `followee_id` BIGINT UNSIGNED NOT NULL COMMENT '被关注者ID（被关注的人）',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_follower_followee` (`follower_id`, `followee_id`) COMMENT '防止重复关注',
    KEY `idx_follower_id` (`follower_id`) COMMENT '用于查询关注列表',
    KEY `idx_followee_id` (`followee_id`) COMMENT '用于查询粉丝列表',
    CONSTRAINT `fk_follows_follower` FOREIGN KEY (`follower_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_follows_followee` FOREIGN KEY (`followee_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CHECK (`follower_id` != `followee_id`) COMMENT '不能关注自己'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='关注表';

-- ============================================
-- 索引说明
-- ============================================
-- users 表：
--   - uk_username: 用户名唯一索引，用于登录和注册校验
--   - idx_created_at: 按创建时间排序
--
-- videos 表：
--   - idx_author_id: 用于查询用户的发布列表
--   - idx_created_at: 用于视频流按时间排序
--   - idx_visit_count: 用于热门排行榜（降序）
--   - idx_title_desc: 用于视频搜索（标题和描述）
--
-- likes 表：
--   - uk_user_video: 唯一索引，防止重复点赞
--   - idx_user_id: 用于查询用户的点赞列表
--   - idx_video_id: 用于统计视频点赞数
--
-- comments 表：
--   - idx_video_id: 用于查询视频的评论列表
--   - idx_user_id: 用于查询用户的评论
--   - idx_created_at: 用于按时间排序评论
--
-- follows 表：
--   - uk_follower_followee: 唯一索引，防止重复关注
--   - idx_follower_id: 用于查询关注列表（我关注了谁）
--   - idx_followee_id: 用于查询粉丝列表（谁关注了我）





