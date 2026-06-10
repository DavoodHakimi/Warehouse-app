package orders

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AllOrders(cID int) (*OrdersInfo, error) {
	return nil, nil
}

func (s *Service) ReadOrder(orderID string) (*OrderInfoResponse, error) {
	return nil, nil
}

func (s *Service) CreateOrder(o *CreateOrderRequest, cid int) error {
	return nil
}

func (s *Service) UpdateOrder(o *UpdateOrderRequest, userRequestedID int) error {
	return nil
}

func (s *Service) DeleteOrder(oID int) error {
	return nil
}

func (s *Service) Approve() {

}

func (s *Service) Pack() {

}

func (s *Service) Ship() {

}

func (s *Service) Receive() {

}
