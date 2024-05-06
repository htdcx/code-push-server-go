package constants

const (
	GIN_USER_ID = "GIN_USER_ID"
	GIN_LANG    = "LANG"
)
const (
	REDIS_TOKEN_INFO  = "TOKEN:"
	REDIS_UPDATE_INFO = "UPDATE_INFO:"
)

const (
	CONFIG_LOGIN_VERIFICATION    = "CONFIG_LOGIN_VERIFICATION"
	CONFIG_REGISTER_VERIFICATION = "CONFIG_REGISTER_VERIFICATION"
)

type ErrObj struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type PageData[T any] struct {
	Data       []T   `json:"data"`
	TotalCount int64 `json:"totalCount"`
}

type PageBean struct {
	Page int `json:"page"`
	Rows int `json:"rows"`
}

func (PageBean) GetNew() PageBean {
	return PageBean{Rows: 10}
}
