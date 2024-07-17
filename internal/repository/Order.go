package repository

import (
	"context"
	"errors"
	"time"

	pb "Github.com/LocalEats/Order-Service/gen-proto/order"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	pq "github.com/lib/pq"
	"github.com/spf13/cast"
	"go.uber.org/zap"

	l "Github.com/LocalEats/Order-Service/internal/config/logger"

	"database/sql"
)

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (o *OrderRepository) CreateDish(ctx context.Context, req *pb.CreateDishRequest) (*pb.CreateDishResponse, error) {

	log, err := l.NewLogger()
	if err != nil {
		return nil, err
	}

	resp := &pb.CreateDishResponse{}
	dish := &pb.Dish{}

	query := `insert into dishes (id, kithen_id,name, description,price, category,allergens, nutrition_info, dietary_info,ingredients,available) values ($1, $2, $3, $4, $5, $6, $7, $8,$9,$10,$11)`

	id := uuid.NewString()
	req.Dish.Id = id

	dish.CreatedAt = cast.ToString(time.Now())
	row := o.DB.QueryRowContext(ctx, query, req.Dish.Id, req.Dish.KitchenId, req.Dish.Name, req.Dish.Description, req.Dish.Price, req.Dish.Category, req.Dish.Allergens, req.Dish.NutritionInfo, req.Dish.DietaryInfo, req.Dish.Ingredients, req.Dish.Available)
	err = row.Scan(&dish.Id, &dish.KitchenId, &dish.Name, &dish.Description, &dish.Price, &dish.Category, pq.Array(req.Dish.Ingredients), pq.Array(req.Dish.Allergens), &req.Dish.NutritionInfo, &dish.Available, pq.Array(req.Dish.DietaryInfo))
	if err != nil {
		log.Error("error inserting dish", zap.Error(err))
		return resp, err
	}
	resp.Dish = dish
	log.Info("insert dish", zap.Any("dish", dish))
	return resp, nil
}

func (o *OrderRepository) UpdateDish(ctx context.Context, req *pb.UpdateDishRequest) (*pb.UpdateDishResponse, error) {

	log, err := l.NewLogger()
	if err != nil {
		return nil, err
	}

	query := `update dishes set name=$1, price=$2,available=$3 where id=$4 and deleted_at is null`

	resp := &pb.UpdateDishResponse{}
	dish := &pb.Dish{}

	dish.UpdatedAt = cast.ToString(time.Now())

	row := o.DB.QueryRowContext(ctx, query, req.Dish.Name, req.Dish.Price, req.Dish.Available, req.Dish.Id)
	err = row.Scan(&dish.Id, &dish.KitchenId, &dish.Name, &dish.Description, &dish.Price, &dish.Category, &dish.Ingredients, &dish.Available, &dish.CreatedAt)
	if err != nil {
		log.Error("error updating dish", zap.Error(err))
		return resp, err
	}
	resp.Dish = dish

	log.Info("updated dish", zap.Any("dish", dish))
	return resp, nil
}

func (o *OrderRepository) DeleteDish(ctx context.Context, req *pb.DeleteDishRequest) (*pb.DeleteDishResponse, error) {
	log, err := l.NewLogger()
	if err != nil {
		return nil, err
	}
	query := `update dishes set deleted_at=$1 where id=$2`

	resp := &pb.DeleteDishResponse{}

	_, err = o.DB.ExecContext(ctx, query, req.DishId)
	if err != nil {
		log.Error("error deleting dish", zap.Error(err))
		return resp, err
	}
	resp.Message = cast.ToString(gin.H{"message": "Dish successfully deleted"})

	log.Info("dish successfully deleted")
	return resp, nil
}

func (o *OrderRepository) GetDishes(ctx context.Context, req *pb.ListDishesRequest) (*pb.ListDishesResponse, error) {

	log, err := l.NewLogger()
	if err != nil {
		return nil, err
	}

	resp := &pb.ListDishesResponse{}

	var offset int32

	query := "select id, name, price, category, available from dishes where deleted_at is null order by id"
	if req.Page > 0 {
		offset = (req.Page * req.Limit) - req.Limit
	}
	query += " limit " + cast.ToString(req.Limit) + "offset " + cast.ToString(offset)

	var dishes []*pb.Dish

	rows, err := o.DB.QueryContext(ctx, query, req.KitchenId, req.Limit, req.Page)
	if err != nil {
		log.Error("error getting dishes", zap.Error(err))
		return resp, err
	}
	defer rows.Close()

	for rows.Next() {
		var dish pb.Dish

		err := rows.Scan(&dish.Id, &dish.Name, &dish.Price, &dish.Category, &dish.Available)
		if err != nil {
			log.Error("error getting dish", zap.Error(err))
			return resp, err
		}
		dishes = append(dishes, &dish)
	}
	resp.Dishes = dishes

	log.Info("Get Dish", zap.Any("dishes", dishes))
	return resp, nil
}

