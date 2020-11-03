package macro

const (
	STATUS_AUTH_FAILED = 1000000
)

var ERR_MSG = map[int64]string{
	STATUS_AUTH_FAILED: "登陆校验失败",
}

type Error struct {
	Code int64
	Msg string
}