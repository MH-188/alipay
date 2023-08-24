package alipay

import (
	"encoding/json"
	"fmt"
)

/*
 * 付款码支付
 * 统一收单交易支付接口
 * alipay.trade.pay
 * 只有支付时设置了回调接口，退款才会收到回调通知
 */

// 订单包含的商品列表信息，json格式。
type GoodsDetails struct {
	GoodsId        string  `json:"goods_id,omitempty"`        //必选	64 商品的编号 apple-01
	GoodsName      string  `json:"goods_name,omitempty"`      //必选	256 商品名称 ipad
	Quantity       float64 `json:"quantity,omitempty"`        //必选	32 商品数量 1
	Price          float64 `json:"price,omitempty"`           //必选	9 商品单价，单位为元 2000
	GoodsCategory  string  `json:"goods_category,omitempty"`  //可选	24 商品类目 34543238
	CategoriesTree string  `json:"categories_tree,omitempty"` //可选	128 商品类目树，从商品类目根节点到叶子节点的类目id组成，类目id值使用|分割 124868003|126232002|126252004
	ShowUrl        string  `json:"show_url,omitempty"`        //可选	400 商品的展示地址 http://www.alipay.com/xxx.jpg
}

// 业务扩展参数
type ExtendParams struct {
	SysServiceProviderId string `json:"sys_service_provider_id,omitempty"` //可选	64 系统商编号 该参数作为系统商返佣数据提取的依据，请填写系统商签约协议的PID 2088511833207846
	SpecifiedSellerName  string `json:"specified_seller_name,omitempty"`   //可选	32 特殊场景下，允许商户指定交易展示的卖家名称 XXX的跨境小铺
	CardType             string `json:"card_type,omitempty"`               //可选	32 卡类型 S0JP0000
}

// 商户传入业务信息，具体值要和支付宝约定，应用于安全，营销等参数直传场景
type BusinessParams struct {
	McCreateTradeIp string `json:"mc_create_trade_ip,omitempty"` //可选	128 商户端创建订单的 IP，须上传正确的用户端外网 IP，支持 ipv4/ipv6 格式；mc_create_trade_ip和mcCreateTradeIp（旧）参数描述相同，首选mc_create_trade_ip入参，请勿重复入参； 如已入参mcCreateTradeIp（旧），无需新增入参mc_create_trade_ip。  127.0.0.1
}

// 优惠明细参数，通过此属性补充营销参数。 注：仅与支付宝协商后可用。
type PromoParam struct {
	ActualOrderTime string   `json:"actual_order_time,omitempty"` //可选	32 存在延迟扣款这一类的场景，用这个时间表明用户发生交易的时间，比如说，在公交地铁场景，用户刷码出站的时间，和商户上送交易的时间是不一样的。  2018-09-25 22:47:33
	StoreId         string   `json:"store_id,omitempty"`          //可选	32 商户门店编号。指商户创建门店时输入的门店编号。  NJ_001
	OperatorId      string   `json:"operator_id,omitempty"`       //可选	28 商户操作员编号。 yx_001
	TerminalId      string   `json:"terminal_id,omitempty"`       //可选	32 商户机具终端编号。NJ_T_001
	QueryOptions    []string `json:"query_options,omitempty"`     //可选	1024 返回参数选项。 商户通过传递该参数来定制同步需要额外返回的信息字段，数组格式。如：["fund_bill_list","voucher_detail_list","discount_goods_detail"]  示例：["fund_bill_list","voucher_detail_list","discount_goods_detail"]
	//query_options 枚举值
	//资金明细信息: fund_bill_list
	//优惠券信息: voucher_detail_list
	//因公付金额信息: enterprise_pay_info
	//惠营宝回票金额信息: hyb_amount
	//商品优惠信息: discount_goods_detail
	//平台优惠金额: discount_amount
	//商家优惠金额: mdiscount_amount
}

// 交易支付使用的资金渠道。 只有在签约中指定需要返回资金明细，或者入参的query_options中指定时才返回该字段信息。
type TradeFundBill struct {
	FundChannel string `json:"fund_channel"`        //必选	32 交易使用的资金渠道，详见 支付渠道列表 ALIPAYACCOUNT
	Amount      string `json:"amount"`              //必选	32 该支付工具类型所使用的金额 10
	RealAmount  string `json:"real_amount"`         //可选	11 渠道实际付款金额 11.21
	FundType    string `json:"fund_type,omitempty"` //可选	32 (退款交易接口响应中使用该字段)渠道所使用的资金类型,目前只在资金渠道(fund_channel)是银行卡渠道(BANKCARD)的情况下才返回该信息(DEBIT_CARD:借记卡,CREDIT_CARD:信用卡,MIXED_CARD:借贷合一卡) DEBIT_CARD
}