func (o *OrderRepository) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {

	log, err := l.NewLogger()
	if err != nil {
		return nil, err
	}

	query := "update dishes set status=$1 where id=$2"

	row := o.DB.QueryRowContext(ctx, query, req.OrderId, req.Status)
	err = row.Scan(&req.OrderId, &req.Status)
	if err != nil {
		log.Error("error updating dish status", zap.Error(err))
		return nil, err
	}
	resp := &pb.UpdateOrderStatusResponse{
		OrderId:   req.OrderId,
		Status:    req.Status,
		UpdatedAt: cast.ToString(time.Now()),
	}
	log.Info("update dish status", zap.String("status", req.Status))
	return resp, nil
}

func (o *OrderRepository) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	resp := &pb.CreateOrderResponse{}
	order := &pb.Order{}

	log, err := l.NewLogger()
	if err != nil {
		return nil, err
	}

	id := uuid.NewString()
	req.Order.Id = id

	query := "INSERT INTO orders (id, user_id, kitchen_id, delivery_address, delivery_time, status, created_at, total_amount) values ($1,$2,$3,$4,$5,$6,$7,$8)"

	row := o.DB.QueryRowContext(ctx, query, req.Order.Id, req.Order.UserId, req.Order.KitchenId, req.Order.DeliveryAddress, req.Order.DeliveryTime, req.Order.Status, req.Order.CreatedAt, req.Order.TotalAmount)

	err = row.Scan(&order.Id, &order.UserId, &order.KitchenId, &order.DeliveryAddress, &order.DeliveryTime, &order.Status, &order.CreatedAt, &order.TotalAmount)
	if err != nil {
		log.Error("error inserting order", zap.Error(err))
		return nil, err
	}
	resp.Order = order

	log.Info("insert order", zap.Any("order", order))
	return resp, nil

}

func (o *OrderRepository) listOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	log, err := l.NewLogger()
	if err != nil {
		return nil, err
	}

	resp := &pb.ListOrdersResponse{}
	kit := &pb.Kitchen{}

	var offset int32

	query := "select u.id, k.kitchen_name, u.total_amount, u.status, u.delivery_time from order u, kitchen k , where deleted_at is null order by id"
	if req.Page > 0 {
		offset = (req.Page * req.Limit) - req.Limit
	}
	query += " limit " + cast.ToString(req.Limit) + "offset " + cast.ToString(offset)

	var orders []*pb.Order

	rows, err := o.DB.QueryContext(ctx, query, req.UserId, req.KitchenId, req.Limit, req.Page)
	if err != nil {
		log.Error("error getting dishes", zap.Error(err))
		return resp, err
	}
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var order pb.Order

		err := rows.Scan(&order.Id, &kit.Name, &order.TotalAmount, &order.Status, &order.DeliveryTime)
		if err != nil {
			log.Error("error getting dishes", zap.Error(err))
			return resp, err
		}
		orders = append(orders, &order)
	}
	resp.Orders = orders

	log.Info("Get Order", zap.Any("Orders", orders))
	return resp, nil
}

func (o *OrderRepository) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {

	log, err := l.NewLogger()
	if err != nil {
		return nil, err
	}

	resp := &pb.GetOrderResponse{}

	query := "select d.id, u.user_name, d.total_amount, d.status, d.delivery_time from orders d, users u, where kitchen_id=$1"

	var offset int32

	if req.Page > 0 {
		offset = (req.Page * req.Limit) - req.Limit
	}
	query += " limit " + cast.ToString(req.Limit) + "offset " + cast.ToString(offset)

	var orders []*pb.Order

	rows, err := o.DB.QueryContext(ctx, query, req.KitchenID)
	if err != nil {
		log.Error("error getting dishes", zap.Error(err))
		return resp, err
	}
	defer rows.Close()

	for rows.Next() {
		var order pb.Order

		err = rows.Scan(&order.Id, &order.KitchenId, &order.TotalAmount, &order.Status, &order.DeliveryTime)
		if err != nil {
			log.Error("error getting dishes", zap.Error(err))
			return resp, err
		}
		orders = append(orders, &order)
	}
	resp.Order = orders
	log.Info("Get Order", zap.Any("Orders", orders))
	if err != nil {
		panic(err)
	}

	return resp, nil
}

