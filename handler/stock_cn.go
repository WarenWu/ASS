package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	crawler "ASS/crawl/vmcrawler"
)

// @Summary 获取所有A股票池信息
// @Id 1
// @Tags A股
// @version 1.0.0
// @Accept application/json
// @Produce application/json
// @Param streamName query string true "流id"
// @Success 200 object handler.GetStreamInfoResponse 成功后返回值
// @Failure 400 object handler.GetStreamInfoResponse 参数不对
// @Router /streamInfo [get]
func GetStockInfos(c *gin.Context) {
	msg := GetCNStockInfosResponse{
		BaseResponse: baseMsg,
		Data:         crawler.StockCrawler_cn.GetStockInfos(nil),
	}

	c.JSON(msg.ResultCode, msg)
}

func GetStockInfo(c *gin.Context) {
	code := c.Query("stockcode")
	msg := GetCNStockInfoResponse{
		BaseResponse: baseMsg,
		Data:         crawler.StockCrawler_cn.GetStockInfo(code, nil),
	}
	c.JSON(msg.ResultCode, msg)
}

func AddStock(c *gin.Context) {
	msg := AddCNStockResponse{
		BaseResponse: baseMsg,
	}
	r := AddCNStockRequest{}
	err := c.ShouldBind(&r)
	if err != nil {
		logrus.Errorln(err)
		msg.ResultCode = http.StatusBadRequest
		msg.ResultMsg = "pararms exception!"
		msg.Successful = false
	} else {
		crawler.StockCrawler_cn.PutStockCode(r.Code)
	}
	c.JSON(msg.ResultCode, msg)
}

func DelStock(c *gin.Context) {
	msg := DelCNStockResponse{
		BaseResponse: baseMsg,
	}
	r := DelCNStockRequest{}
	err := c.ShouldBind(&r)
	if err != nil {
		logrus.Errorln(err)
		msg.ResultCode = http.StatusBadRequest
		msg.ResultMsg = "pararms exception!"
		msg.Successful = false
	} else {
		crawler.StockCrawler_cn.DelStockCode(r.Code)
	}
	c.JSON(msg.ResultCode, msg)
}

func GetStockCondtion(c *gin.Context) {
	msg := GetCNStockConditionResponse{
		BaseResponse: baseMsg,
		Data:         crawler.StockCrawler_cn.GetCodition(),
	}
	c.JSON(msg.ResultCode, msg)
}

func SetStockCondtion(c *gin.Context) {
	msg := SetCNStockConditionResponse{
		BaseResponse: baseMsg,
	}
	r := SetCNStockConditionRequest{}
	err := c.ShouldBind(&r)
	if err != nil {
		logrus.Errorln(err)
		msg.ResultCode = http.StatusBadRequest
		msg.ResultMsg = "pararms exception!"
		msg.Successful = false
	} else {
		crawler.StockCrawler_cn.SetCodition(r.Condition)
	}
	c.JSON(msg.ResultCode, msg)
}
