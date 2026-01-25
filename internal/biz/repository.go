package biz

import (
	"context"

	"gorm.io/gorm"
)

// Repository 通用仓储接口
type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity T) (int, error)
	Delete(ctx context.Context, id any) (int, error)
	DeleteBatch(ctx context.Context, ids any) (int, error)
	DeleteByCondition(ctx context.Context, condition string, args ...any) (int, error)
	FindByID(ctx context.Context, id any) (T, error)
	FindList(ctx context.Context, query QueryBuilder[T]) ([]T, int64, error)
	FindAll(ctx context.Context, query QueryBuilder[T]) ([]T, error)
	FindOne(ctx context.Context, query QueryBuilder[T]) (T, error)
	Count(ctx context.Context, query QueryBuilder[T]) (int64, error)
	Exists(ctx context.Context, condition string, args ...any) (bool, error)
}

// repositoryImpl 泛型实现类
type RepositoryImpl[T any] struct {
	DB *gorm.DB
}

// NewRepository 实例化仓储
func NewRepository[T any](db *gorm.DB) Repository[T] {
	return &RepositoryImpl[T]{DB: db}
}

// --- QueryBuilder 定义 ---

// QueryBuilder 负责构建查询逻辑
type QueryBuilder[T any] struct {
	db       *gorm.DB
	page     int
	pageSize int
}

// NewQueryBuilder 实例化构建器
func NewQueryBuilder[T any](db *gorm.DB) *QueryBuilder[T] {
	return &QueryBuilder[T]{
		db: db.Model(new(T)), // 预设模型
	}
}

// Where 添加基础查询条件
func (qb *QueryBuilder[T]) Where(query any, args ...any) *QueryBuilder[T] {
	qb.db = qb.db.Where(query, args...)
	return qb
}

// Order 添加排序
func (qb *QueryBuilder[T]) Order(value any) *QueryBuilder[T] {
	qb.db = qb.db.Order(value)
	return qb
}

// Paginate 设置分页参数
func (qb *QueryBuilder[T]) Paginate(page, pageSize int) *QueryBuilder[T] {
	qb.page = page
	qb.pageSize = pageSize
	return qb
}

// Apply 核心方法：支持应用外部 Scope（如 search.MakeCondition）
func (qb *QueryBuilder[T]) Apply(scopes ...func(*gorm.DB) *gorm.DB) *QueryBuilder[T] {
	for _, scope := range scopes {
		if scope != nil {
			qb.db = qb.db.Scopes(scope)
		}
	}
	return qb
}

// Build 核心优化：使用 Session 分离状态
func (qb *QueryBuilder[T]) Build() *gorm.DB {
	// .Session(&gorm.Session{}) 会创建一个新的 DB 实例副本
	// 确保本次 Build 添加的 Limit/Offset 不会影响原始 qb.DB 状态
	db := qb.db.Session(&gorm.Session{})
	if qb.pageSize > 0 {
		offset := (qb.page - 1) * qb.pageSize
		db = db.Offset(offset).Limit(qb.pageSize)
	}
	return db
}

// Preload 预加载关联字段
func (qb *QueryBuilder[T]) Preload(query string, args ...any) *QueryBuilder[T] {
	qb.db = qb.db.Preload(query, args...)
	return qb
}

// Joins 关联查询
func (qb *QueryBuilder[T]) Joins(query string, args ...any) *QueryBuilder[T] {
	qb.db = qb.db.Joins(query, args...)
	return qb
}

// --- Repository 接口实现 (基于 gorm.G) ---

func (r *RepositoryImpl[T]) Create(ctx context.Context, entity *T) error {
	result := gorm.WithResult()
	return gorm.G[T](r.DB, result).Create(ctx, entity)
}

func (r *RepositoryImpl[T]) Update(ctx context.Context, entity T) (int, error) {
	// 使用 Updates 保证兼容性
	return gorm.G[T](r.DB).Updates(ctx, entity)
}

func (r *RepositoryImpl[T]) Delete(ctx context.Context, id any) (int, error) {
	return gorm.G[T](r.DB).Where("id = ?", id).Delete(ctx)
}

func (r *RepositoryImpl[T]) DeleteBatch(ctx context.Context, ids any) (int, error) {
	return gorm.G[T](r.DB).Where("id IN ?", ids).Delete(ctx)
}

func (r *RepositoryImpl[T]) DeleteByCondition(ctx context.Context, condition string, args ...any) (int, error) {
	return gorm.G[T](r.DB).Where(condition, args...).Delete(ctx)
}

func (r *RepositoryImpl[T]) FindByID(ctx context.Context, id any) (T, error) {
	return gorm.G[T](r.DB).Where("id = ?", id).First(ctx)
}

func (r *RepositoryImpl[T]) FindList(ctx context.Context, qb QueryBuilder[T]) (list []T, total int64, err error) {
	// 1. 使用独立的 Session 执行 Count，确保不带 Offset/Limit
	// 即使 qb.DB 里有 Scopes，Count 也会正确执行
	if err = qb.db.Session(&gorm.Session{}).WithContext(ctx).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 2. 如果总数为 0，直接返回，节省一次 Find 查询
	if total == 0 {
		return []T{}, 0, nil
	}

	// 3. 执行带分页的 Find。直接使用 &list 接收，避开 gorm.G 可能存在的重置问题
	err = qb.Build().WithContext(ctx).Find(&list).Error
	// list, err := gorm.G[T](qb.Build()).Find(ctx) 这种方式会导致问题，这种内含重置导致scope失效
	return list, total, err
}

func (r *RepositoryImpl[T]) FindAll(ctx context.Context, qb QueryBuilder[T]) ([]T, error) {
	return gorm.G[T](qb.Build()).Find(ctx)
}

func (r *RepositoryImpl[T]) FindOne(ctx context.Context, qb QueryBuilder[T]) (T, error) {
	return gorm.G[T](qb.Build()).First(ctx)
}

func (r *RepositoryImpl[T]) Count(ctx context.Context, qb QueryBuilder[T]) (int64, error) {
	var total int64
	err := qb.db.Count(&total).Error
	return total, err
}

func (r *RepositoryImpl[T]) Exists(ctx context.Context, condition string, args ...any) (bool, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(new(T)).Where(condition, args...).Count(&count).Error
	return count > 0, err
}

// Transaction 执行事务
func Transaction(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	return db.Transaction(fn)
}

// TransactionWithResult 执行事务并返回结果
func TransactionWithResult[T any](db *gorm.DB, fn func(tx *gorm.DB) (T, error)) (T, error) {
	var result T
	err := db.Transaction(func(tx *gorm.DB) error {
		var err error
		result, err = fn(tx)
		return err
	})
	return result, err
}