func (o *OrderRepository) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*pb.CreateReviewResponse, error) {
	log, err := l.NewLogger()
	if err != nil {
		return nil, err
	}
	resp := &pb.CreateReviewResponse{}

	query := "insert into reviews(id, order_id, user_id, kitchen_id,rating, comment, created_at) values ($1,$2,$3,$4,$5,$6,$7)"
	_, err = o.DB.ExecContext(ctx, query)
	if err != nil {
		log.Error("error inserting review", zap.Error(err))
		return nil, err
	}
	return resp, nil
}

func (o *OrderRepository) UpdateDishNutritionInfo(ctx context.Context, req *pb.UpdateDishNutritionInfoRequest) (*pb.UpdateDishNutritionInfoResponse, error) {

	log, err := l.NewLogger()
	if err != nil {
		return nil, err
	}
	query := `UPDATE dishes SET allergens = $1, calories = $2, protein = $3, carbohydrates = $4, fat = $5, dietary_info = $6 WHERE id = $7`
	_, err = o.DB.ExecContext(ctx, query, pq.Array(req.Allergens), req.Calories, req.Protein, req.Carbohydrates, req.Fat, pq.Array(req.DietaryInfo), req.DishId)
	if err != nil {
		log.Error("error updating dishes", zap.Error(err))
		return nil, err
	}

	_, err = getDish(o.DB, req.DishId)
	if err != nil {
		log.Error("error getting dish", zap.Error(err))
		return nil, err
	}
	if err != nil {
		panic(err)
	}

	return &pb.UpdateDishNutritionInfoResponse{}, nil
}

func (o *OrderRepository) ListReviews(ctx context.Context, name string, req *pb.ListReviewsRequest) (*pb.ListReviewsResponse, error) {
	log, err := l.NewLogger()
	if err != nil {
		return nil, err
	}

	resp := &pb.ListReviewsResponse{}

	query := "SELECT  " +
		"reviews.id," +
		"users.username," +
		"reviews.rating," +
		"reviews.comment, " +
		"reviews.created_at" +
		"FROM reviews JOIN  users ON reviews.user_id = users.id" +
		"WHERE reviews.kitchen_id = $1'" +
		"ORDER BY reviews.created_at DESC "

	var offset int32

	if req.Page > 0 {
		offset = (req.Page * req.Limit) - req.Limit
	}
	query += " limit " + cast.ToString(req.Limit) + "offset " + cast.ToString(offset)

	var reviews []*pb.Review

	rows, err := o.DB.QueryContext(ctx, query, req.KitchenId, req.Limit, req.Page)
	if err != nil {
		log.Error("error getting reviews", zap.Error(err))
		return resp, err
	}
	defer rows.Close()

	for rows.Next() {
		var review pb.Review

		err = rows.Scan(&review.Id, &review.KitchenId, &review.OrderId, &review.UserId, &review.Rating, &review.Comment, &review.CreatedAt)
		if err != nil {
			log.Error("error getting reviews", zap.Error(err))
			return resp, err
		}

		reviews = append(reviews, &review)

	}

	log.Info("List Reviews", zap.Any("Reviews", reviews))
	return resp, nil
}

func (o *OrderRepository) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	log, err := l.NewLogger()
	if err != nil {
		return nil, err
	}

	if req.Payment == nil {
		return nil, errors.New("invalid payment request")
	}
	if req.Payment.OrderId == "" || req.Payment.PaymentMethod == "" {
		return nil, errors.New("order_id and payment_method are required")
	}

	txID, err := ProcessPayment(req.Payment.OrderId, req.Payment.PaymentMethod, req.Payment.CardNumber, "12/24", "123", req.Payment.Amount)
	if err != nil {
		return nil, err
	}

	paymentID := uuid.New().String()
	payment := &pb.Payment{
		Id:            paymentID,
		OrderId:       req.Payment.OrderId,
		Amount:        req.Payment.Amount,
		Status:        "success",
		PaymentMethod: req.Payment.PaymentMethod,
		TransactionId: txID,
		CreatedAt:     time.Now().Format(time.RFC3339),
	}

	req.Payment = payment

	resp := &pb.CreatePaymentResponse{
		Payment: payment,
	}

	log.Info("Payment created successfully:")
	return resp, nil
}

