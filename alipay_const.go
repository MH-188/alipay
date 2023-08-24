package alipay

// 错误码
const (
	RESPONSE_SUCCESS string = "10000"
)

const (
	// 请求环境
	CLIENT_ENVIRONMENT_SB   string = "sandbox"
	CLIENT_ENVIRONMENT_PROD string = "production"

	// 网关地址
	GATEWAY_URL_SB   string = "https://openapi-sandbox.dl.alipaydev.com/gateway.do"
	GATEWAY_URL_PROD string = "https://openapi.alipay.com/gateway.do"
)

// 当面付支付场景
const (
	PAY_SCENE_BAR_CODE      string = "bar_code"      //当面付条码支付场景(默认)
	PAY_SCENE_SECURITY_CODE string = "security_code" //当面付刷脸支付场景，对应的auth_code为fp开头的刷脸标识串
)

// 当面付查询参数
const (
	QUERY_OPTIONS_TRADE_SETTLE_INFO      string = "trade_settle_info"
	QUERY_OPTIONS_FUND_BILL_LIST         string = "fund_bill_list"
	QUERY_OPTIONS_VOUCHER_DETAIL_LIST    string = "voucher_detail_list"
	QUERY_OPTIONS_DISCOUNT_GOODS_DETAIL  string = "discount_goods_detail"
	QUERY_OPTIONS_MDISCOUNT_AMOUNT       string = "mdiscount_amount"
	QUERY_OPTIONS_MEDICAL_INSURANCE_INFO string = "medical_insurance_info"
)

// 换取应用令牌的两种方式
const (
	AUTH_TOKEN_APP_SUTHORIZATION_CODE string = "authorization_code" //使用应用授权码换取应用授权令牌app_auth_token。
	AUTH_TOKEN_APP_RFRESH_TOKEN       string = "refresh_token"      //使用app_refresh_token刷新获取新授权令牌。
)

// 查询交易订单的状态
const (
	TRADE_STATUS_WAIT_BUYER_PAY string = "WAIT_BUYER_PAY" //交易创建，等待买家付款
	TRADE_STATUS_TRADE_CLOSED   string = "TRADE_CLOSED"   //未付款交易超时关闭，或支付完成后全额退款
	TRADE_STATUS_TRADE_SUCCESS  string = "TRADE_SUCCESS"  //交易支付成功
	TRADE_STATUS_TRADE_FINISHED string = "TRADE_FINISHED" //交易结束，不可退款
)

// 撤销订单返回状态
const (
	TRADE_CANCEL_ACTION_CLOSE  string = "close"
	TRADE_CANCEL_ACTION_REFUND string = "refund"
)

// 退款查询返回状态
const (
	TRADE_REFUND_QUERY_STATUS_SUCCESS = "REFUND_SUCCESS"
)