// 本交易支付时使用的所有优惠券信息。 只有在query_options中指定时才返回该字段信息。
type VoucherDetail struct {
	Id   string `json:"id"`   //必选	32 券id 2015102600073002039000002D5O
	Name string `json:"name"` //必选	64 券名称 XX超市5折优惠
	Type string `json:"type"` //必选	32 券类型
	//注：不排除将来新增其他类型的可能，商家接入时注意兼容性避免硬编码
	//枚举值
	//全场代金券: ALIPAY_FIX_VOUCHER
	//折扣券: ALIPAY_DISCOUNT_VOUCHER
	//单品优惠券: ALIPAY_ITEM_VOUCHER
	//现金抵价券: ALIPAY_CASH_VOUCHER
	//商家全场券: ALIPAY_BIZ_VOUCHER
	Amount                     string `json:"amount"`                       //必选	8 优惠券面额，它应该会等于商家出资加上其他出资方出资 10.00
	MerchantContribute         string `json:"merchant_contribute"`          //可选	8 商家出资（特指发起交易的商家出资金额） 9.00
	OtherContribute            string `json:"other_contribute"`             //可选	8 其他出资方出资金额，可能是支付宝，可能是品牌商，或者其他方，也可能是他们的一起出资 1.00
	Memo                       string `json:"memo"`                         //可选	256 优惠券备注信息 学生专用优惠
	TemplateId                 string `json:"template_id"`                  //可选	64 券模板id 20171030000730015359000EMZP0
	PurchaseBuyerContribute    string `json:"purchase_buyer_contribute"`    //可选	8 如果使用的这张券是用户购买的，则该字段代表用户在购买这张券时用户实际付款的金额 2.01
	PurchaseMerchantContribute string `json:"purchase_merchant_contribute"` //可选	8 如果使用的这张券是用户购买的，则该字段代表用户在购买这张券时商户优惠的金额 1.03
	PurchaseAntContribute      string `json:"purchase_ant_contribute"`      //可选	8 如果使用的这张券是用户购买的，则该字段代表用户在购买这张券时平台优惠的金额 0.82
}

// 统一收单交易支付接口 请求 参数
type TradePayRequest struct {
	OutTradeNo     string         `json:"out_trade_no"`              //必选	64 商户订单号。由商家自定义，64个字符以内，仅支持字母、数字、下划线且需保证在商户端不重复。   20150320010101001
	TotalAmount    float64        `json:"total_amount"`              //必选	11 订单总金额。单位为元，精确到小数点后两位，取值范围：[0.01,100000000]。 88.88
	Subject        string         `json:"subject"`                   //必选	256 订单标题。注意：不可使用特殊字符，如 /，=，& 等。 Iphone6 16G
	AuthCode       string         `json:"auth_code"`                 //必选	64 支付授权码。 当面付场景传买家的付款码（25~30开头的长度为16~24位的数字，实际字符串长度以开发者获取的付款码长度为准）或者刷脸标识串（fp开头的35位字符串）。 28763443825664394
	Scene          string         `json:"scene"`                     //必选	32 支付场景。枚举值： bar_code：当面付条码支付场景； security_code：当面付刷脸支付场景，对应的auth_code为fp开头的刷脸标识串； 默认值为bar_code。 枚举值 当面付条码支付场景: bar_code 当面付刷脸支付场景，对应的auth_code为fp开头的刷脸标识串: security_code bar_code
	ProductCode    string         `json:"product_code,omitempty"`    //可选	64  产品码。 商家和支付宝签约的产品码。 当面付场景下，如果签约的是当面付快捷版，则传 OFFLINE_PAYMENT; 其它支付宝当面付产品传 FACE_TO_FACE_PAYMENT； 不传则默认使用FACE_TO_FACE_PAYMENT。  FACE_TO_FACE_PAYMENT
	SellerId       string         `json:"seller_id,omitempty"`       //可选	28 卖家支付宝用户ID。 当需要指定收款账号时，通过该参数传入，如果该值为空，则默认为商户签约账号对应的支付宝用户ID。 收款账号优先级规则：门店绑定的收款账户>请求传入的seller_id>商户签约账号对应的支付宝用户ID； 注：直付通和机构间联场景下seller_id无需传入或者保持跟pid一致；如果传入的seller_id与pid不一致，需要联系支付宝小二配置收款关系；  2088102146225135
	GoodsDetail    []GoodsDetails `json:"goods_detail,omitempty"`    //可选    订单包含的商品列表信息，json格式。
	ExtendParams   ExtendParams   `json:"extend_params,omitempty"`   //业务扩展参数
	BusinessParams BusinessParams `json:"business_params,omitempty"` //可选 商户传入业务信息，具体值要和支付宝约定，应用于安全，营销等参数直传场景，格式为json格式
	PromoParams    PromoParam     `json:"promo_params,omitempty"`    //可选 优惠明细参数，通过此属性补充营销参数。 注：仅与支付宝协商后可用。
}

