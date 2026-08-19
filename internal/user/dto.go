package user

type Response struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func ToResponse(user *User) Response {
	return Response{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role.Name,
	}
}