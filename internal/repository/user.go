package repository

import (
	"context"
	"database/sql"
	"time"
	"xiaoweishu/internal/domain"
	"xiaoweishu/internal/repository/cache"
	"xiaoweishu/internal/repository/dao"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
}

type CachedUserRepository struct {
	dao   dao.UserDao
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDao, c cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: c,
	}
}

func (r *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.domainToEntity(u))
}

func (r *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), err
}

func (r *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), err
}

// 非异步
//func (r *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
//	u, err := r.cache.Get(ctx, id)
//	if err == nil {
//		return u, nil
//	}
//	//if errors.Is(err, cache.ErrKeyNotFound) {
//	//	// 去数据库找
//	//}
//	ue, err := r.dao.FindById(ctx, id)
//	if err != nil {
//		return domain.User{}, err
//	}
//	u = r.entityToDomain(ue)
//	_ = r.cache.Set(ctx, u)
//	//go func() {
//	//	err = r.cache.Set(ctx, u)
//	//	if err != nil {
//	//		//return domain.User{}, err
//	//		// ignore cache error
//	//		// 打日志，记录缓存失败
//	//	}
//	//}()
//	return u, nil
//}

// 异步
func (r *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		return u, nil
	}
	//if errors.Is(err, cache.ErrKeyNotFound) {
	//	// 去数据库找
	//}
	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = r.entityToDomain(ue)

	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			//return domain.User{}, err
			// ignore cache error
			// 打日志，记录缓存失败
		}
	}()
	return u, nil
}

func (r *CachedUserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Password: u.Password,
		NickName: u.NicName,
		Birthday: u.BirthDay.String(),
		AboutMe:  u.AboutMe,
		Ctime:    time.UnixMilli(u.Ctime),
	}
}

func (r *CachedUserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{

		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Ctime:    u.Ctime.UnixMilli(),
	}
}
