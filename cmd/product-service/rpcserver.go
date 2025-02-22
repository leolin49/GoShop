package main

import (
	"context"
	productpb "goshop/api/protobuf/product"
	"goshop/models"
	errorcode "goshop/pkg/error"
	"net"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

type ProductRpcService struct {
	productpb.UnimplementedProductServiceServer
}

func rpcServerStart() bool {
	lis, err := net.Listen("tcp", ":49100")
	if err != nil {
		glog.Fatalf("[ProductServer] rpcserver failed to listen: %v", err)
		return false
	}
	rpcServer := grpc.NewServer()
	productpb.RegisterProductServiceServer(rpcServer, new(ProductRpcService))
	glog.Infoln("[ProductServer] Starting rpc server on :49100")
	go func() {
		if err := rpcServer.Serve(lis); err != nil {
			glog.Fatalf("[ProductServer] rpcserver failed to start: %v", err)
			return
		}
	}()
	return true
}

func (s *ProductRpcService) AddProduct(ctx context.Context, req *productpb.ReqAddProduct) (*productpb.RspAddProduct, error) {
	cq := models.NewCategoryQuery(Mysql())
	if exist, err := models.NewProductQuery(Mysql()).ProductExisted(req.Product.Name); exist {
		return &productpb.RspAddProduct{
			ErrorCode: errorcode.ProductAlreadyExist,
		}, err
	}
	ids, err := cq.GetIdsByNames(req.Product.Categories)
	if err != nil {
		return &productpb.RspAddProduct{
			ErrorCode: errorcode.UnknowError,
		}, err
	}
	categories := make([]models.Category, len(ids))
	for i, id := range ids {
		categories[i].ID = uint(id)
		categories[i].Name = req.Product.Categories[i]
	}
	product := &models.Product{
		Name:        req.Product.Name,
		Description: req.Product.Description,
		Picture:     req.Product.Picture,
		Price:       req.Product.Price,
		Categories:  categories,
	}
	res := Mysql().Create(product)
	if err := res.Error; err != nil {
		return &productpb.RspAddProduct{
			ErrorCode: errorcode.UnknowError,
		}, err
	}
	return &productpb.RspAddProduct{
		ErrorCode: errorcode.Ok,
	}, nil
}

func (s *ProductRpcService) ListProducts(ctx context.Context, req *productpb.ReqListProducts) (*productpb.RspListProducts, error) {
	res, err := models.NewCategoryQuery(Mysql()).ListProducts(req.CategoryName, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	ret := &productpb.RspListProducts{}
	for _, pd := range res {
		ret.Products = append(ret.Products, &productpb.Product{
			Name:        pd.Name,
			Description: pd.Description,
			Picture:     pd.Picture,
			Price:       pd.Price,
		})
	}
	return ret, nil
}

func (s *ProductRpcService) GetProduct(ctx context.Context, req *productpb.ReqGetProduct) (*productpb.RspGetProduct, error) {
	res, err := models.NewProductQuery(Mysql()).GetById(int32(req.Id))
	if err != nil {
		return nil, err
	}
	return &productpb.RspGetProduct{
		Product: &productpb.Product{
			Id:          uint32(res.ID),
			Name:        res.Name,
			Description: res.Description,
			Picture:     res.Picture,
			Price:       res.Price,
		},
	}, nil
}

func (s *ProductRpcService) SearchProducts(ctx context.Context, req *productpb.ReqSearchProducts) (*productpb.RspSearchProducts, error) {
	res, err := models.NewProductQuery(Mysql()).SearchProducts(req.Query)
	if err != nil {
		return nil, err
	}
	ret := &productpb.RspSearchProducts{}
	for _, pd := range res {
		ret.Results = append(ret.Results, &productpb.Product{
			Name:        pd.Name,
			Description: pd.Description,
			Picture:     pd.Picture,
			Price:       pd.Price,
		})
	}
	return ret, nil
}
