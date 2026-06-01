package auth

type SignUpRequest struct {
	CompanyName          string `json:"company_name" binding:"required,min=4,max=100"`
	FullName             string `json:"full_name" binding:"required,min=4,max=100"`
	UserName             string `json:"user_name" binding:"required,min=5,max=16"`
	Password             string `json:"password" binding:"required,min=8"`
	PasswordConfirmation string `json:"password_confirmation" binding:"eqfield=Password"`
	PhoneNumber          string `json:"phone_number" form:"phone_number" binding:"required,numeric,len=11,startswith=09"`
	Email                string `json:"email" form:"email" binding:"required,min=8,max=32,email"`
}

type LogInRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type meResponse struct {
	Token    string `json:"token"`
	UserName string `json:"user_name"`
}
