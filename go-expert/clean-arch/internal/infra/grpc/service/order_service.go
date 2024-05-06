package service

import (
	"context"

	"github.com/albuquerq/full-cycle/go-expert/clean-arch/internal/infra/grpc/pb"
	"github.com/albuquerq/full-cycle/go-expert/clean-arch/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	createOrderUseCase usecase.CreateOrderUseCase
	listOrderUseCase   usecase.ListOrdersUseCase
}

func NewOrderService(
	createOrderUseCase usecase.CreateOrderUseCase,
	listOrdersUseCase usecase.ListOrdersUseCase,
) *OrderService {
	return &OrderService{
		createOrderUseCase: createOrderUseCase,
		listOrderUseCase:   listOrdersUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.createOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, in *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := s.listOrderUseCase.Execute()
	if err != nil {
		return nil, err
	}

	out := make([]*pb.ListOrdersResponse_OrderResponse, len(orders))

	for i := range orders {
		out[i] = &pb.ListOrdersResponse_OrderResponse{
			Id:         orders[i].ID,
			Price:      float32(orders[i].Price),
			Tax:        float32(orders[i].Tax),
			FinalPrice: float32(orders[i].FinalPrice),
		}
	}

	return &pb.ListOrdersResponse{Orders: out}, nil
}
