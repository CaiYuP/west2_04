# 数据库设计说明

## 数据库表结构

### 1. 用户表 (users)

存储用户基本信息。

**字段说明：**
- `id`: 用户ID，主键，自增
- `username`: 用户名，唯一索引，用于登录
- `password`: 密码，BCrypt 加密存储
- `email`: 邮箱（可选）
- `nickname`: 昵称
- `avatar_url`: 头像URL
- `description`: 个人简介
- `created_at`: 创建时间
- `updated_at`: 更新时间

**索引：**
- 主键：`id`
- 唯一索引：`username`（防止重复注册）
- 普通索引：`created_at`（按时间排序）

---

### 2. 视频表 (videos)

存储视频信息。

**字段说明：**
- `id`: 视频ID，主键，自增
- `author_id`: 作者ID，外键关联 `users.id`
- `play_url`: 视频播放地址
- `cover_url`: 封面图URL
- `title`: 视频标题
- `description`: 视频描述
- `favorite_count`: 点赞数（冗余字段，用于性能优化）
- `comment_count`: 评论数（冗余字段）
- `visit_count`: 点击/播放次数（用于热门排行榜）
- `created_at`: 创建时间
- `updated_at`: 更新时间

**索引：**
- 主键：`id`
- 外键：`author_id` → `users.id`
- 普通索引：
  - `author_id`（查询用户的发布列表）
  - `created_at`（视频流按时间排序）
  - `visit_count DESC`（热门排行榜）
  - `title, description`（视频搜索）

**设计考虑：**
- `favorite_count` 和 `comment_count` 是冗余字段，通过触发器或应用层维护，避免每次查询都 COUNT
- `visit_count` 用于热门排行榜，需要 Redis 缓存优化

---

### 3. 点赞表 (likes)

存储用户对视频的点赞关系。

**字段说明：**
- `id`: 点赞ID，主键，自增
- `user_id`: 用户ID，外键关联 `users.id`
- `video_id`: 视频ID，外键关联 `videos.id`
- `created_at`: 创建时间

**索引：**
- 主键：`id`
- 唯一索引：`(user_id, video_id)`（防止重复点赞）
- 普通索引：
  - `user_id`（查询用户的点赞列表）
  - `video_id`（统计视频点赞数）

**设计考虑：**
- 唯一索引确保一个用户对一个视频只能点赞一次
- 点赞/取消点赞时，需要同步更新 `videos.favorite_count`

---

### 4. 评论表 (comments)

存储视频评论。

**字段说明：**
- `id`: 评论ID，主键，自增
- `video_id`: 视频ID，外键关联 `videos.id`
- `user_id`: 用户ID，外键关联 `users.id`
- `content`: 评论内容
- `created_at`: 创建时间
- `updated_at`: 更新时间

**索引：**
- 主键：`id`
- 外键：`video_id` → `videos.id`，`user_id` → `users.id`
- 普通索引：
  - `video_id`（查询视频的评论列表）
  - `user_id`（查询用户的评论）
  - `created_at`（按时间排序）

**设计考虑：**
- 本次作业不要求实现"对评论进行评论"，所以没有 `parent_comment_id` 字段
- 删除评论时需要校验：只能删除自己的评论

---

### 5. 关注表 (follows)

存储用户之间的关注关系。

**字段说明：**
- `id`: 关注ID，主键，自增
- `follower_id`: 关注者ID（粉丝），外键关联 `users.id`
- `followee_id`: 被关注者ID（被关注的人），外键关联 `users.id`
- `created_at`: 创建时间

**索引：**
- 主键：`id`
- 唯一索引：`(follower_id, followee_id)`（防止重复关注）
- 普通索引：
  - `follower_id`（查询关注列表：我关注了谁）
  - `followee_id`（查询粉丝列表：谁关注了我）

**设计考虑：**
- 唯一索引确保一个用户不能重复关注同一个人
- 好友列表 = 互相关注（`follower_id` 和 `followee_id` 互换查询）
- 添加 CHECK 约束：不能关注自己

---

## 数据库设计原则

1. **外键约束**：使用 CASCADE 删除，确保数据一致性
2. **冗余字段**：`favorite_count`、`comment_count` 用于性能优化，需要应用层维护
3. **唯一索引**：防止重复数据（用户名、点赞、关注）
4. **索引优化**：根据查询场景建立合适的索引
5. **字符集**：使用 `utf8mb4` 支持 emoji 等特殊字符

---

## 使用说明

### 初始化数据库

```bash
# 方式1：使用 MySQL 客户端
mysql -u root -p < database/init.sql

# 方式2：在 MySQL 中执行
source database/init.sql;
```

### 注意事项

1. 确保 MySQL 版本 >= 5.7（支持 JSON 和 CHECK 约束）
2. 建议使用 InnoDB 引擎（支持事务和外键）
3. 生产环境建议调整字符集和排序规则
4. 根据实际数据量调整索引策略

---

## 后续优化建议

1. **分表策略**：当数据量增大时，可以考虑按时间分表（如按月分表）
2. **读写分离**：主库写，从库读
3. **缓存策略**：
   - 热门排行榜使用 Redis 缓存
   - 用户信息、视频信息使用 Redis 缓存
4. **全文搜索**：视频搜索可以使用 ElasticSearch 或 MySQL 全文索引





