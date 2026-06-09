package seed

import (
	"context"
	"fmt"
	"log"

	"gorm.io/gorm"

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

func (s *Seeder) seedUserTypes(ctx context.Context) error {
	for _, ut := range userTypes {
		result := s.db.WithContext(ctx).
			Where(users.UserType{Name: ut.Name}).
			FirstOrCreate(&users.UserType{
				Name:        ut.Name,
				PersianName: ut.PersianName,
				Description: ut.Description,
			})

		if result.Error != nil {
			return fmt.Errorf("upserting user type %q: %w", ut.Name, result.Error)
		}
	}

	log.Printf("User types seeded (%d)", len(userTypes))
	return nil
}

func (s *Seeder) seedPermissionTypes(ctx context.Context) error {
	count := 0
	for resource, actions := range permissionTypes {
		for _, action := range actions {
			name := resource + "." + action

			result := s.db.WithContext(ctx).
				Where(users.PermissionType{Name: name}).
				FirstOrCreate(&users.PermissionType{Name: name})

			if result.Error != nil {
				return fmt.Errorf("upserting permission type %q: %w", name, result.Error)
			}
			count++
		}
	}

	log.Printf("Permission types seeded (%d)", count)
	return nil
}

func (s *Seeder) seedPermissions(ctx context.Context) error {
	count := 0
	for roleName, perms := range rolePermissions {

		var userType users.UserType
		if err := s.db.WithContext(ctx).
			Where("name = ?", roleName).
			First(&userType).Error; err != nil {
			return fmt.Errorf("finding user type %q: %w", roleName, err)
		}

		for _, permName := range perms {

			var permType users.PermissionType
			if err := s.db.WithContext(ctx).
				Where("name = ?", permName).
				First(&permType).Error; err != nil {
				return fmt.Errorf("finding permission type %q: %w", permName, err)
			}

			result := s.db.WithContext(ctx).
				Where(users.Permission{
					UserTypeID:       userType.ID,
					PermissionTypeID: permType.ID,
				}).
				FirstOrCreate(&users.Permission{
					UserTypeID:       userType.ID,
					PermissionTypeID: permType.ID,
				})

			if result.Error != nil {
				return fmt.Errorf(
					"linking %q -> %q: %w",
					roleName, permName, result.Error,
				)
			}
			count++
		}
	}

	log.Printf("Permissions seeded (%d links)", count)
	return nil
}
