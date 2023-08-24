package alipay

import (
	"encoding/json"
	"errors"
)

/*
 * 付款码支付
 * 统一收单交易撤销接口
 * alipay.trade.cancel
 */

// 统一收单交易撤销接口 请求 参数
type TradeCancelRequest struct {
	OutTradeNo string `json:"out_trade_no,omitempty"` //特殊可选	64 原支付请求的商户订单号,和支付宝交易号不能同时为空 20150320010101001
	TradeNo    string `json:"trade_no,omitempty"`     //特殊可选	64 支付宝交易号，和商户订单号不能同时为空 2014112611001004680073956707
}

func (t *TradeCancelRequest) HttpMethod() string {
	return "POST"
}

func (t *TradeCancelRequest) ApiParamMethod() string {
	return "alipay.trade.cancel"
}

// 从响应数据中生成结构体
func (t *TradeCancelRequest) GenResponse(data []byte) (IResponse, error) {
	type tempResponse struct {
		TradeCancelResponse TradeCancelResponse `json:"alipay_trade_cancel_response"`
	}
	var tResp tempResponse
	err := json.Unmarshal(data, &tResp)
	if err != nil {
		return nil, err
	}

	return tResp.TradeCancelResponse, nil
}

func (t *TradeCancelRequest) DoValidate() error {
	if len(t.OutTradeNo) < 1 && len(t.TradeNo) < 1 {
		err := errors.New("原支付请求的商户订单号,和支付宝交易号不能同时为空")
		return err
	}
	return nil
}

// 统一收单交易撤销接口 响应 参数
type TradeCancelResponse struct {
	TradeNo    string `json:"trade_no"`     //必选	64 支付宝交易号; 当发生交易关闭或交易退款时返回； 2013112011001004330000121536
	OutTradeNo string `json:"out_trade_no"` //必选	64 商户订单号 6823789339978248
	RetryFlag  string `json:"retry_flag"`   //必选	1 是否需要重试 N
	Action     string `json:"action"`       //必选	10 本次撤销触发的交易动作,接口调用成功且交易存在时返回。
	//可能的返回值：
	//close：交易未支付，触发关闭交易动作，无退款；
	//refund：交易已支付，触发交易退款动作；
	//未返回：未查询到交易，或接口调用失败；
}

func (t TradeCancelResponse) ApiParamMethod() string {
	return "alipay.trade.cancel"
}
