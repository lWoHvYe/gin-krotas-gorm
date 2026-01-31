package service

import (
	"context"
	"helloworld-go/internal/data/persistence"

	pb "helloworld-go/api/product/v1"
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
