package service

import (
	pb "Github.com/LocalEats/Order-Service/gen-proto/order"
	"Github.com/LocalEats/Order-Service/internal/repository"
	"context"
)

type OrderService struct {
	OrderRepo *repository.OrderRepository
	pb.UnimplementedOrderServiceServer
}

func NewOrderService(orderRepo repository.OrderRepository) *OrderService {
	return &OrderService{
		OrderRepo: &orderRepo,
	}
}
func (s *OrderService) CreateDish(ctx context.Context, req *pb.CreateDishRequest) (*pb.CreateDishResponse, error) {
	return s.OrderRepo.CreateDish(ctx, req)
}

func (s *OrderService) UpdateDish(ctx context.Context, req *pb.UpdateDishRequest) (*pb.UpdateDishResponse, error) {
	return s.OrderRepo.UpdateDish(ctx, req)
}

func (s *OrderService) DeleteDish(ctx context.Context, req *pb.DeleteDishRequest) (*pb.DeleteDishResponse, error) {
	return s.OrderRepo.DeleteDish(ctx, req)
}

func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	return s.OrderRepo.CreateOrder(ctx, req)
}

func (s *OrderService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	return s.OrderRepo.GetOrder(ctx, req)
}

func (s *OrderService) GetDishes(ctx context.Context, req *pb.ListDishesRequest) (*pb.ListDishesResponse, error) {
	return s.OrderRepo.GetDishes(ctx, req)
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
	return s.OrderRepo.UpdateOrderStatus(ctx, req)
}

func (s *OrderService) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	return s.ListOrders(ctx, req)
}

func (s *OrderService) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*pb.CreateReviewResponse, error) {
	return s.OrderRepo.CreateReview(ctx, req)
}

func (s *OrderService) ListReviews(ctx context.Context, req *pb.ListReviewsRequest) (*pb.ListReviewsResponse, error) {
	return s.OrderRepo.ListReviews(ctx, "ali", req)
}

func (s *OrderService) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	return s.OrderRepo.CreatePayment(ctx, req)
}

func (s *OrderService) GetDishRecommendations(ctx context.Context, req *pb.GetDishRecommendationsRequest) (*pb.GetDishRecommendationsResponse, error) {
	return s.OrderRepo.GetDishRecommendations(ctx, req)
}

func (s *OrderService) GetKitchenStatics(ctx context.Context, req *pb.GetKitchenStatisticsRequest) (*pb.GetKitchenStatisticsResponse, error) {
	return s.OrderRepo.GetKitchenStatistics(ctx, req)
}

func (s *OrderService) GetUserActivity(ctx context.Context, req *pb.GetUserActivityRequest) (*pb.GetUserActivityResponse, error) {
	return s.OrderRepo.GetUserActivity(ctx, req)
}

func (s *OrderService) UpdateWorkingHours(ctx context.Context, req *pb.UpdateWorkingHoursRequest) (*pb.UpdateWorkingHoursResponse, error) {
	return s.OrderRepo.UpdateWorkingHours(ctx, req)
}

func (s *OrderService) UpdateDishNutritionInfo(ctx context.Context, req *pb.UpdateDishNutritionInfoRequest) (*pb.UpdateDishNutritionInfoResponse, error) {
	return s.OrderRepo.UpdateDishNutritionInfo(ctx, req)
}
