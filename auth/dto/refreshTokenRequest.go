package dto

type RefreshTokenRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r RefreshTokenRequest) RefreshTokenExists() bool {
	return r.RefreshToken != ""
}
