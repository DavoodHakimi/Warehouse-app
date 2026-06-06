package users

type UsersInfo struct {
	Users []UserInfoResponse `json:"users"`
}

type UserInfoResponse struct {
	FullName    string `json:"full_name" binding:"required,min=4,max=100"`
	UserName    string `json:"user_name" binding:"required,min=5,max=16"`
	Password    string `json:"-" binding:"required,min=8"`
	UserTypeID  int    `json:"user_type_id" binding:"required,numeric"`
	PhoneNumber string `json:"phone_number" form:"phone_number" binding:"required,numeric,len=11,startswith=09"`
	Email       string `json:"email" form:"email" binding:"required,min=8,max=32,email"`
}

type CreateUserRequest struct {
	FullName             string `json:"full_name" binding:"required,min=4,max=100"`
	UserName             string `json:"user_name" binding:"required,min=5,max=16"`
	UserTypeID           int    `json:"user_type_id" binding:"required,numeric"`
	Password             string `json:"password" binding:"required,min=8"`
	PasswordConfirmation string `json:"password_confirmation" binding:"eqfield=Password"`
	PhoneNumber          string `json:"phone_number" form:"phone_number" binding:"required,numeric,len=11,startswith=09"`
	Email                string `json:"email" form:"email" binding:"required,min=8,max=32,email"`
}
