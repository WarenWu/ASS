// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "wuweiming",
            "email": "wuweimingoen@163.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/cn/condition/get": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "A股"
                ],
                "summary": "获取股票爬取条件",
                "operationId": "5",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.GetCNStockConditionResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler.GetCNStockConditionResponse"
                        }
                    }
                }
            }
        },
        "/cn/condition/set": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "A股"
                ],
                "summary": "设置股票爬取条件",
                "operationId": "6",
                "parameters": [
                    {
                        "description": "股票代码",
                        "name": "code",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.SetCNStockConditionRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.SetCNStockConditionResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler.SetCNStockConditionResponse"
                        }
                    }
                }
            }
        },
        "/cn/judge/get": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "A股"
                ],
                "summary": "获取股票买卖判断结果",
                "operationId": "7",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.GetCNStockJudgeResultResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler.GetCNStockJudgeResultResponse"
                        }
                    }
                }
            }
        },
        "/cn/stock/add": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "A股"
                ],
                "summary": "增加股票",
                "operationId": "3",
                "parameters": [
                    {
                        "description": "股票代码",
                        "name": "code",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.AddCNStockRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.AddCNStockResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler.AddCNStockResponse"
                        }
                    }
                }
            }
        },
        "/cn/stock/del": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "A股"
                ],
                "summary": "删除股票",
                "operationId": "4",
                "parameters": [
                    {
                        "description": "股票代码",
                        "name": "code",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.DelCNStockRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.DelCNStockResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler.DelCNStockResponse"
                        }
                    }
                }
            }
        },
        "/cn/stockInfo": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "A股"
                ],
                "summary": "获取所指定A股信息",
                "operationId": "2",
                "parameters": [
                    {
                        "type": "string",
                        "description": "股票代码",
                        "name": "code",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.GetCNStockInfoResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler.GetCNStockInfoResponse"
                        }
                    }
                }
            }
        },
        "/cn/stockInfos": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "A股"
                ],
                "summary": "获取所有A股信息",
                "operationId": "1",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.GetCNStockInfosResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler.GetCNStockInfosResponse"
                        }
                    }
                }
            }
        },
        "/cn/strategy/set": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "A股"
                ],
                "summary": "设置股票买卖策略",
                "operationId": "8",
                "parameters": [
                    {
                        "description": "股票买卖策略信息",
                        "name": "code",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler.SetCNStockStrategyRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.SetCNStockConditionResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler.SetCNStockConditionResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handler.AddCNStockRequest": {
            "type": "object",
            "required": [
                "code"
            ],
            "properties": {
                "code": {
                    "type": "string"
                }
            }
        },
        "handler.AddCNStockResponse": {
            "type": "object",
            "properties": {
                "resultCode": {
                    "type": "integer"
                },
                "resultMsg": {
                    "type": "string"
                },
                "successful": {
                    "type": "boolean"
                }
            }
        },
        "handler.DelCNStockRequest": {
            "type": "object",
            "required": [
                "code"
            ],
            "properties": {
                "code": {
                    "type": "string"
                }
            }
        },
        "handler.DelCNStockResponse": {
            "type": "object",
            "properties": {
                "resultCode": {
                    "type": "integer"
                },
                "resultMsg": {
                    "type": "string"
                },
                "successful": {
                    "type": "boolean"
                }
            }
        },
        "handler.GetCNStockConditionResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "string"
                },
                "resultCode": {
                    "type": "integer"
                },
                "resultMsg": {
                    "type": "string"
                },
                "successful": {
                    "type": "boolean"
                }
            }
        },
        "handler.GetCNStockInfoResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.StockInfo"
                    }
                },
                "resultCode": {
                    "type": "integer"
                },
                "resultMsg": {
                    "type": "string"
                },
                "successful": {
                    "type": "boolean"
                }
            }
        },
        "handler.GetCNStockInfosResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "array",
                        "items": {
                            "$ref": "#/definitions/model.StockInfo"
                        }
                    }
                },
                "resultCode": {
                    "type": "integer"
                },
                "resultMsg": {
                    "type": "string"
                },
                "successful": {
                    "type": "boolean"
                }
            }
        },
        "handler.GetCNStockJudgeResultResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/model.StockJudgeResult"
                },
                "resultCode": {
                    "type": "integer"
                },
                "resultMsg": {
                    "type": "string"
                },
                "successful": {
                    "type": "boolean"
                }
            }
        },
        "handler.SetCNStockConditionRequest": {
            "type": "object",
            "required": [
                "condition"
            ],
            "properties": {
                "condition": {
                    "type": "string"
                }
            }
        },
        "handler.SetCNStockConditionResponse": {
            "type": "object",
            "properties": {
                "resultCode": {
                    "type": "integer"
                },
                "resultMsg": {
                    "type": "string"
                },
                "successful": {
                    "type": "boolean"
                }
            }
        },
        "handler.SetCNStockStrategyRequest": {
            "type": "object",
            "required": [
                "code",
                "strategy"
            ],
            "properties": {
                "code": {
                    "type": "string"
                },
                "strategy": {
                    "$ref": "#/definitions/model.StrategyCN"
                }
            }
        },
        "model.StockInfo": {
            "type": "object",
            "properties": {
                "asset_liability_ratio": {
                    "type": "string"
                },
                "cash_ratio": {
                    "type": "string"
                },
                "code": {
                    "type": "string"
                },
                "dividend_ratio": {
                    "type": "string"
                },
                "gross_profit_ratio": {
                    "type": "string"
                },
                "interest_ratio": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "pe": {
                    "type": "string"
                },
                "period": {
                    "type": "string"
                },
                "price": {
                    "type": "string"
                },
                "roe": {
                    "type": "string"
                }
            }
        },
        "model.StockJudgeResult": {
            "type": "object",
            "properties": {
                "can_buy": {
                    "type": "boolean"
                },
                "can_sell": {
                    "type": "boolean"
                }
            }
        },
        "model.StrategyCN": {
            "type": "object",
            "properties": {
                "aim_max_pe": {
                    "description": "卖出目标市场市盈率",
                    "type": "number"
                },
                "aim_max_stock_pe": {
                    "description": "股票卖出市盈率",
                    "type": "number"
                },
                "aim_min_pe": {
                    "description": "买入目标市场市盈率",
                    "type": "number"
                },
                "aim_min_stock_pe": {
                    "description": "股票买入市盈率",
                    "type": "number"
                },
                "pe_type": {
                    "type": "integer"
                },
                "use_yield": {
                    "type": "boolean"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0.0",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "股票信息获取",
	Description: "提供股票信息相关接口",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register("swagger", &s{})
}