func (o *OrderRepository) GetDishRecommendations(ctx context.Context, req *pb.GetDishRecommendationsRequest) (*pb.GetDishRecommendationsResponse, error) {
	log, err := l.NewLogger()
	if err != nil {
		return nil, err
	}

	query := `SELECT id, name, description, price FROM dishes 
               WHERE user_id = $1 
               ORDER BY rating DESC`

	rows, err := o.DB.QueryContext(ctx, query, req.UserId)
	if err != nil {
		log.Error("error getting dish recommendations", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var dishes []*pb.Dish
	for rows.Next() {
		dish := &pb.Dish{}
		err := rows.Scan(&dish.Id, &dish.Name, &dish.Description, &dish.Price)
		if err != nil {
			log.Error("error scanning dish", zap.Error(err))
			return nil, err
		}
		dishes = append(dishes, dish)
	}

	log.Info("Dish recommendations retrieved successfully", zap.Int("count", len(dishes)))
	return &pb.GetDishRecommendationsResponse{Recommendations: dishes}, nil
}

func (o *OrderRepository) GetKitchenStatistics(ctx context.Context, req *pb.GetKitchenStatisticsRequest) (*pb.GetKitchenStatisticsResponse, error) {
	log, err := l.NewLogger()
	if err != nil {
		return nil, err
	}

	query := `SELECT COUNT(*) as total_orders, AVG(total_amount) as avg_order_value, 
               SUM(total_amount) as total_revenue 
               FROM orders 
               WHERE kitchen_id = $1 AND created_at >= $2 AND created_at <= $3`

	var stats pb.GetKitchenStatisticsResponse
	err = o.DB.QueryRowContext(ctx, query, req.KitchenId, req.StartDate, req.EndDate).Scan(
		&stats.TotalOrders, &stats.AverageRating, &stats.TotalRevenue)
	if err != nil {
		log.Error("error getting kitchen statistics", zap.Error(err))
		return nil, err
	}

	log.Info("Kitchen statistics retrieved successfully", zap.Any("stats", stats))
	return &stats, nil
}

func (o *OrderRepository) GetUserActivity(ctx context.Context, req *pb.GetUserActivityRequest) (*pb.GetUserActivityResponse, error) {
	log, err := l.NewLogger()
	if err != nil {
		return nil, err
	}

	query := `SELECT o.id, o.total_amount, o.status, o.created_at, k.name as kitchen_name 
               FROM orders o 
               JOIN kitchens k ON o.kitchen_id = k.id 
               WHERE o.user_id = $1 
               ORDER BY o.created_at DESC`

	rows, err := o.DB.QueryContext(ctx, query, req.UserId)
	if err != nil {
		log.Error("error getting user activity", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var activities []*pb.UserActivity
	for rows.Next() {
		activity := &pb.UserActivity{}
		err := rows.Scan(&activity.OrderId, &activity.Amount, &activity.Status, &activity.CreatedAt, &activity.KitchenName)
		if err != nil {
			log.Error("error scanning user activity", zap.Error(err))
			return nil, err
		}
		activities = append(activities, activity)
	}

	log.Info("User activity retrieved successfully", zap.Int("count", len(activities)))
	return &pb.GetUserActivityResponse{UserActivity: activities}, nil
}

func (o *OrderRepository) UpdateWorkingHours(ctx context.Context, req *pb.UpdateWorkingHoursRequest) (*pb.UpdateWorkingHoursResponse, error) {
	log, err := l.NewLogger()
	if err != nil {
		return nil, err
	}

	query := `UPDATE kitchens SET working_hours = $1 WHERE id = $2`
	_, err = o.DB.ExecContext(ctx, query, pq.Array(req.WorkingHours), req.KitchenId)
	if err != nil {
		log.Error("error updating working hours", zap.Error(err))
		return nil, err
	}

	kitchen, err := getDish(o.DB, req.KitchenId)
	if err != nil {
		log.Error("error getting kitchen", zap.Error(err))
		return nil, err
	}

	log.Info("Working hours updated successfully", zap.Any("kitchen", kitchen))
	return &pb.UpdateWorkingHoursResponse{WorkingHours: nil}, nil
}

func getDish(db *sql.DB, id string) (*pb.Dish, error) {
	var dish pb.Dish
	var ingredients, allergens, dietaryInfo []string
	err := db.QueryRow(`SELECT id, kitchen_id, name, description, price, category, ingredients, allergens, nutrition_info, dietary_info, available, created_at, updated_at FROM dishes WHERE id = $1`, id).
		Scan(&dish.Id, &dish.KitchenId, &dish.Name, &dish.Description, &dish.Price, &dish.Category, pq.Array(&ingredients), pq.Array(&allergens), &dish.NutritionInfo, pq.Array(&dietaryInfo), &dish.Available, &dish.CreatedAt, &dish.UpdatedAt)
	if err != nil {
		return nil, err
	}
	dish.Ingredients = ingredients
	dish.Allergens = allergens
	dish.DietaryInfo = dietaryInfo
	return &dish, nil
}

func ProcessPayment(orderID, paymentMethod, cardNumber, expiryDate, cvv string, amount float64) (string, error) {
	if cardNumber == "4111111111111111" && cvv == "123" {
		return "tx789", nil
	}
	return "", errors.New("payment failed")
}
