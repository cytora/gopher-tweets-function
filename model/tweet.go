package model

type PostTweetRequest struct {
	Body string `json:"body"`
}

type PostTweetResponse struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}

type AliveResponse struct {
	Alive bool `json:"alive"`
}

type TwitterAuthMessage struct {
	Message string `json:"message"`
}

type TwitterProfile struct {
	IsLoggedIn bool   `json:"isLoggedIn"`
	UserName   string `json:"username"`
	Error      string `json:"error"`
}
