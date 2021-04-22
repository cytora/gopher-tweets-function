package model

type AuthenticatedUser struct {
	AccessToken  string
	AccessSecret string
}

type UserProfile struct {
	IsLoggedIn bool   `json:"isLoggedIn"`
	ScreenName string `json:"screenName"`
}
