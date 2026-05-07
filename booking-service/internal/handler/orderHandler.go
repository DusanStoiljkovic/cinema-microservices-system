package handler

import (
	"booking-service/internal/dto"
	"booking-service/internal/mapper"
	"booking-service/internal/models"
	"booking-service/internal/utils"
	"context"
	"encoding/json"
	"net/http"
)

type OrderService interface {
	GetAllOrders(ctx context.Context) ([]*models.Order, error)
	GetOrderByID(ctx context.Context, id uint) (*models.Order, error)
	GetOrdersByUserID(ctx context.Context, id uint) ([]*models.Order, error)
	GetMyOrders(ctx context.Context) ([]*models.Order, error)
	CreateOrder(ctx context.Context, orderReq *dto.OrderRequest) (*models.Order, error)
	PayOrder(ctx context.Context, id uint) error
	CancelOrder(ctx context.Context, id uint) error
	DeleteOrder(ctx context.Context, id uint) error
}

type OrderHandler struct {
	service OrderService
}

func NewOrderHandler(service OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (handler *OrderHandler) HandleGetAllOrders(w http.ResponseWriter, r *http.Request) error {
	orders, err := handler.service.GetAllOrders(r.Context())
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, mapper.OrdersToResponse(orders))
}

func (handler *OrderHandler) HandleGetOrderByID(w http.ResponseWriter, r *http.Request) error {
	id, err := parseParamID(r, "id")
	if err != nil {
		return utils.NewInvalidInput("Invalid order id", err)
	}

	order, err := handler.service.GetOrderByID(r.Context(), id)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, order)
}

func (handler *OrderHandler) HandleGetOrdersByUserID(w http.ResponseWriter, r *http.Request) error {
	id, err := parseParamID(r, "id")
	if err != nil {
		return utils.NewInvalidInput("Invalid user id", err)
	}

	orders, err := handler.service.GetOrdersByUserID(r.Context(), id)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, mapper.OrdersToResponse(orders))
}

func (handler *OrderHandler) HandleGetMyOrders(w http.ResponseWriter, r *http.Request) error {
	orders, err := handler.service.GetMyOrders(r.Context())
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, mapper.OrdersToResponse(orders))
}

func (handler *OrderHandler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) error {
	req := &dto.OrderRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.NewInvalidInput("Invalid request body", err)
	}

	order, err := handler.service.CreateOrder(r.Context(), req)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, mapper.OrderToResponse(order))
}

func (handler *OrderHandler) HandlePayOrder(w http.ResponseWriter, r *http.Request) error {
	id, err := parseParamID(r, "id")
	if err != nil {
		return utils.NewInvalidInput("Invalid order id", err)
	}

	err = handler.service.PayOrder(r.Context(), id)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func (handler *OrderHandler) HandleCancelOrder(w http.ResponseWriter, r *http.Request) error {
	id, err := parseParamID(r, "id")
	if err != nil {
		return utils.NewInvalidInput("Invalid order id", err)
	}

	err = handler.service.CancelOrder(r.Context(), id)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func (handler *OrderHandler) HandleDeleteOrder(w http.ResponseWriter, r *http.Request) error {
	id, err := parseParamID(r, "id")
	if err != nil {
		return utils.NewInvalidInput("Invalid order id", err)
	}

	err = handler.service.DeleteOrder(r.Context(), id)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
