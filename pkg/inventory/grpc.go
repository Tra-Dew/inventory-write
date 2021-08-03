package inventory

import (
	"context"

	"github.com/d-leme/tradew-inventory-write/pkg/inventory/proto"
	"github.com/sirupsen/logrus"
)

type grpcService struct {
	proto.UnimplementedInventoryServiceServer

	service Service
}

// NewGRPCService ...
func NewGRPCService(service Service) proto.InventoryServiceServer {
	return &grpcService{
		service: service,
	}
}

// LockItems ...
func (s *grpcService) LockItems(ctx context.Context, req *proto.LockItemsRequest) (*proto.Empty, error) {

	logrus.Info("lock items called by GRPC")

	servReq := &LockItemsRequest{
		OwnerID:            req.OwnerID,
		LockedBy:           req.LockedBy,
		WantedItemsOwnerID: req.WantedItemsOwnerID,
		OfferedItems:       make([]*LockItemModel, len(req.OfferedItems)),
		WantedItems:        make([]*LockItemModel, len(req.WantedItems)),
	}

	for i, item := range req.OfferedItems {
		servReq.OfferedItems[i] = &LockItemModel{
			ID:       item.Id,
			Quantity: item.Quantity,
		}
	}

	for i, item := range req.WantedItems {
		servReq.WantedItems[i] = &LockItemModel{
			ID:       item.Id,
			Quantity: item.Quantity,
		}
	}

	if err := s.service.LockItems(ctx, servReq); err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

// TradeItems ...
func (s *grpcService) TradeItems(ctx context.Context, req *proto.TradeItemsRequest) (*proto.Empty, error) {

	logrus.Info("trade items called by GRPC")

	servReq := &TradeItemsRequest{
		TradeID:            req.TradeID,
		OwnerID:            req.OwnerID,
		WantedItemsOwnerID: req.WantedItemsOwnerID,
		OfferedItems:       make([]*TradeItemModel, len(req.OfferedItems)),
		WantedItems:        make([]*TradeItemModel, len(req.WantedItems)),
	}

	for i, item := range req.OfferedItems {
		servReq.OfferedItems[i] = &TradeItemModel{
			ID:       item.Id,
			Quantity: item.Quantity,
		}
	}

	for i, item := range req.WantedItems {
		servReq.WantedItems[i] = &TradeItemModel{
			ID:       item.Id,
			Quantity: item.Quantity,
		}
	}

	if err := s.service.TradeItems(ctx, servReq); err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}
