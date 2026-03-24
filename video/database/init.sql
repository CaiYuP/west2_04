-- =================================================================
-- west2-video Database Initialization Script
-- Generated based on GORM models from all microservices.
-- =================================================================

-- Create database if it doesn't exist
CREATE DATABASE IF NOT EXISTS `west2_video` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `west2_video`;

-- Drop tables in reverse order of dependency to avoid foreign key errors
DROP TABLE IF EXISTS `likes`;
DROP TABLE IF EXISTS `follows`;
DROP TABLE IF EXISTS `comments`;
DROP TABLE IF EXISTS `videos`;
DROP TABLE IF EXISTS `users`;

-- =================================================================
-- 1. Users Table (users)
-- Based on userService/internal/data/user.go
-- =================================================================
CREATE TABLE `users` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `username` VARCHAR(32) NOT NULL COMMENT '用户名',
    `password` VARCHAR(255) NOT NULL COMMENT '密码',
    `email` VARCHAR(128) NULL COMMENT '邮箱',
    `nickname` VARCHAR(64) NULL COMMENT '昵称',
    `avatar_url` VARCHAR(512) NULL COMMENT '头像URL',
    `description` VARCHAR(512) NULL COMMENT '个人简介',
    `is_mfa_enabled` TINYINT(1) DEFAULT 0 COMMENT '是否启用MFA',
    `mfa_secret` VARCHAR(64) NULL COMMENT 'MFA Secret',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- =================================================================
-- 2. Videos Table (videos)
-- Based on videoService/internal/data/video.go
-- =================================================================
CREATE TABLE `videos` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '视频ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '视频作者',
    `video_url` VARCHAR(512) NOT NULL COMMENT '视频链接',
    `cover_url` VARCHAR(512) NOT NULL COMMENT '封面链接',
    `title` VARCHAR(128) NOT NULL COMMENT '标题',
    `description` VARCHAR(512) NOT NULL COMMENT '描述',
    `visit_count` BIGINT UNSIGNED DEFAULT 0 COMMENT '访问量',
    `like_count` BIGINT UNSIGNED DEFAULT 0 COMMENT '点赞数',
    `comment_count` BIGINT UNSIGNED DEFAULT 0 COMMENT '评论数',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    CONSTRAINT `fk_videos_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='视频表';

-- =================================================================
-- 3. Comments Table (comments)
-- Based on interactionService/internal/data/comment.go
-- =================================================================
CREATE TABLE `comments` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '评论ID',
    `video_id` BIGINT UNSIGNED NOT NULL COMMENT '视频ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '发表者ID',
    `parent_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '父评论ID',
    `like_count` BIGINT UNSIGNED DEFAULT 0 COMMENT '点赞数',
    `child_count` BIGINT UNSIGNED DEFAULT 0 COMMENT '子评论数',
    `content` TEXT NOT NULL COMMENT '评论内容',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    -- 按视频查询评论的索引
    KEY `idx_video_id` (`video_id`),
    -- 按用户查询评论的索引
    KEY `idx_user_id` (`user_id`),
    -- 按父评论查询子评论的索引（递归删除、子评论列表）
    KEY `idx_comments_parent_id` (`parent_id`),
    -- 视频 + 创建时间的联合索引（按时间分页）
    KEY `idx_comments_video_created_at` (`video_id`, `created_at`),
    CONSTRAINT `fk_comments_video` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_comments_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论表';

-- =================================================================
-- 4. Likes Table (likes)
-- 视频/评论点赞共用表，区分 is_comment
-- =================================================================
CREATE TABLE `likes` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '点赞ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `video_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '视频ID；评论点赞时为0',
    `comment_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '评论ID；视频点赞时为0',
    `is_comment` TINYINT(1) NOT NULL COMMENT '是否为评论点赞',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    -- 视频点赞唯一：同一用户对同一视频只能点赞一次
    UNIQUE KEY `uk_like_user_video` (`user_id`, `video_id`, `is_comment`),
    -- 评论点赞唯一：同一用户对同一评论只能点赞一次
    UNIQUE KEY `uk_like_user_comment` (`user_id`, `comment_id`, `is_comment`),
    -- 用户点赞列表（按 is_comment 过滤）
    KEY `idx_like_user_is_comment` (`user_id`, `is_comment`),
    CONSTRAINT `fk_likes_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='点赞表（视频/评论共用）';

-- =================================================================
-- 5. Follows Table (follows)
-- Based on socialService/internal/data/follow.go
-- =================================================================
CREATE TABLE `follows` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '关注ID',
    `follower_id` BIGINT UNSIGNED NOT NULL COMMENT '关注者ID',
    `followee_id` BIGINT UNSIGNED NOT NULL COMMENT '被关注者ID',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    -- 防止重复关注 + 加速“我关注了谁 / 是否关注”
    UNIQUE KEY `uk_follower_followee` (`follower_id`, `followee_id`),
    -- 加速“谁关注了我”（粉丝列表）
    KEY `idx_follows_followee_id` (`followee_id`),
    CONSTRAINT `fk_follows_follower` FOREIGN KEY (`follower_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_follows_followee` FOREIGN KEY (`followee_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='关注表';

