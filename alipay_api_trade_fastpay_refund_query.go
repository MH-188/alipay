package alipay

import (
	"encoding/json"
	"fmt"
)

/*
 * 付款码支付
 * 统一收单交易退款查询
 * alipay.trade.fastpay.refund.query
 */

// 退分账明细信息，当前仅在直付通产品中返回。
type RefundRoyaltyResult struct {
	RefundAmount  string `json:"refund_amount,omitempty"`   //必选	9 退分账金额。单位：元。 10
	RoyaltyType   string `json:"royalty_type,omitempty"`    //可选	32 分账类型. 字段为空默认为普通分账类型transfer 枚举值 普通分账类型: transfer 补差分账类型: replenish    transfer
	ResultCode    string `json:"result_code,omitempty"`     //必选	32 退分账结果码 SUCCESS
	TransOut      string `json:"trans_out,omitempty"`       //可选	28 转出人支付宝账号对应用户ID 2088102210397302
	TransOutEmail string `json:"trans_out_email,omitempty"` //可选	64 转出人支付宝账号 alipay-test03@alipay.com
	TransIn       string `json:"trans_in,omitempty"`        //可选	28 转入人支付宝账号对应用户ID 2088102210397302
	TransInEmail  string `json:"trans_in_email,omitempty"`  //可选	64 转入人支付宝账号 zen_gwen@hotmail.com
}

type DepositBackInfo struct {
	HasDepositBack     string `json:"has_deposit_back,omitempty"`      //可选	10 是否存在银行卡冲退信息。 true
	DbackStatus        string `json:"dback_status,omitempty"`          //可选	8 银行卡冲退状态。S-成功，F-失败，P-处理中。银行卡冲退失败，资金自动转入用户支付宝余额。 S
	DbackAmount        string `json:"dback_amount,omitempty"`          //可选	9 银行卡冲退金额。单位：元。 1.01
	BankAckTime        string `json:"bank_ack_time,omitempty"`         //可选	32 银行响应时间，格式为yyyy-MM-dd HH:mm:ss 2020-06-02 14:03:48
	EstBankReceiptTime string `json:"est_bank_receipt_time,omitempty"` //可选	32 预估银行到账时间，格式为yyyy-MM-dd HH:mm:ss 2020-06-02 14:03:48
}

// 统一收单交易支付接口 请求 参数
type TradeFastPayRefundQueryRequest struct {
	TradeNo      string   `json:"trade_no,omitempty"`       //特殊可选	64 支付宝交易号。 和商户订单号不能同时为空 2021081722001419121412730660
	OutTradeNo   string   `json:"out_trade_no,omitempty"`   //特殊可选	64 商户订单号。 订单支付时传入的商户订单号,和支付宝交易号不能同时为空。 trade_no,out_trade_no如果同时存在优先取trade_no 2014112611001004680073956707
	OutRequestNo string   `json:"out_request_no,omitempty"` //必选	64 退款请求号。  请求退款接口时，传入的退款请求号，如果在退款请求时未传入，则该值为创建交易时的商户订单号。  HZ01RF001
	QueryOptions []string `json:"query_options,omitempty"`  //可选	1024 查询选项，商户通过上送该参数来定制同步需要额外返回的信息字段，数组格式。
	//枚举值
	//本次退款使用的资金渠道: refund_detail_item_list
	//退款执行成功的时间: gmt_refund_pay
	//银行卡冲退信息: deposit_back_info
	//refund_detail_item_list
}

func (t *TradeFastPayRefundQueryRequest) HttpMethod() string {
	return "POST"
}

func (t *TradeFastPayRefundQueryRequest) ApiParamMethod() string {
	return "alipay.trade.fastpay.refund.query"
}

// 从响应数据中生成结构体
func (t *TradeFastPayRefundQueryRequest) GenResponse(data []byte) (IResponse, error) {
	type tempResponse struct {
		TradeFastPayRefundQueryResponse TradeFastPayRefundQueryResponse `json:"alipay_trade_fastpay_refund_query_response"`
	}
	var tResp tempResponse
	err := json.Unmarshal(data, &tResp)
	if err != nil {
		return nil, err
	}

	return tResp.TradeFastPayRefundQueryResponse, nil
}

func (t *TradeFastPayRefundQueryRequest) DoValidate() error {
	if len(t.OutTradeNo) == 0 && len(t.TradeNo) == 0 {
		return fmt.Errorf("商户订单号out_trade_no和支付宝交易号trade_no不能同时为空")
	}
	if len(t.OutRequestNo) < 1 {
		return fmt.Errorf("退款请求号不能为空")
	}
	return nil
}

// 统一收单交易退款查询接口 响应 参数
type TradeFastPayRefundQueryResponse struct {
	ErrorParam
	TradeNo              string                `json:"trade_no,omitempty"`                //特殊可选	64 支付宝交易号 2014112611001004680073956707
	OutTradeNo           string                `json:"out_trade_no,omitempty"`            //特殊可选	64 创建交易传入的商户订单号 20150320010101001
	OutRequestNo         string                `json:"out_request_no,omitempty"`          //特殊可选	64 本笔退款对应的退款请求号 20150320010101001
	TotalAmount          string                `json:"total_amount,omitempty"`            //特殊可选	11 该笔退款所对应的交易的订单金额。单位：元。 100.20
	RefundAmount         string                `json:"refund_amount,omitempty"`           //特殊可选	11 本次退款请求，对应的退款金额。单位：元。 12.33
	RefundStatus         string                `json:"refund_status,omitempty"`           //特殊可选	32 退款状态。 枚举值： REFUND_SUCCESS 退款处理成功； 未返回该字段表示退款请求未收到或者退款失败； 注：如果退款查询发起时间早于退款时间，或者间隔退款发起时间太短，可能出现退款查询时还没处理成功，后面又处理成功的情况，建议商户在退款发起后间隔10秒以上再发起退款查询请求。
	RefundRoyaltys       []RefundRoyaltyResult `json:"refund_royaltys,omitempty"`         //特殊可选 退分账明细信息，当前仅在直付通产品中返回。
	GmtRefundPay         string                `json:"gmt_refund_pay,omitempty"`          //特殊可选	32 退款时间。默认不返回该信息，需要在入参的query_options中指定"gmt_refund_pay"值时才返回该字段信息。 2014-11-27 15:45:57
	RefundDetailItemList []TradeFundBill       `json:"refund_detail_item_list,omitempty"` //特殊可选 本次退款使用的资金渠道； 默认不返回该信息，需要在入参的query_options中指定"refund_detail_item_list"值时才返回该字段信息。
	SendBackFee          string                `json:"send_back_fee,omitempty"`           //特殊可选	11 本次商户实际退回金额；单位：元。 默认不返回该信息，需要在入参的query_options中指定"refund_detail_item_list"值时才返回该字段信息。 88
	DepositBackInfo      DepositBackInfo       `json:"deposit_back_info,omitempty"`       //特殊可选 银行卡冲退信息； 默认不返回该信息，需要在入参的query_options中指定"deposit_back_info"值时才返回该字段信息。
	RefundHybAmount      string                `json:"refund_hyb_amount,omitempty"`       //可选	11 本次退款金额中退惠营宝的金额。单位：元。 10.24
	RefundChargeInfoList []RefundChargeInfo    `json:"refund_charge_info_list,omitempty"` //可选 退费信息
}

func (t TradeFastPayRefundQueryResponse) ApiParamMethod() string {
	return "alipay.trade.fastpay.refund.query"
}
