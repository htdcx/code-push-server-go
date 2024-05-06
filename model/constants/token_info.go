package constants

type TokenInfo struct {
	Uid      *int    `json:"uid"`
	UserName *string `json:"userName"`
	Token    *string `json:"token"`
	Money    *int64  `json:"money"`
	VipTime  *int64  `json:"vipTime"`
	// PlatformId int
}
