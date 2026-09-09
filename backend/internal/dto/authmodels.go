package dto

type UserLoginRequestParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponseParams struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	AccessToken string `json:"access_token"`
}

type RefreshResponseParams struct {
	AccessToken string `json:"access_token"`
}
