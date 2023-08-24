package alipay

/*
 * 付款码支付
 * 收单退款冲退完成通知
 * alipay.trade.refund.depositback.completed
 */

// 支付宝服务器发送的通知 请求
type TradeRefundDepositBackCompletedRequest struct {
	NotifyId     string `json:"notify_id"`              //必选	50 通知ID 5608cccc09ddb39d41c2e3c06e3d9fejh2
	UtcTimestamp string `json:"utc_timestamp"`          //必选	13 消息发送时的服务端时间 1514210452731
	MsgMethod    string `json:"msg_method"`             //必选	100 消息接口名称 alipay.trade.refund.depositback.completed
	AppId        string `json:"app_id"`                 //必选	20 消息接受方的应用id 2014060600164699
	MsgType      string `json:"msg_type,omitempty"`     //可选	5 消息类型。目前支持类型：sys：系统消息；usr，用户消息；app，应用消息 sys
	MsgAppId     string `json:"msg_app_id,omitempty"`   //可选	20 消息归属方的应用id。应用消息时非空 2016032301002387
	Version      string `json:"version"`                //必选	5 版本号(1.1版本为标准消息) 1.0或者1.1
	BizContent   string `json:"biz_content"`            //必选    消息报文 参见具体的消息接口文档
	Sign         string `json:"sign"`                   //必选    签名 WcO+t3D8Kg71dTlKwN7r9PzUOXeaBJwp8/FOuSxcuSkXsoVYxBpsAidprySCjHCjmaglNcjoKJQLJ28/Asl93joTW39FX6i07lXhnbPknezAlwmvPdnQuI01HZsZF9V1i6ggZjBiAd5lG8bZtTxZOJ87ub2i9GuJ3Nr/NUc9VeY=
	SignType     string `json:"sign_type"`              //必选	10 签名类型 RSA2
	EncryptType  string `json:"encrypt_type,omitempty"` //可选	10 加密算法 AES
	Charset      string `json:"charset"`                //必选	10 编码集，该字符集为验签和解密所需要的字符集 UTF-8
	NotifyType   string `json:"notify_type,omitempty"`  //可选	20 [1.0版本老协议参数]通知类型，1.1接口没有该参数 trade_status
	NotifyTime   string `json:"notify_time,omitempty"`  //可选	19 [1.0版本老协议参数]通知时间，北京时区，时间格式为：yyyy-MM-dd HH:mm:ss，如果服务器部署在国外请注意做时区转换。若version=1.1则可以使用utc_timestamp识别时间  1970-01-01 00:00:00
	AuthAppId    string `json:"auth_app_id,omitempty"`  //可选	20 [1.0版本老协议参数]授权方的应用id 2016032301002387
}

// 需要响应给支付宝服务器的内容 响应
const (
	TRADE_REFUND_RESPONSE_FAIL    string = "fail"
	TRADE_REFUND_RESPONSE_SUCCESS string = "success"
)
