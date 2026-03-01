package dao

import (
	"context"
)

// Article 制作库的
// 如何设计索引，和 WHERE 相关
// 对于帖子来说，是什么查询场景
// 对于创作者来说，看草稿箱，看到所有自己的文章
// SELECT * FROM artiles WHERE author_id = 1
// 单独查询一篇文章
// SELECT * FROM articles WHERE id = 1
// 假设要求按照创建时间的倒序排序
type Article struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 长度1024
	Title    string `gorm:"type=varchar(1024)"`
	Content  string `gorm:"type=BLOB"`
	AuthorId int64
	Ctime    int64
	Utime    int64
}

type ArticleDao interface {
	Insert(ctx context.Context, article Article) (int64, error)
}
