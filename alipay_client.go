package alipay

import (
	"crypto"
	"crypto/rsa"
	"fmt"
	"net/http"
	"time"
)

type AliClient struct {
	appId      string //应用id
	format     string //JSON
	signType   string //商户生成签名字符串所使用的签名算法类型，目前支持RSA2和RSA，推荐使用RSA2
	charset    string //编码格式utf-8
	aliGateway string //支付宝网关地址，沙箱：https://openapi-sandbox.dl.alipaydev.com/gateway.do 正式：https://openapi.alipay.com/gateway.do

	privateKey  *rsa.PrivateKey //商户私钥
	publicKey   *rsa.PublicKey  //支付宝公钥
	environment string          //"sandbox"沙箱；"production"生产环境
	httpClient  *http.Client    //http客户端
}

func NewAliClient(appId, format, charset, signType, environment string, aliPrivateKey *rsa.PrivateKey, aliPublicKey *rsa.PublicKey) *AliClient {
	client := AliClient{
		appId:       appId,
		format:      format,
		signType:    signType,
		environment: environment,
		charset:     charset,
		privateKey:  aliPrivateKey,
		publicKey:   aliPublicKey,
		httpClient:  http.DefaultClient, //先使用默认客户端
	}
	if environment == CLIENT_ENVIRONMENT_PROD {
		client.aliGateway = GATEWAY_URL_PROD
	} else if environment == CLIENT_ENVIRONMENT_SB {
		client.aliGateway = GATEWAY_URL_SB
	}
	return &client
}

// SyncVerifySign 同步响应签名校验: scene验签场景 sign待验证的签名 buff待验证的内容 sertSn支付宝公钥证书序列号(公钥证书模式下，同步响应时使用)
func (client *AliClient) SyncVerifySign(sign, buff string) error {
	//同步响应验签
	err := RSAVerifyWithKey([]byte(buff), []byte(sign), client.publicKey, crypto.SHA256)
	if err != nil {
		return err
	}

	return nil
}

// 发起请求
func (client *AliClient) DoRequest(appId, appAuthToken, notifyUrl string, apiReq IRequester) (IResponse, error) {
	commonReq := CommonRequestParam{
		AppId:        appId,
		Method:       apiReq.ApiParamMethod(),
		Format:       client.format,
		Charset:      client.charset,
		SignType:     client.signType,
		Sign:         "",
		Timestamp:    time.Now().Format("2006-01-02 15:04:05"),
		Version:      "1.0",
		NotifyUrl:    notifyUrl,
		BizContent:   "",           //除公共参数外的所有参数
		AppAuthToken: appAuthToken, //灵活数据，要作为参数传入
	}
	_, err := commonReq.SetRequestBizContent(apiReq)
	if err != nil {
		return nil, err
	}

	// 设置请求签名
	_, err = commonReq.SetRequestSign(client.privateKey)
	if err != nil {
		return nil, err
	}

	// 构建请求参数
	paramStr, err := commonReq.GenRequestParamStr()
	if err != nil {
		return nil, err
	}

	// 发起请求
	resp, err := client.doRequest(apiReq, []byte(paramStr))
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// 处理http层数据
func (client *AliClient) doRequest(apiReq IRequester, param []byte) (IResponse, error) {
	// 参数校验
	err := apiReq.DoValidate()
	if err != nil {
		return nil, err
	}

	// 发起请求
	method := apiReq.HttpMethod()
	resp, err := HttpDoRequest(method, client.aliGateway, param)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(resp))

	respBuff := RBuffer(resp)
	// 准备验签数据
	awaitSignData, err := respBuff.GetSyncRespSignData()
	if err != nil {
		return nil, err
	}

	sign, err := respBuff.GetSyncRespSign()
	if err != nil {
		return nil, err
	}

	// 验签
	err = client.SyncVerifySign(sign, awaitSignData)
	if err != nil {
		return nil, err
	}

	// 处理业务层错误
	//if respBuff.IsBadResponse() {
	//	type tempResponse struct {
	//		AlipayTradePayResponse ErrorParam `json:"alipay_trade_pay_response"`
	//	}
	//	var errResp tempResponse
	//	err = json.Unmarshal(resp, &errResp)
	//	if err != nil {
	//		return nil, err
	//	}
	//	err = errors.New(errResp.AlipayTradePayResponse.GetBadResponseDesc())
	//	return nil, err
	//}

	apiResp, err := apiReq.GenResponse(resp)
	if err != nil {
		return nil, err
	}
	fmt.Println(apiResp)

	err = apiResp.IsBadResponse()
	if err != nil {
		return nil, err
	}

	return apiResp, nil
}

