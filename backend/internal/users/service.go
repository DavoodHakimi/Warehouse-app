package users

import (
	"errors"
	"strconv"

	"github.com/DavoodHakimi/warehouse-app/internal/audit"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AllUsers(cID int) (*UsersInfo, error) {
	users, err := s.repo.ReadCompanyUsers(cID)
	allUsers := UsersInfo{Users: make([]UserInfoResponse, 0, len(users))}

	if err != nil {
		return nil, err
	}

	for _, item := range users {
		allUsers.Users = append(allUsers.Users, UserInfoResponse{
			FullName:    item.FullName,
			UserName:    item.UserName,
			Email:       item.Email,
			PhoneNumber: item.PhoneNumber,
			UserTypeID:  int(item.UserTypeID),
		})
	}
	return &allUsers, nil
}

func (s *Service) ReadUser(userID string) (*UserInfoResponse, error) {
	val, _ := strconv.Atoi(userID)
	user, err := s.repo.FindByID(val)
	if err != nil {
		return nil, err
	}
	return &UserInfoResponse{
		ID:          int(user.ID),
		FullName:    user.FullName,
		UserName:    user.UserName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		UserTypeID:  int(user.UserTypeID),
	}, err
}

func (s *Service) CreateUser(u *CreateUserRequest, cid int) error {
	hashedPassword, err := HashPassword(u.Password)
	if err != nil {
		return err
	}

	user := User{
		FullName:    u.FullName,
		UserName:    u.UserName,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		UserTypeID:  uint(u.UserTypeID),
		CompanyID:   uint(cid),
		Password:    string(hashedPassword),
	}
	return s.repo.Create(&user)
}

func (s *Service) UpdateUser(u *UpdateUserRequest, userRequestedID int) error {
	changedFields := s.modifiedFields(u)
	if len(changedFields) == 0 {
		return errors.New("no changes detected")
	}

	user := &User{
		FullName:    u.FullName,
		UserName:    u.UserName,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		UserTypeID:  uint(u.UserTypeID),
	}
	user.ID = uint(u.ID)

	err := s.repo.Update(user)
	if err != nil {
		return err
	}
	for field, values := range changedFields {
		log := audit.Log{
			EntityType: "user",
			EntityID:   uint(u.ID),
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

func (s *Service) DeleteUser(uID int) error {
	user, err := s.repo.FindByID(uID)
	if err != nil {
		return err
	}
	return s.repo.Delete(user)
}

func (s *Service) modifiedFields(u *UpdateUserRequest) map[string][2]string {
	oldValues, err := s.repo.FindByID(u.ID)
	if err != nil {
		return nil
	}

	changes := make(map[string][2]string)

	if u.FullName != "" && u.FullName != oldValues.FullName {
		changes["FullName"] = [2]string{oldValues.FullName, u.FullName}
	}
	if u.Email != "" && u.Email != oldValues.Email {
		changes["Email"] = [2]string{oldValues.Email, u.Email}
	}
	if u.PhoneNumber != "" && u.PhoneNumber != oldValues.PhoneNumber {
		changes["PhoneNumber"] = [2]string{oldValues.PhoneNumber, u.PhoneNumber}
	}
	if u.UserTypeID != int(oldValues.UserTypeID) {
		changes["UserTypeID"] = [2]string{strconv.Itoa(int(oldValues.UserTypeID)), strconv.Itoa(u.UserTypeID)}
	}
	return changes
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
