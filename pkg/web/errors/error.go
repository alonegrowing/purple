package errors

import "purple/pkg/macro"

const httpMemberNotExist = 10000


var ErrMsg  = map[int64]string {
	httpMemberNotExist: "用户不存在",
}

type HttpErrorInterface interface{
	MemberNotExistError() macro.Error
}

type HtppErrorImpl struct {}

var HtppError HttpErrorInterface


func init() {
	HtppError = NewHtppErrorImpl()
}
func NewHtppErrorImpl() *HtppErrorImpl {
	return &HtppErrorImpl{}
}

func (r *HtppErrorImpl) MemberNotExistError() macro.Error {
	return macro.Error{
		Code: httpMemberNotExist,
		Msg: ErrMsg[httpMemberNotExist],
	}
}