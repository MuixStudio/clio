package auth

type (
	LoginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	LoginResp struct {
		Username string `json:"username"`
		UserID   uint32 `json:"userId"`
	}
)
