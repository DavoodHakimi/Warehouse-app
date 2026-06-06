package users

import "strconv"

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
		return &allUsers, err
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
		return &UserInfoResponse{}, nil
	}
	return &UserInfoResponse{
		FullName:    user.FullName,
		UserName:    user.UserName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		UserTypeID:  int(user.UserTypeID),
	}, err
}

func (s *Service) CreateUser(u *CreateUserRequest) error {
	return nil
}

// func (s *Service)
// func (s *Service)
