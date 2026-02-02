package service

import (
	"context"
	"database/sql"
	"helloworld-go/internal/product/data/persistence"

	pb "helloworld-go/api/product/v1"

	"github.com/dtm-labs/client/dtmgrpc"
	"gorm.io/gorm"
)

type ProductService struct {
	pb.UnimplementedProductServer
	repo *persistence.ProductRepo
}

func NewProductService(repo *persistence.ProductRepo) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetProductDetail(ctx context.Context, req *pb.GetProductDetailReq) (*pb.GetProductDetailReply, error) {
	// 1. 调用 Biz/Data 层获取数据 (预加载 Skus)
	// 假设 repo.GetProductWithSkus 内部使用了 .Preload("Skus")
	spu, err := s.repo.GetProductWithSkus(ctx, int64(req.Id))
	if err != nil {
		return nil, err
	}

	// 2. 将数据模型转换为 Proto 结构体
	reply := &pb.GetProductDetailReply{
		Id:          uint32(spu.ID),
		Name:        spu.Name,
		Description: spu.Description,
		MainImage:   spu.MainImage,
		CategoryId:  uint32(spu.CategoryID),
		BrandId:     uint32(spu.BrandID),
	}

	// 3. 循环转换子列表 SKU
	for _, item := range spu.Skus {
		reply.Skus = append(reply.Skus, &pb.SkuInfo{
			Id:        uint32(item.ID),
			SkuCode:   item.SkuCode,
			AttrValue: item.AttrValue,
			Price:     item.Price,
			Stock:     int32(item.Stock),
		})
	}

	return reply, nil
}
func (s *ProductService) ReduceStock(ctx context.Context, req *pb.ReduceStockReq) (*pb.SimpleResponse, error) {
	// 获取屏障对象
	barrier, err := dtmgrpc.BarrierFromGrpc(ctx)
	if err != nil {
		return nil, err
	}

	// 在屏障内执行数据库操作
	// UnderlyingDB() 获取的是你在 Data 层初始化的 *gorm.DB
	gormDB := s.repo.GetDB(ctx)
	sqlDB, err := gormDB.DB() // GORM 1.31+ 获取 *sql.DB 的标准方法
	if err != nil {
		return nil, err
	}
	// 这里传入 *sql.DB，但在回调中使用 gormDB 的事务
	err = barrier.CallWithDB(sqlDB, func(tx *sql.Tx) error {
		// 关键点：将原生 sql.Tx 转换回 GORM 的事务对象供业务使用
		// 使用 Use(gormDB.WithContext(ctx).Begin(opts...)) 或者直接：
		return gormDB.Transaction(func(gtx *gorm.DB) error {
			// 这里执行你的业务逻辑，如 q := query.Use(gtx)
			return s.repo.ReduceStockInTx(ctx, gtx, req.OrderSn, req.SkuId, req.Num)
		})
	})
	return &pb.SimpleResponse{Msg: "success"}, err
}
func (s *ProductService) CompensateStock(ctx context.Context, req *pb.ReduceStockReq) (*pb.SimpleResponse, error) {
	err := s.repo.CompensateStock(ctx, req.OrderSn, req.SkuId, req.Num)
	return &pb.SimpleResponse{Msg: "done"}, err
}
