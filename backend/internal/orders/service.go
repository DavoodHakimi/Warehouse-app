package orders

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/DavoodHakimi/warehouse-app/internal/audit"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AllOrders(cID int) (*OrdersInfo, error) {
	orders, err := s.repo.ReadCompanyOrders(cID)
	allOrders := OrdersInfo{Orders: make([]OrderInfoResponse, 0, len(orders))}
	if err != nil {
		return nil, err
	}

	for _, item := range orders {
		allOrders.Orders = append(allOrders.Orders, OrderInfoResponse{
			ID:                  uint(item.ID),
			OrderType:           item.OrderType,
			OrderNumber:         item.OrderNumber,
			Status:              item.Status,
			BusinessPartnerName: item.BusinessPartner.Name,
			Currency:            item.Currency.Name,
			ExchangeRate:        item.ExchangeRate,
		})
	}
	return &allOrders, nil
}

func (s *Service) ReadOrder(orderID string, companyID int) (*OrderInfoResponse, error) {
	order, _, err := s.CheckOrderExist(orderID, companyID)
	if err != nil {
		return nil, err
	}

	return &OrderInfoResponse{
		ID:                  order.ID,
		OrderType:           order.OrderType,
		OrderNumber:         order.OrderNumber,
		Status:              order.Status,
		BusinessPartnerName: order.BusinessPartner.Name,
		Currency:            order.Currency.Name,
		ExchangeRate:        order.ExchangeRate}, nil
}

func (s *Service) CreateOrder(o *CreateOrderRequest, cid int) error {
	allItems := make([]OrderItem, 0, len(o.OrderItems))
	for _, item := range o.OrderItems {
		allItems = append(allItems, OrderItem{
			ProductID:    item.ProductID,
			Quantity:     item.Quantity,
			PerItemPrice: item.PerItemPrice,
		})
	}
	order := Order{
		OrderType:         o.OrderType,
		OrderNumber:       fmt.Sprintf("ORD-%d", time.Now().UnixNano()),
		CompanyID:         uint(cid),
		Status:            "Pending",
		BusinessPartnerID: o.BusinessPartnerID,
		CurrencyID:        o.CurrencyID,
		ExchangeRate:      o.ExchangeRate,
		OrderItems:        allItems,
	}
	return s.repo.Create(&order)
}

func (s *Service) UpdateOrder(o *UpdateOrderRequest, userRequestedID int, companyID int) error {
	changedFields := s.modifiedFields(o, companyID)
	if len(changedFields) == 0 {
		return errors.New("no changes detected")
	}

	order := Order{
		OrderType:         o.OrderType,
		BusinessPartnerID: o.BusinessPartnerID,
		CurrencyID:        o.CurrencyID,
		ExchangeRate:      o.ExchangeRate,
	}
	order.ID = uint(o.ID)

	existing, err := s.repo.FindByID(order.ID, uint(companyID))
	if err != nil {
		return err
	}
	order.Status = existing.Status
	err = s.repo.Update(&order, uint(companyID))
	if err != nil {
		return err
	}
	for field, values := range changedFields {
		log := audit.Log{
			EntityType: "order",
			EntityID:   uint(o.ID),
			Event:      "updated",
			Field:      field,
			OldValue:   values[0],
			NewValue:   values[1],
			ByUserID:   uint(userRequestedID),
		}
		s.repo.RecordAudit(&log)
	}
	return nil
}
func (s *Service) modifiedFields(o *UpdateOrderRequest, cid int) map[string][2]string {
	oldValues, err := s.repo.FindByID(o.ID, uint(cid))
	if err != nil {
		return nil
	}

	changes := make(map[string][2]string)

	if o.OrderType != "" && o.OrderType != oldValues.OrderType {
		changes["OrderType"] = [2]string{oldValues.OrderType, o.OrderType}
	}
	if o.BusinessPartnerID != oldValues.BusinessPartnerID {
		changes["BusinessPartnerID"] = [2]string{strconv.Itoa(int(oldValues.BusinessPartnerID)), strconv.Itoa(int(o.BusinessPartnerID))}
	}
	if o.CurrencyID != oldValues.CurrencyID {
		changes["CurrencyID"] = [2]string{strconv.Itoa(int(oldValues.CurrencyID)), strconv.Itoa(int(o.CurrencyID))}
	}
	if o.ExchangeRate != oldValues.ExchangeRate {
		changes["ExchangeRate"] = [2]string{strconv.FormatFloat(oldValues.ExchangeRate, 'f', -1, 64), strconv.FormatFloat(o.ExchangeRate, 'f', -1, 64)}
	}
	return changes
}

