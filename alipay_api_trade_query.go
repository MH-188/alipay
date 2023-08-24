package alipay

import (
	"encoding/json"
	"fmt"
)

/*
 * 付款码支付
 * 统一收单线下交易查询接口
 * alipay.trade.query
 */

// 统一收单线下交易查询接口请求
type TradeQueryRequest struct {
	OutTradeNo   string   `json:"out_trade_no,omitempty"`  //特殊可选	64 订单支付时传入的商户订单号,和支付宝交易号不能同时为空。 trade_no,out_trade_no如果同时存在优先取trade_no 20150320010101001
	TradeNo      string   `json:"trade_no,omitempty"`      //特殊可选	64 支付宝交易号，和商户订单号不能同时为空  2014112611001004680 073956707
	QueryOptions []string `json:"query_options,omitempty"` //可选	1024  查询选项，商户通过上送该参数来定制同步需要额外返回的信息字段，数组格式。
	// 支持枚举如下：
	//交易结算信息: trade_settle_info
	//交易支付使用的资金渠道: fund_bill_list
	//交易支付时使用的所有优惠券信息: voucher_detail_list
	//交易支付使用单品券优惠的商品优惠信息: discount_goods_detail
	//商家优惠金额: mdiscount_amount
	//医保信息: medical_insurance_info
}

func (t *TradeQueryRequest) HttpMethod() string {
	return "POST"
}

func (t *TradeQueryRequest) ApiParamMethod() string {
	return "alipay.trade.query"
}

// 从响应数据中生成结构体
func (t *TradeQueryRequest) GenResponse(data []byte) (IResponse, error) {
	type tempResponse struct {
		AlipayTradeQueryResponse TradeQueryResponse `json:"alipay_trade_query_response"`
	}
	var tResp tempResponse
	err := json.Unmarshal(data, &tResp)
	if err != nil {
		return nil, err
	}

	return tResp.AlipayTradeQueryResponse, nil
}

func (t *TradeQueryRequest) DoValidate() error {
	if len(t.OutTradeNo) == 0 && len(t.TradeNo) == 0 {
		return fmt.Errorf("商户订单号out_trade_no和支付宝交易号trade_no不能同时为空")
	}
	return nil
}

// 统一收单线下交易查询接口响应
type TradeQueryResponse struct {
	ErrorParam
	TradeNo         string          `json:"trade_no"`                   //必选	64 支付宝交易号 2013112011001004330000121536
	OutTradeNo      string          `json:"out_trade_no"`               //必选	64  商家订单号 6823789339978248
	BuyerLogonId    string          `json:"buyer_logon_id"`             //必选	100 买家支付宝账号 159****5620
	TradeStatus     string          `json:"trade_status"`               //必选	32 交易状态：WAIT_BUYER_PAY（交易创建，等待买家付款）、TRADE_CLOSED（未付款交易超时关闭，或支付完成后全额退款）、TRADE_SUCCESS（交易支付成功）、TRADE_FINISHED（交易结束，不可退款） TRADE_CLOSED
	TotalAmount     string          `json:"total_amount"`               //必选	11 交易的订单金额，单位为元，两位小数。该参数的值为支付时传入的total_amount 88.88
	BuyerPayAmount  string          `json:"buyer_pay_amount,omitempty"` //可选	11 买家实付金额，单位为元，两位小数。该金额代表该笔交易买家实际支付的金额，不包含商户折扣等金额 8.88
	PointAmount     string          `json:"point_amount,omitempty"`     //可选	11 积分支付的金额，单位为元，两位小数。该金额代表该笔交易中用户使用积分支付的金额，比如集分宝或者支付宝实时优惠等 10
	InvoiceAmount   string          `json:"invoice_amount,omitempty"`   //可选	11 交易中用户支付的可开具发票的金额，单位为元，两位小数。该金额代表该笔交易中可以给用户开具发票的金额 12.11
	SendPayDate     string          `json:"send_pay_date,omitempty"`    //特殊可选	32 本次交易打款给卖家的时间 2014-11-27 15:45:57
	ReceiptAmount   string          `json:"receipt_amount,omitempty"`   //特殊可选	11 实收金额，单位为元，两位小数。该金额为本笔交易，商户账户能够实际收到的金额 15.25
	StoreId         string          `json:"store_id,omitempty"`         //特殊可选	32 商户门店编号 NJ_S_001
	TerminalId      string          `json:"terminal_id,omitempty"`      //特殊可选	32 商户机具终端编号 NJ_T_001
	FundBillList    []TradeFundBill `json:"fund_bill_list"`             //必选 交易支付使用的资金渠道。 只有在签约中指定需要返回资金明细，或者入参的query_options中指定时才返回该字段信息。
	StoreName       string          `json:"store_name,omitempty"`       //特殊可选	512 请求交易支付中的商户店铺的名称 证大五道口店
	BuyerUserId     string          `json:"buyer_user_id"`              //必选	16 买家在支付宝的用户id 2088101117955611
	BuyerUserType   string          `json:"buyer_user_type,omitempty"`  //特殊可选	18 买家用户类型。CORPORATE:企业用户；PRIVATE:个人用户。
	MdiscountAmount string          `json:"mdiscount_amount,omitempty"` //特殊可选	11 商家优惠金额 88.88
	DiscountAmount  string          `json:"discount_amount,omitempty"`  //特殊可选	11 平台优惠金额 88.88
	ExtInfos        string          `json:"ext_infos,omitempty"`        //特殊可选	1024 交易额外信息，特殊场景下与支付宝约定返回。 json格式。  {"action":"cancel"}
}

func (t TradeQueryResponse) ApiParamMethod() string {
	return "alipay.trade.query"
}
