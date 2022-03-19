package handler

type BaseResponse struct {
	ResultCode int    `json:"resultCode"`
	Successful bool   `json:"successful"`
	ResultMsg  string `json:"resultMsg"`
}

type GetCNStockInfosResponse struct {
	BaseResponse
	Data string `json:"data"`
}

type GetCNStockInfoResponse struct {
	BaseResponse
	Data string `json:"data"`
}

type AddCNStockRequest struct {
	Code string `json:"code" binding:"required"`
}

type AddCNStockResponse struct {
	BaseResponse
}

type DelCNStockRequest struct {
	Code string `json:"code" binding:"required"`
}

type DelCNStockResponse struct {
	BaseResponse
}

type SetCNStockConditionRequest struct {
	Condition string `json:"condition" binding:"required"`
}

type SetCNStockConditionResponse struct {
	BaseResponse
}

type GetCNStockConditionResponse struct {
	BaseResponse
	Data string `json:"data"`
}