func (t *TradePayRequest) HttpMethod() string {
	return "POST"
}

func (t *TradePayRequest) ApiParamMethod() string {
	return "alipay.trade.pay"
}

// 从响应数据中生成结构体
func (t *TradePayRequest) GenResponse(data []byte) (IResponse, error) {
	type tempResponse struct {
		AlipayTradePayResponse TradePayResponse `json:"alipay_trade_pay_response"`
	}
	var tResp tempResponse
	err := json.Unmarshal(data, &tResp)
	if err != nil {
		return nil, err
	}

	return tResp.AlipayTradePayResponse, nil
}

func (t *TradePayRequest) DoValidate() error {
	var err error
	if len(t.OutTradeNo) > 64 {
		err = fmt.Errorf("商户订单号最大长度为64字符")
		return err
	}
	for i := 0; i < len(t.Subject); i++ {
		if t.Subject[i] == '/' || t.Subject[i] == '=' || t.Subject[i] == '&' {
			err = fmt.Errorf("订单标题不可含有特殊字符")
			return err
		}
	}
	if !(t.Scene == PAY_SCENE_BAR_CODE || t.Scene == PAY_SCENE_SECURITY_CODE) {
		err = fmt.Errorf("当面付场景只能为条码支付和刷脸支付")
		return err
	}
	return nil
}

// 统一收单交易支付接口 响应 参数
type TradePayResponse struct {
	ErrorParam
	TradeNo             string          `json:"trade_no"`                        //必选	64 支付宝交易号 2013112011001004330000121536
	OutTradeNo          string          `json:"out_trade_no"`                    //必选	64 商户订单号 6823789339978248
	BuyerLogonId        string          `json:"buyer_logon_id"`                  //必选	100 买家支付宝账号 159****5620
	TotalAmount         string          `json:"total_amount"`                    //必选	11 交易金额 120.88
	ReceiptAmount       string          `json:"receipt_amount"`                  //必选	11 实收金额 88.88
	BuyerPayAmount      string          `json:"buyer_pay_amount,omitempty"`      //可选	11 买家付款的金额 8.88
	PointAmount         string          `json:"point_amount,omitempty"`          //可选	11 使用集分宝付款的金额 8.12
	InvoiceAmount       string          `json:"invoice_amount,omitempty"`        //可选	11 交易中可给用户开具发票的金额 12.50
	GmtPayment          string          `json:"gmt_payment"`                     //必选	32 交易支付时间 2014-11-27 15:45:57
	FundBillList        []TradeFundBill `json:"fund_bill_list"`                  //必选 交易支付使用的资金渠道。 只有在签约中指定需要返回资金明细，或者入参的query_options中指定时才返回该字段信息。
	StoreName           string          `json:"store_name,omitempty"`            //可选	512 发生支付交易的商户门店名称 证大五道口店
	DiscountGoodsDetail string          `json:"discount_goods_detail,omitempty"` //可选	5120 本次交易支付所使用的单品券优惠的商品优惠信息。 只有在query_options中指定时才返回该字段信息。 [{"goods_id":"STANDARD1026181538","goods_name":"雪碧","discount_amount":"100.00","voucher_id":"2015102600073002039000002D5O"}]
	BuyerUserId         string          `json:"buyer_user_id"`                   //必选	28 买家在支付宝的用户id 088101117955611
	VoucherDetailList   []VoucherDetail `json:"voucher_detail_list,omitempty"`   //可选 本交易支付时使用的所有优惠券信息。 只有在query_options中指定时才返回该字段信息。
	MdiscountAmount     string          `json:"mdiscount_amount,omitempty"`      //特殊可选	11 商家优惠金额 88.88
	DiscountAmount      string          `json:"discount_amount,omitempty"`       //特殊可选	11 平台优惠金额 88.88
}

func (t TradePayResponse) ApiParamMethod() string {
	return "alipay.trade.pay"
}
