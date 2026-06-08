package partners

type PartnersInfo struct {
	Partners []PartnerInfoResponse `json:"partners"`
}

type PartnerInfoResponse struct {
	ID                  int    `json:"id"` // do not show this in frontend, BUT STORE IT TO USE IN UPDATE
	Name                string `json:"name"`
	BusinessPartnerType string `json:"business_partner_type"`
	PhoneNumber         string `json:"phone_number" form:"phone_number"`
	Email               string `json:"email" form:"email"`
	ContactName         string `json:"contact_name"`
	ContactPhoneNumber  string `json:"contact_phone_number"`
}

type CreatePartnerRequest struct {
	Name                  string `json:"name" binding:"required,min=4,max=100"`
	BusinessPartnerTypeID uint   `json:"business_partner_type_id" binding:"required"`
	PhoneNumber           string `json:"phone_number" form:"phone_number" binding:"required,numeric,len=11,startswith=09"`
	Email                 string `json:"email" form:"email" binding:"required,min=8,max=32,email"`
	ContactName           string `json:"contact_name"`
	ContactPhoneNumber    string `json:"contact_phone_number"`
}

type UpdatePartnerRequest struct {
	ID                    int    `json:"id" binding:"required,numeric"`
	Name                  string `json:"name" binding:"required,min=4,max=100"`
	BusinessPartnerTypeID uint   `json:"business_partner_type_id"`
	PhoneNumber           string `json:"phone_number" form:"phone_number" binding:"required,numeric,len=11,startswith=09"`
	Email                 string `json:"email" form:"email" binding:"required,min=8,max=32,email"`
	ContactName           string `json:"contact_name,omitempty"`
	ContactPhoneNumber    string `json:"contact_phone_number,omitempty"`
}
