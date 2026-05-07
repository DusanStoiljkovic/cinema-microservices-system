package service

import (
	"booking-service/internal/dto"
	"booking-service/internal/mapper"
	"booking-service/internal/models"
	"booking-service/internal/utils"
	"context"
	"errors"
	"fmt"
)

type OrderRepository interface {
	GetAll(ctx context.Context) ([]*models.Order, error)
	GetByID(ctx context.Context, id uint) (*models.Order, error)
	GetByUserID(ctx context.Context, id uint) ([]*models.Order, error)
	Create(ctx context.Context, order *models.Order) (*models.Order, error)
	Pay(ctx context.Context, id uint) error
	Cancel(ctx context.Context, id uint) error
	Delete(ctx context.Context, id uint) error
}

type OrderService struct {
	repo OrderRepository
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (service *OrderService) GetAllOrders(ctx context.Context) ([]*models.Order, error) {
	return service.repo.GetAll(ctx)
}

func (service *OrderService) GetOrderByID(ctx context.Context, id uint) (*models.Order, error) {
	return service.repo.GetByID(ctx, id)
}

func (service *OrderService) GetOrdersByUserID(ctx context.Context, id uint) ([]*models.Order, error) {
	return service.repo.GetByUserID(ctx, id)
}

func (service *OrderService) GetMyOrders(ctx context.Context) ([]*models.Order, error) {
	val := ctx.Value("userID")
	userID, ok := val.(uint)
	if !ok {
		return nil, utils.NewAuthFailed("Unauthorized", errors.New("OrderService.GetMyOrders -> userID is not available/valid"))
	}

	return service.repo.GetByUserID(ctx, userID)
}

func (service *OrderService) CreateOrder(ctx context.Context, orderReq *dto.OrderRequest) (*models.Order, error) {
	val := ctx.Value("userID")
	userID, ok := val.(uint)
	if !ok {
		return nil, utils.NewAuthFailed("Unauthorized", errors.New(fmt.Sprintf("OrderService.CreateOrder -> userID is not available/valid, userID=%d", userID)))
	}

	order := mapper.OrderFromRequest(orderReq)
	order.UserID = userID

	return service.repo.Create(ctx, order)
}

func (service *OrderService) PayOrder(ctx context.Context, id uint) error {
	return service.repo.Pay(ctx, id)
}

func (service *OrderService) CancelOrder(ctx context.Context, id uint) error {
	return service.repo.Cancel(ctx, id)
}

func (service *OrderService) DeleteOrder(ctx context.Context, id uint) error {
	return service.repo.Delete(ctx, id)
}
