package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	crawler "ASS/crawl/vmcrawler"
)

// @Summary 获取所有A股信息
// @Id 1
// @Tags A股
// @version 1.0.0
// @Accept application/json
// @Produce application/json
// @Success 200 object handler.GetCNStockInfosResponse 成功后返回值
// @Failure 400 object handler.GetCNStockInfosResponse 参数不对
// @Router /cn/stockInfos [get]
func GetStockInfos(c *gin.Context) {
	msg := GetCNStockInfosResponse{
		BaseResponse: baseMsg,
		Data:         crawler.StockCrawler_cn.GetStockInfos(nil),
	}

	c.JSON(msg.ResultCode, msg)
}

// @Summary 获取所指定A股信息
// @Id 2
// @Tags A股
// @version 1.0.0
// @Accept application/json
// @Produce application/json
// @Param code query string true "股票代码"
// @Success 200 object handler.GetCNStockInfoResponse 成功后返回值
// @Failure 400 object handler.GetCNStockInfoResponse 参数不对
// @Router /cn/stockInfo [get]
func GetStockInfo(c *gin.Context) {
	code := c.Query("code")
	msg := GetCNStockInfoResponse{
		BaseResponse: baseMsg,
		Data:         crawler.StockCrawler_cn.GetStockInfo(code, nil),
	}
	c.JSON(msg.ResultCode, msg)
}

// @Summary 增加股票
// @Id 3
// @Tags A股
// @version 1.0.0
// @Accept application/json
// @Produce application/json
// @Param code body handler.AddCNStockRequest true "股票代码"
// @Success 200 object handler.AddCNStockResponse 成功后返回值
// @Failure 400 object handler.AddCNStockResponse 参数不对
// @Router /cn/stock/add [post]
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

// @Summary 删除股票
// @Id 4
// @Tags A股
// @version 1.0.0
// @Accept application/json
// @Produce application/json
// @Param code body handler.DelCNStockRequest true "股票代码"
// @Success 200 object handler.DelCNStockResponse 成功后返回值
// @Failure 400 object handler.DelCNStockResponse 参数不对
// @Router /cn/stock/del [post]
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

// @Summary 获取股票爬取条件
// @Id 5
// @Tags A股
// @version 1.0.0
// @Accept application/json
// @Produce application/json
// @Success 200 object handler.GetCNStockConditionResponse 成功后返回值
// @Failure 400 object handler.GetCNStockConditionResponse 参数不对
// @Router /cn/condition/get [get]
func GetStockCondtion(c *gin.Context) {
	msg := GetCNStockConditionResponse{
		BaseResponse: baseMsg,
		Data:         crawler.StockCrawler_cn.GetCodition(),
	}
	c.JSON(msg.ResultCode, msg)
}

// @Summary 设置股票爬取条件
// @Id 6
// @Tags A股
// @version 1.0.0
// @Accept application/json
// @Produce application/json
// @Param code body handler.SetCNStockConditionRequest true "股票代码"
// @Success 200 object handler.SetCNStockConditionResponse 成功后返回值
// @Failure 400 object handler.SetCNStockConditionResponse 参数不对
// @Router /cn/condition/set [post]
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

// @Summary 获取股票买卖判断结果
// @Id 7
// @Tags A股
// @version 1.0.0
// @Accept application/json
// @Produce application/json
// @Param   code query string true "股票代码"
// @Success 200 object handler.GetCNStockJudgeResultResponse 成功后返回值
// @Failure 400 object handler.GetCNStockJudgeResultResponse 参数不对
// @Router /cn/judge/get [get]
func GetJudgeResult(c *gin.Context) {
	code := c.Query("code")
	msg := GetCNStockJudgeResultResponse{
		BaseResponse: baseMsg,
		Data:         crawler.Judger_cn.GetJudgeResult(code),
	}
	c.JSON(msg.ResultCode, msg)
}

// @Summary 设置股票买卖策略
// @Id 8
// @Tags A股
// @version 1.0.0
// @Accept application/json
// @Produce application/json
// @Param code body handler.SetCNStockStrategyRequest true "股票买卖策略信息"
// @Success 200 object handler.SetCNStockConditionResponse 成功后返回值
// @Failure 400 object handler.SetCNStockConditionResponse 参数不对
// @Router /cn/strategy/set [post]
func SetStockStrategy(c *gin.Context) {
	msg := SetCNStockConditionResponse{
		BaseResponse: baseMsg,
	}
	r := SetCNStockStrategyRequest{}
	err := c.ShouldBind(&r)
	if err != nil {
		logrus.Errorln(err)
		msg.ResultCode = http.StatusBadRequest
		msg.ResultMsg = "pararms exception!"
		msg.Successful = false
	} else {
		crawler.Judger_cn.SetStrategy(r.Code, r.Strategy)
	}
	c.JSON(msg.ResultCode, msg)
}

// @Summary 获取股票爬取条件
// @Id 9
// @Tags A股
// @version 1.0.0
// @Accept application/json
// @Produce application/json
// @Success 200 object handler.GetCanBuyStocksResponse 成功后返回值
// @Failure 400 object handler.GetCanBuyStocksResponse 参数不对
// @Router /cn/canbuy/get [get]
func GetCanBuyStocks(c *gin.Context) {
	msg := GetCanBuyStocksResponse{
		BaseResponse: baseMsg,
		Data:         crawler.Judger_cn.GetCanBuyStocks(),
	}
	c.JSON(msg.ResultCode, msg)
}

// @Summary 获取股票爬取条件
// @Id 10
// @Tags A股
// @version 1.0.0
// @Accept application/json
// @Produce application/json
// @Success 200 object handler.GetCanSellStocksResponse 成功后返回值
// @Failure 400 object handler.GetCanSellStocksResponse 参数不对
// @Router /cn/cansell/get [get]
func GetCanSellStocks(c *gin.Context) {
	msg := GetCanSellStocksResponse{
		BaseResponse: baseMsg,
		Data:         crawler.Judger_cn.GetCanSellStocks(),
	}
	c.JSON(msg.ResultCode, msg)
}
