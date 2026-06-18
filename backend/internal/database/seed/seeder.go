package seed

import (
	"context"
	"fmt"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/DavoodHakimi/warehouse-app/internal/users"
)

type Seeder struct {
	db *gorm.DB
}

func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

func (s *Seeder) Run(ctx context.Context) error {
	log.Println("Running database seeder...")

	if err := s.seedCurrencies(ctx); err != nil {
		return fmt.Errorf("seeding Currencies: %w", err)
	}
	if err := s.seedPartnerTypes(ctx); err != nil {
		return fmt.Errorf("seeding Partner types: %w", err)
	}
	if err := s.seedUserTypes(ctx); err != nil {
		return fmt.Errorf("seeding user types: %w", err)
	}
	if err := s.seedPermissionTypes(ctx); err != nil {
		return fmt.Errorf("seeding permission types: %w", err)
	}
	if err := s.seedPermissions(ctx); err != nil {
		return fmt.Errorf("seeding permissions: %w", err)
	}

	log.Println("Seeding complete.")
	return nil
}

func (s *Seeder) seedCurrencies(ctx context.Context) error {
	records := make([]Currency, len(currencies))
	for _, currency := range currencies {
		records = append(records, currency)
	}

	result := s.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&records)

	if result.Error != nil {
		return fmt.Errorf("seeding Currencies: %w", result.Error)
	}
	log.Printf("Currencies seeded (%d)", len(records))
	return nil
}

func (s *Seeder) seedPartnerTypes(ctx context.Context) error {
	records := make([]BusinessPartnerType, len(partnerTypes))
	for _, partnerType := range partnerTypes {
		records = append(records, partnerType)
	}
	result := s.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&records)

	if result.Error != nil {
		return fmt.Errorf("seeding user types: %w", result.Error)
	}
	log.Printf("Partner types seeded (%d)", len(records))
	return nil
}

func (s *Seeder) seedUserTypes(ctx context.Context) error {
	records := make([]users.UserType, len(userTypes))
	for i, ut := range userTypes {
		records[i] = users.UserType{
			Model:       gorm.Model{ID: ut.ID},
			Name:        ut.Name,
			PersianName: ut.PersianName,
			Description: ut.Description,
		}
	}

	result := s.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&records)

	if result.Error != nil {
		return fmt.Errorf("seeding user types: %w", result.Error)
	}

	log.Printf("User types seeded (%d)", len(records))
	return nil
}

func (s *Seeder) seedPermissionTypes(ctx context.Context) error {
	var records []users.PermissionType
	for resource, actions := range permissionTypes {
		for _, action := range actions {
			records = append(records, users.PermissionType{
				Name: resource + "." + action,
			})
		}
	}

	result := s.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&records)

	if result.Error != nil {
		return fmt.Errorf("seeding permission types: %w", result.Error)
	}

	log.Printf("Permission types seeded (%d)", len(records))
	return nil
}

func (s *Seeder) seedPermissions(ctx context.Context) error {
	var permTypes []users.PermissionType
	if err := s.db.WithContext(ctx).Find(&permTypes).Error; err != nil {
		return fmt.Errorf("loading permission types: %w", err)
	}

	permTypeIDs := make(map[string]uint, len(permTypes))
	for _, pt := range permTypes {
		permTypeIDs[pt.Name] = pt.ID
	}

	var toInsert []users.Permission
	for roleID, perms := range rolePermissions {
		for _, permName := range perms {
			ptID, ok := permTypeIDs[permName]
			if !ok {
				return fmt.Errorf("unknown permission type %q — did seedPermissionTypes run?", permName)
			}
			toInsert = append(toInsert, users.Permission{
				UserTypeID:       uint(roleID),
				PermissionTypeID: ptID,
			})
		}
	}

	result := s.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&toInsert)

	if result.Error != nil {
		return fmt.Errorf("inserting permissions: %w", result.Error)
	}

	log.Printf("Permissions seeded (%d links)", len(toInsert))
	return nil
}
