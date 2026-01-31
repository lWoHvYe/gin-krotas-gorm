package service

import (
	"context"
	"helloworld-go/internal/biz/product"

	pb "helloworld-go/api/product/v1"
)

type ProductService struct {
	pb.UnimplementedProductServer
	productRepo product.ProductRepository
}

func NewProductService(repo product.ProductRepository) *ProductService {
	return &ProductService{productRepo: repo}
}

func (s *ProductService) GetProductDetail(ctx context.Context, req *pb.GetProductDetailReq) (*pb.GetProductDetailReply, error) {
	// 1. 调用 Biz/Data 层获取数据 (预加载 Skus)
	// 假设 repo.GetProductWithSkus 内部使用了 .Preload("Skus")
	spu, err := s.productRepo.GetProductWithSkus(ctx, uint(req.Id))
	if err != nil {
		return nil, err
	}

	// 2. 将数据模型转换为 Proto 结构体
	reply := &pb.GetProductDetailReply{
		Id:          uint32(spu.ID),
		Name:        spu.Name,
		Description: spu.Description,
		MainImage:   spu.MainImage,
		CategoryId:  uint32(spu.CategoryId),
		BrandId:     uint32(spu.BrandId),
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
