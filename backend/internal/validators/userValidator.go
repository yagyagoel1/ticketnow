package validators

type CreateUserReq struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
}

type SigninUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateProfile struct {
	Name string `json:"name" validate:"required"`
}

type PostBooking struct {
	ShowId      uint                   `json:"showId" validate:"required"`
	TicketTypes map[string]interface{} `json:"ticketTypes" validate:"required"`
}

type UpdatePassword struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}
