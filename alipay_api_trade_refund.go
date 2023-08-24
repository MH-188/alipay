package alipay

import (
	"encoding/json"
	"fmt"
)

/*
 * 付款码支付
 * 统一收单交易退款接口
 * alipay.trade.refund
 */

// 退款包含的商品列表信息
type RefundGoodsDetail struct {
	GoodsId      string  `json:"goods_id,omitempty"`      //必选	32 商品编号。 对应支付时传入的goods_id apple-01
	RefundAmount float64 `json:"refund_amount,omitempty"` //必选	9 该商品的退款总金额，单位为元 19.50
	OutItemId    string  `json:"out_item_id,omitempty"`   //可选	64 商家侧小程序商品ID，具体使用方式请参考：https://opendocs.alipay.com/pre-open/06uila?pathHash=87297d0a  outItem_01
	OutSkuId     string  `json:"out_sku_id,omitempty"`    //可选	64 商家侧小程序商品sku ID，具体使用方式请参考：https://opendocs.alipay.com/pre-open/06uila?pathHash=87297d0a outSku_01
}

// 退分账明细信息
type OpenApiRoyaltyDetailInfoPojo struct {
	RoyaltyType  string  `json:"royalty_type,omitempty"`   //可选	32  分账类型. 普通分账为：transfer; 补差为：replenish; 为空默认为分账transfer; transfer
	TransOut     string  `json:"trans_out,omitempty"`      //可选	16 支出方账户。如果支出方账户类型为userId，本参数为支出方的支付宝账号对应的支付宝唯一用户号，以2088开头的纯16位数字；如果支出方类型为loginName，本参数为支出方的支付宝登录号。 泛金融类商户分账时，该字段不要上送。  2088101126765726
	TransOutType string  `json:"trans_out_type,omitempty"` //可选	64 支出方账户类型。userId表示是支付宝账号对应的支付宝唯一用户号;loginName表示是支付宝登录号； 泛金融类商户分账时，该字段不要上送。  userId
	TransInType  string  `json:"trans_in_type,omitempty"`  //可选	64 收入方账户类型。userId表示是支付宝账号对应的支付宝唯一用户号;cardAliasNo表示是卡编号;loginName表示是支付宝登录号； userId
	TransIn      string  `json:"trans_in,omitempty"`       //必选	16 收入方账户。如果收入方账户类型为userId，本参数为收入方的支付宝账号对应的支付宝唯一用户号，以2088开头的纯16位数字；如果收入方类型为cardAliasNo，本参数为收入方在支付宝绑定的卡编号；如果收入方类型为loginName，本参数为收入方的支付宝登录号； 2088101126708402
	Amount       float64 `json:"amount,omitempty"`         //可选	9 分账的金额，单位为元 0.1
	Desc         string  `json:"desc,omitempty"`           //可选	1000 分账描述 分账给2088101126708402
	RoyaltyScene string  `json:"royalty_scene,omitempty"`  //可选	256 可选值：达人佣金、平台服务费、技术服务费、其他  达人佣金
	TransInName  string  `json:"trans_in_name,omitempty"`  //可选	64 分账收款方姓名，上送则进行姓名与支付宝账号的一致性校验，校验不一致则分账失败。不上送则不进行姓名校验 张三
}

// 组合支付退费明细
type RefundSubFee struct {
	RefundChargeFee string `json:"refund_charge_fee,omitempty"` //可选	11 实退费用。单位：元。 0.10
	SwitchFeeRate   string `json:"switch_fee_rate,omitempty"`   //可选	64 签约费率 0.01
}

// 退费信息
type RefundChargeInfo struct {
	RefundChargeFee        string         `json:"refund_charge_fee,omitempty"`          //可选	11 实退费用。单位：元。 0.01
	SwitchFeeRate          string         `json:"switch_fee_rate,omitempty"`            //可选	64 签约费率 0.01
	ChargeType             string         `json:"charge_type,omitempty"`                //可选	64 收单手续费trade，花呗分期手续hbfq，其他手续费charge  trade
	RefundSubFeeDetailList []RefundSubFee `json:"refund_sub_fee_detail_list,omitempty"` //可选 组合支付退费明细
}