// AsyncRequestCheckSign 异步通知请求验签
func (client *AliClient) AsyncRequestGetParamAndCheckSign(reqBuff []byte) (map[string]string, error) {
	reqParam, signData, sign, err := RBuffer(reqBuff).GetAsyncRequestSignStr()
	if err != nil {
		return reqParam, err
	}

	err = RSAVerifyWithKey([]byte(signData), []byte(sign), client.publicKey, crypto.SHA256)
	if err != nil {
		return reqParam, err
	}
	return reqParam, nil
}

// AliPayOpenAuthTokenApp 换取应用授权令牌
func (client *AliClient) AliPayOpenAuthTokenApp(openAuthTokenAppRequest OpenAuthTokenAppRequest) (OpenAuthTokenAppResponse, error) {
	resp, err := client.DoRequest(client.appId, "", "", &openAuthTokenAppRequest)
	if err != nil {
		return OpenAuthTokenAppResponse{}, err
	}

	tradePayResponse, ok := resp.(OpenAuthTokenAppResponse)
	if !ok {
		err = fmt.Errorf("响应结构不匹配方法%s", openAuthTokenAppRequest.ApiParamMethod())
		return OpenAuthTokenAppResponse{}, err
	}

	return tradePayResponse, nil
}

// AliPayTradPay 统一收单交易支付接口(alipay.trade.pay): 收银员使用扫码枪读取用户手机支付宝支付码
func (client *AliClient) AliPayTradPay(appAuthToken string, notifyUrl string, tradePayRequest TradePayRequest) (TradePayResponse, error) {
	resp, err := client.DoRequest(client.appId, appAuthToken, notifyUrl, &tradePayRequest)
	if err != nil {
		return TradePayResponse{}, err
	}

	tradePayResponse, ok := resp.(TradePayResponse)
	if !ok {
		err = fmt.Errorf("响应结构不匹配方法%s", tradePayRequest.ApiParamMethod())
		return TradePayResponse{}, err
	}

	return tradePayResponse, nil
}

// AliPayTradeQuery 统一收单交易查询(alipay.trade.query)
func (client *AliClient) AliPayTradeQuery(appAuthToken string, tradeQueryRequest TradeQueryRequest) (TradeQueryResponse, error) {
	resp, err := client.DoRequest(client.appId, appAuthToken, "", &tradeQueryRequest)
	if err != nil {
		return TradeQueryResponse{}, err
	}

	tradeQueryResponse, ok := resp.(TradeQueryResponse)
	if !ok {
		err = fmt.Errorf("响应结构不匹配方法%s", tradeQueryRequest.ApiParamMethod())
		return TradeQueryResponse{}, err
	}
	return tradeQueryResponse, nil
}

// AliPayTradeRefund 统一收单交易退款接口(alipay.trade.refund) 只有在支付时填写回调地址，才会收到退款通知
func (client *AliClient) AliPayTradeRefund(appAuthToken string, notifyUrl string, tradeRefundRequest TradeRefundRequest) (TradeRefundResponse, error) {
	resp, err := client.DoRequest(client.appId, appAuthToken, notifyUrl, &tradeRefundRequest)
	if err != nil {
		return TradeRefundResponse{}, err
	}

	tradeRefundResponse, ok := resp.(TradeRefundResponse)
	if !ok {
		err = fmt.Errorf("响应结构不匹配方法%s", tradeRefundRequest.ApiParamMethod())
		return TradeRefundResponse{}, err
	}

	return tradeRefundResponse, nil
}

// AliPayTradeFastPayRefundQuery 统一收单交易退款查询(alipay.trade.fastpay.refund.query)
func (client *AliClient) AliPayTradeFastPayRefundQuery(appAuthToken string, tradeFastPayRefundQueryRequest TradeFastPayRefundQueryRequest) (TradeFastPayRefundQueryResponse, error) {
	resp, err := client.DoRequest(client.appId, appAuthToken, "", &tradeFastPayRefundQueryRequest)
	if err != nil {
		return TradeFastPayRefundQueryResponse{}, err
	}
	tradeFastPayRefundQueryResponse, ok := resp.(TradeFastPayRefundQueryResponse)
	if !ok {
		err = fmt.Errorf("响应结构不匹配方法%s", tradeFastPayRefundQueryRequest.ApiParamMethod())
		return TradeFastPayRefundQueryResponse{}, err
	}
	return tradeFastPayRefundQueryResponse, nil
}

// AliPayTradeCancel 统一收单交易撤销接口(alipay.trade.cancel)
func (client *AliClient) AliPayTradeCancel(appAuthToken string, tradeCancelRequest TradeCancelRequest) (TradeCancelResponse, error) {
	resp, err := client.DoRequest(client.appId, appAuthToken, "", &tradeCancelRequest)
	if err != nil {
		return TradeCancelResponse{}, err
	}
	tradeCancelResponse, ok := resp.(TradeCancelResponse)
	if !ok {
		err = fmt.Errorf("响应结构不匹配方法%s", tradeCancelRequest.ApiParamMethod())
		return TradeCancelResponse{}, err
	}
	return tradeCancelResponse, nil
}
