package products

import (
	"errors"
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

func (s *Service) AllProducts(cID int) (*ProductsInfo, error) {
	prods, err := s.repo.ReadCompanyProducts(cID)
	allProducts := ProductsInfo{Products: make([]ProductInfoResponse, 0, len(prods))}

	if err != nil {
		return &allProducts, err
	}

	for _, item := range prods {
		allProducts.Products = append(allProducts.Products, ProductInfoResponse{
			Name:          item.Name,
			ProductNumber: item.ProductNumber,
			IsFrozen:      item.IsFrozen,
			DefaultPrice:  item.DefaultPrice,
		})
	}
	return &allProducts, nil
}

func (s *Service) ReadProduct(productNumber string) (*ProductInfoResponse, error) {
	prod, err := s.repo.FindByID(productNumber)
	if err != nil {
		return &ProductInfoResponse{}, nil
	}
	return &ProductInfoResponse{
		ID:            int(prod.ID),
		Name:          prod.Name,
		ProductNumber: prod.ProductNumber,
		IsFrozen:      prod.IsFrozen,
		DefaultPrice:  prod.DefaultPrice,
	}, err
}

func (s *Service) CreateProduct(u *ProductRequest, cid int) error {

	prod := Product{
		Name:          u.Name,
		ProductNumber: "PRD-" + strconv.Itoa(int(time.Now().Unix())),
		IsFrozen:      u.IsFrozen,
		DefaultPrice:  u.DefaultPrice,
		CompanyID:     uint(cid),
	}
	return s.repo.Create(&prod)
}

func (s *Service) UpdateProduct(p *UpdateProductRequest, userRequestedID int) error {
	changedFields := s.modifiedFields(p)
	if len(changedFields) == 0 {
		return errors.New("no changes detected")
	}

	prod := &Product{
		Name:         p.Name,
		IsFrozen:     p.IsFrozen,
		DefaultPrice: p.DefaultPrice,
	}
	prod.ID = uint(p.ID)

	err := s.repo.Update(prod)
	if err != nil {
		return err
	}
	for field, values := range changedFields {
		log := audit.Log{
			EntityType: "product",
			EntityID:   uint(p.ID),
			Event:      "updated",
			Field:      field,
			OldValue:   values[0],
			NewValue:   values[1],
			ByUserID:   uint(userRequestedID),
		}
		audit.Record(s.repo.db, &log)
	}
	return nil
}

func (s *Service) DeleteProduct(pID string) error {
	prod, err := s.repo.FindByID(pID)
	if err != nil {
		return err
	}
	return s.repo.Delete(prod)
}

func (s *Service) modifiedFields(p *UpdateProductRequest) map[string][2]string {
	oldValues, err := s.repo.FindByID(p.ProductNumber)
	if err != nil {
		return nil
	}

	changes := make(map[string][2]string)

	if p.Name != "" && p.Name != oldValues.Name {
		changes["Name"] = [2]string{oldValues.Name, p.Name}
	}
	if p.IsFrozen != oldValues.IsFrozen {
		changes["IsFrozen"] = [2]string{strconv.FormatBool(oldValues.IsFrozen), strconv.FormatBool(p.IsFrozen)}
	}
	if p.DefaultPrice != oldValues.DefaultPrice {
		changes["DefaultPrice"] = [2]string{strconv.FormatFloat(oldValues.DefaultPrice, 'f', -1, 64), strconv.FormatFloat(p.DefaultPrice, 'f', -1, 64)}
	}
	return changes
}