// 统一收单交易退款接口 请求 参数
type TradeRefundRequest struct {
	OutTradeNo              string                         `json:"out_trade_no,omitempty"`              //特殊可选	64 商户订单号。 订单支付时传入的商户订单号，商家自定义且保证商家系统中唯一。与支付宝交易号 trade_no 不能同时为空 20150320010101001
	TradeNo                 string                         `json:"trade_no,omitempty"`                  //特殊可选	64 支付宝交易号。 和商户订单号 out_trade_no 不能同时为空。 2014112611001004680073956707
	RefundAmount            float64                        `json:"refund_amount"`                       //必选	11 退款金额。 需要退款的金额，该金额不能大于订单金额，单位为元，支持两位小数。注：如果正向交易使用了营销，该退款金额包含营销金额，支付宝会按业务规则分配营销和买家自有资金分别退多少，默认优先退买家的自有资金。如交易总金额100元，用户支付时使用了80元自有资金和20元无资金流的营销券，商家实际收款80元。如果首次请求退款60元，则60元全部从商家收款资金扣除退回给用户自有资产；如果再请求退款40元，则从商家收款资金扣除20元退回用户资产以及把20元的营销券退回给用户（券是否可再使用取决于券的规则配置）。  200.12
	RefundReason            string                         `json:"refund_reason,omitempty"`             //可选	256 退款原因说明。 商家自定义，将在会在商户和用户的pc退款账单详情中展示 正常退款
	OutRequestNo            string                         `json:"out_request_no,omitempty"`            //可选	64 退款请求号。 标识一次退款请求，需要保证在交易号下唯一，如需部分退款，则此参数必传。 注：针对同一次退款请求，如果调用接口失败或异常了，重试时需要保证退款请求号不能变更，防止该笔交易重复退款。支付宝会保证同样的退款请求号多次请求只会退一次。 HZ01RF001
	RefundGoodsDetail       []RefundGoodsDetail            `json:"refund_goods_detail,omitempty"`       //可选 退款包含的商品列表信息
	RefundRoyaltyParameters []OpenApiRoyaltyDetailInfoPojo `json:"refund_royalty_parameters,omitempty"` //可选 退分账明细信息。
	//注： 1.当面付且非直付通模式无需传入退分账明细，系统自动按退款金额与订单金额的比率，从收款方和分账收入方退款，不支持指定退款金额与退款方。
	//2.直付通模式，电脑网站支付，手机 APP 支付，手机网站支付产品，须在退款请求中明确是否退分账，从哪个分账收入方退，退多少分账金额；如不明确，默认从收款方退款，收款方余额不足退款失败。不支持系统按比率退款。
	Query_options []string //可选	1024 查询选项。 商户通过上送该参数来定制同步需要额外返回的信息字段，数组格式。支持：refund_detail_item_list：退款使用的资金渠道；deposit_back_info：触发银行卡冲退信息通知；
	//枚举值
	//本次退款使用的资金渠道: refund_detail_item_list
	//银行卡冲退信息: deposit_back_info
}

func (t *TradeRefundRequest) HttpMethod() string {
	return "POST"
}

func (t *TradeRefundRequest) ApiParamMethod() string {
	return "alipay.trade.refund"
}

// 从响应数据中生成结构体
func (t *TradeRefundRequest) GenResponse(data []byte) (IResponse, error) {
	type tempResponse struct {
		AlipayTradeRefundResponse TradeRefundResponse `json:"alipay_trade_refund_response"`
	}
	var tResp tempResponse
	err := json.Unmarshal(data, &tResp)
	if err != nil {
		return nil, err
	}

	return tResp.AlipayTradeRefundResponse, nil
}

func (t *TradeRefundRequest) DoValidate() error {
	if len(t.OutTradeNo) == 0 && len(t.TradeNo) == 0 {
		return fmt.Errorf("商户订单号out_trade_no和支付宝交易号trade_no不能同时为空")
	}
	if t.RefundAmount < 0 {
		return fmt.Errorf("请求退款金额不能为负值")
	}
	return nil
}

// 统一收单交易退款接口 响应 参数
type TradeRefundResponse struct {
	ErrorParam
	TradeNo              string             `json:"trade_no"`                          //必选	64 2013112011001004330000121536 支付宝交易号
	OutTradeNo           string             `json:"out_trade_no"`                      //必选	64 商户订单号 6823789339978248
	BuyerLogonId         string             `json:"buyer_logon_id"`                    //必选	100 用户的登录id 159****5620
	FundChange           string             `json:"fund_change"`                       //必选	1 本次退款是否发生了资金变化 Y
	RefundFee            string             `json:"refund_fee"`                        //必选	11 退款总金额。单位：元。 指该笔交易累计已经退款成功的金额。 88.88
	RefundDetailItemList []TradeFundBill    `json:"refund_detail_item_list,omitempty"` //特殊可选 退款使用的资金渠道。 只有在签约中指定需要返回资金明细，或者入参的query_options中指定时才返回该字段信息。
	StoreName            string             `json:"store_name,omitempty"`              //特殊可选	512 交易在支付时候的门店名称  望湘园联洋店
	BuyerUserId          string             `json:"buyer_user_id,omitempty"`           //特殊可选	28 买家在支付宝的用户id 2088101117955611
	SendBackFee          string             `json:"send_back_fee,omitempty"`           //特殊可选	11 本次商户实际退回金额。单位：元。 说明：如需获取该值，需在入参query_options中传入 refund_detail_item_list。 1.8
	RefundHybAmount      string             `json:"refund_hyb_amount,omitempty"`       //可选	11 本次请求退惠营宝金额。单位：元。 10.24
	RefundChargeInfoList []RefundChargeInfo `json:"refund_charge_info_list,omitempty"` //可选 退费信息
	GmtRefundPay         string             `json:"gmt_refund_pay,omitempty"`          //文档中没有说明，实际返回值中存在 时间字符串 2023-08-18 14:08:20
}

func (t TradeRefundResponse) ApiParamMethod() string {
	return "alipay.trade.refund"
}