func (s *Service) DeleteOrder(orderID uint, companyID int) error {

	order, err := s.repo.FindByID(orderID, uint(companyID))
	if err != nil {
		return errors.New("Order Not found")
	}

	return s.repo.Delete(order)
}

func (s *Service) Approve(orderID string, companyID int) error {
	order, val, err := s.CheckOrderExist(orderID, companyID)
	if err != nil {
		return err
	}

	if order.Status != "Pending" {
		return errors.New("This Order status can not changed to approve.")
	}

	if order.OrderType == "sale" {
		return s.repo.ApproveSaleTransaction(uint(val), order.OrderItems)
	}
	return s.repo.StatusUpdate(uint(val), "Approved")
}

func (s *Service) Pack(orderID string, companyID int) error {
	order, val, err := s.CheckOrderExist(orderID, companyID)
	if err != nil {
		return err
	}

	if order.Status != "Approved" || order.OrderType != "sale" {
		return errors.New("This Order status can not changed to Packed.")
	}
	return s.repo.StatusUpdate(uint(val), "Packed")
}

func (s *Service) Ship(orderID string, companyID int) error {
	order, val, err := s.CheckOrderExist(orderID, companyID)
	if err != nil {
		return err
	}

	if order.Status != "Packed" || order.OrderType != "sale" {
		return errors.New("This Order status can not changed to Shipped.")
	}

	return s.repo.ShipSaleTransaction(uint(val), order.OrderItems)
}

func (s *Service) MarkWaiting(orderID string, companyID int) error {
	order, val, err := s.CheckOrderExist(orderID, companyID)
	if err != nil {
		return err
	}

	if order.Status != "Approved" || order.OrderType != "purchase" {
		return errors.New("This Order status can not changed to Waiting.")
	}
	return s.repo.StatusUpdate(uint(val), "Waiting")
}

func (s *Service) Receive(orderID string, companyID int) error {
	order, val, err := s.CheckOrderExist(orderID, companyID)
	if err != nil {
		return err
	}

	if order.Status != "Waiting" || order.OrderType != "purchase" {
		return errors.New("This Order status can not changed to Received.")
	}

	return s.repo.ReceivePurchaseTransaction(uint(val), order.OrderItems)
}

func (s *Service) Cancel(orderID string, companyID int) error {
	order, val, err := s.CheckOrderExist(orderID, companyID)
	if err != nil {
		return err
	}

	if order.Status == "Canceled" {
		return errors.New("This Order is already canceled.")
	}

	reserved := order.OrderType == "sale" && (order.Status == "Approved" || order.Status == "Packed")
	if !reserved {
		return s.repo.StatusUpdate(uint(val), "Canceled")
	}

	return s.repo.CancelSaleTransaction(uint(val), order.OrderItems)
}

func (s *Service) CheckOrderExist(orderID string, companyID int) (*Order, int, error) {
	val, err := strconv.Atoi(orderID)
	if err != nil {
		return nil, val, errors.New("Invalid Order ID")
	}
	order, err := s.repo.FindByID(uint(val), uint(companyID))
	if err != nil {
		return nil, val, err
	}
	return order, val, nil
}
