package handler

import "net/http"

var (
	baseMsg = BaseResponse{
		ResultCode: http.StatusOK,
		ResultMsg:  "success",
		Successful: true,
	}
)
