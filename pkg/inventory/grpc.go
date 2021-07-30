package inventory

import (
	"context"

	"github.com/d-leme/tradew-inventory-write/pkg/inventory/proto"
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

	servReq := &LockItemsRequest{
		OwnerID:  req.OwnerID,
		LockedBy: req.LockedBy,
		Items:    make([]*LockItemModel, len(req.Items)),
	}

	for i, item := range req.Items {
		servReq.Items[i] = &LockItemModel{
			ID:       item.Id,
			Quantity: item.Quantity,
		}
	}

	if err := s.service.LockItems(ctx, req.UserID, servReq); err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

// TradeItems ...
func (s *grpcService) TradeItems(ctx context.Context, req *proto.TradeItemsRequest) (*proto.Empty, error) {

	servReq := &LockItemsRequest{
		LockedBy: req.LockedBy,
		Items:    make([]*LockItemModel, len(req.Items)),
	}

	for i, item := range req.Items {
		servReq.Items[i] = &LockItemModel{
			ID:       item.Id,
			Quantity: item.Quantity,
		}
	}

	if err := s.service.TradeItems(ctx, servReq); err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}
