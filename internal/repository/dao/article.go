package dao

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Article 制作库的
// 如何设计索引，和 WHERE 相关
// 对于帖子来说，是什么查询场景
// 对于创作者来说，看草稿箱，看到所有自己的文章
// SELECT * FROM artiles WHERE author_id = 1
// 单独查询一篇文章
// SELECT * FROM articles WHERE id = 1
// 假设要求按照创建时间的倒序排序
// SELECT * FROM artiles WHERE author_id = 1 ORDER BY ctime DESC
// 因此最佳实现，需要在 author_id 和 ctime 创建联合索引
type Article struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 长度1024
	Title    string `gorm:"type=varchar(1024)"`
	Content  string `gorm:"type=BLOB"`
	AuthorId int64  `gorm:"index=aid_ctime"`
	Ctime    int64  `gorm:"index=aid_ctime"`
	Utime    int64
}

type ArticleDao interface {
	Insert(ctx context.Context, article Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
}

type GormArticleDao struct {
	db *gorm.DB
}

func NewGormArticleDao(db *gorm.DB) ArticleDao {
	return &GormArticleDao{
		db: db,
	}
}

func (dao *GormArticleDao) Insert(ctx context.Context, article Article) (int64, error) {
	now := time.Now().UnixMilli()
	article.Ctime = now
	article.Utime = now
	err := dao.db.WithContext(ctx).Create(&article).Error
	return article.Id, err
}

func (dao *GormArticleDao) UpdateById(ctx context.Context, article Article) error {
	now := time.Now().UnixMilli()
	article.Utime = now
	// gorm 忽略零值特性，使用主键进行更新
	res := dao.db.WithContext(ctx).Model(&article).
		Where("id=? AND author_id=?", article.Id, article.AuthorId).
		Updates(map[string]any{
			"title":   article.Title,
			"content": article.Content,
			"utime":   article.Utime,
		})
	// 需不需要检查是否真的更新
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 { // 更新行数
		return fmt.Errorf("更新失败，文章不存在或非作者本人, id: %d, author_id: %d", article.Id, article.AuthorId)
	}
	return res.Error
}
