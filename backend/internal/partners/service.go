package partners

import (
	"errors"
	"strconv"

	"github.com/DavoodHakimi/warehouse-app/internal/audit"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AllPartners(cID int) (*PartnersInfo, error) {
	partners, err := s.repo.ReadCompanyPartners(cID)
	allPartners := PartnersInfo{Partners: make([]PartnerInfoResponse, 0, len(partners))}

	if err != nil {
		return &allPartners, err
	}

	for _, item := range partners {
		allPartners.Partners = append(allPartners.Partners, PartnerInfoResponse{
			Name:                item.Name,
			BusinessPartnerType: item.BusinessPartnerType.Name,
			Email:               item.Email,
			PhoneNumber:         item.PhoneNumber,
			ContactName:         item.ContactName,
			ContactPhoneNumber:  item.ContactPhoneNumber,
		})
	}
	return &allPartners, nil
}

func (s *Service) ReadPartner(partnerID string) (*PartnerInfoResponse, error) {
	val, _ := strconv.Atoi(partnerID)
	partner, err := s.repo.FindByID(val)
	if err != nil {
		return &PartnerInfoResponse{}, nil
	}
	return &PartnerInfoResponse{
		ID:                  int(partner.ID),
		Name:                partner.Name,
		BusinessPartnerType: partner.BusinessPartnerType.Name,
		Email:               partner.Email,
		PhoneNumber:         partner.PhoneNumber,
		ContactName:         partner.ContactName,
		ContactPhoneNumber:  partner.ContactPhoneNumber,
	}, err
}

func (s *Service) CreatePartner(p *CreatePartnerRequest, cid int) error {

	partner := BusinessPartner{
		Name:                  p.Name,
		BusinessPartnerTypeID: p.BusinessPartnerTypeID,
		Email:                 p.Email,
		PhoneNumber:           p.PhoneNumber,
		CompanyID:             uint(cid),
		ContactName:           p.ContactName,
		ContactPhoneNumber:    p.ContactPhoneNumber,
	}
	return s.repo.Create(&partner)
}

func (s *Service) UpdatePartner(p *UpdatePartnerRequest, userRequestedID int) error {
	changedFields := s.modifiedFields(p)
	if len(changedFields) == 0 {
		return errors.New("no changes detected")
	}

	partner := &BusinessPartner{
		Name:                  p.Name,
		Email:                 p.Email,
		PhoneNumber:           p.PhoneNumber,
		BusinessPartnerTypeID: p.BusinessPartnerTypeID,
		ContactName:           p.ContactName,
		ContactPhoneNumber:    p.ContactPhoneNumber,
	}
	partner.ID = uint(p.ID)

	err := s.repo.Update(partner)
	if err != nil {
		return err
	}
	for field, values := range changedFields {
		log := audit.Log{
			EntityType: "partner",
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

func (s *Service) DeletePartner(pID int) error {
	partner, err := s.repo.FindByID(pID)
	if err != nil {
		return err
	}
	return s.repo.Delete(partner)
}

func (s *Service) modifiedFields(p *UpdatePartnerRequest) map[string][2]string {
	oldValues, err := s.repo.FindByID(p.ID)
	if err != nil {
		return nil
	}

	changes := make(map[string][2]string)

	if p.Name != "" && p.Name != oldValues.Name {
		changes["Name"] = [2]string{oldValues.Name, p.Name}
	}
	if p.Email != "" && p.Email != oldValues.Email {
		changes["Email"] = [2]string{oldValues.Email, p.Email}
	}
	if p.PhoneNumber != "" && p.PhoneNumber != oldValues.PhoneNumber {
		changes["PhoneNumber"] = [2]string{oldValues.PhoneNumber, p.PhoneNumber}
	}
	if p.BusinessPartnerTypeID != oldValues.BusinessPartnerTypeID {
		changes["BusinessPartnerTypeID"] = [2]string{strconv.Itoa(int(oldValues.BusinessPartnerTypeID)), strconv.Itoa(int(p.BusinessPartnerTypeID))}

	}
	return changes
}
