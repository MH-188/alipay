package alipay

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
)

func TestAliClient_AliPayTradPay(t *testing.T) {
	privateKey, err := LoadPrivateKeyWithPath(CLIENT_ENVIRONMENT_PROD, "./cert/alipay_secret_key_20230818.pem")
	//privateKey, err := LoadPrivateKeyWithPath(CLIENT_ENVIRONMENT_SB, "./cert/secret_key.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	publicKey, err := LoadPublicKeyWithPath("./cert/alipay_public_key_20230818.pem")
	//publicKey, err := LoadPublicKeyWithPath("./cert/public_key.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	//appId := "9021000125123456"
	appId := "123456123456123456" //绑定的第三方应用appId
	client := NewAliClient(appId, "JSON", "utf-8", "RSA2", CLIENT_ENVIRONMENT_PROD, privateKey, publicKey)

	tradePayRequest := TradePayRequest{
		OutTradeNo:  "shh202308152311", //商户订单号
		TotalAmount: 0.01,
		Subject:     "shh202308152311",    //订单标题
		AuthCode:    "289863297813534406", //支付授权码
		Scene:       PAY_SCENE_BAR_CODE,
	}
	_, err = client.AliPayTradPay("", "http://your_addr/check_sign", tradePayRequest)
	fmt.Println(err)

	return
}

func TestAliClient_AliPayTradeQuery(t *testing.T) {
	privateKey, err := LoadPrivateKeyWithPath(CLIENT_ENVIRONMENT_SB, "./cert/secret_key.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	publicKey, err := LoadPublicKeyWithPath("./cert/public_key.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	appId := "123456123456123456"
	outTradeNo := "shh202308152234"
	tradeNo := "202308182200143641050061234"
	options := make([]string, 0)
	client := NewAliClient(appId, "JSON", "utf-8", "RSA2", CLIENT_ENVIRONMENT_SB, privateKey, publicKey)

	tradeQueryRequest := TradeQueryRequest{
		OutTradeNo:   outTradeNo,
		TradeNo:      tradeNo,
		QueryOptions: options,
	}
	query, err := client.AliPayTradeQuery("", tradeQueryRequest)
	fmt.Println(query, err)
}

func TestAliClient_AliPayTradRefund(t *testing.T) {
	privateKey, err := LoadPrivateKeyWithPath(CLIENT_ENVIRONMENT_SB, "./cert/secret_key.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	publicKey, err := LoadPublicKeyWithPath("./cert/public_key.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	appId := "123456123456123456"
	outTradeNo := "shh202308152311"
	client := NewAliClient(appId, "JSON", "utf-8", "RSA2", CLIENT_ENVIRONMENT_SB, privateKey, publicKey)

	tradeRefundRequest := TradeRefundRequest{
		OutTradeNo:   outTradeNo,
		RefundAmount: 100,
		RefundReason: "只买一半",
		OutRequestNo: "TK0000000",
		//query_options: []string{"deposit_back_info"},
	}
	refund, err := client.AliPayTradeRefund("", "http://47.103.50.40:8080/check_sign", tradeRefundRequest)

	bytes, _ := json.Marshal(refund)
	fmt.Println(string(bytes), err)
}

func TestAliClient_AliPayTradeFastPayRefundQuery(t *testing.T) {
	privateKey, err := LoadPrivateKeyWithPath(CLIENT_ENVIRONMENT_SB, "./cert/secret_key.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	publicKey, err := LoadPublicKeyWithPath("./cert/public_key.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	appId := "123456123456123456"
	outTradeNo := "shh202308152310"
	client := NewAliClient(appId, "JSON", "utf-8", "RSA2", CLIENT_ENVIRONMENT_SB, privateKey, publicKey)

	tradeFastPayRefundQueryRequest := TradeFastPayRefundQueryRequest{
		OutTradeNo:   outTradeNo,
		OutRequestNo: "TK0000001",
	}
	query, err := client.AliPayTradeFastPayRefundQuery("", tradeFastPayRefundQueryRequest)
	bytes, _ := json.Marshal(query)

	fmt.Println(string(bytes), err)
}

func TestAsyncRequestCheckSign(t *testing.T) {
	privateKey, err := LoadPrivateKeyWithPath(CLIENT_ENVIRONMENT_SB, "./cert/secret_key.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	publicKey, err := LoadPublicKeyWithPath("./cert/public_key.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	appId := "123456123456123456"
	client := NewAliClient(appId, "JSON", "utf-8", "RSA2", CLIENT_ENVIRONMENT_SB, privateKey, publicKey)

	bytes := []byte("gmt_create=2023-08-19+11%3A08%3A14&charset=utf-8&seller_email=cyfegy7191%40sandbox.com&subject=shh202308152306&sign=TccV1D1aPzlPS6VuhLAVCEHRUNjPrKH3vW%2FlbF71G4bt4r9bX%2FReWO99tz2KYxlw5Wp6qcbT4i2A%2FtXwJLn5eSU2N5x8Gx2JJqLltLqS%2BIgEWASbEii43GEl3ZQC2oxH5DfCSnrqln5aNoFD7XZJzKneTvNmbkbwSNa3wDSTxL%2FpoOepixF9vYKnVquhntujoIhoYISWk4h6CFCdRE9GLj4geT6spEPJ8A2Sb%2FiStSr4U%2BWOeuowAEpNjEeWBQYdVX%2BXZZ3%2Bg5fbRUB5CVqdNObUZJj%2BD2cVCplyz8JSv3HK9Z78xHaA0vfLK7pveBnb3OwcTa7vGJr19ldCWH3CNQ%3D%3D&buyer_id=2088722008336414&invoice_amount=1000.00&notify_id=2023081901222110816136410500676618&fund_bill_list=%5B%7B%22amount%22%3A%221000.00%22%2C%22fundChannel%22%3A%22ALIPAYACCOUNT%22%7D%5D&notify_type=trade_status_sync&trade_status=TRADE_SUCCESS&receipt_amount=1000.00&app_id=9021000125631019&buyer_pay_amount=1000.00&sign_type=RSA2&seller_id=2088721008358437&gmt_payment=2023-08-19+11%3A08%3A15&notify_time=2023-08-19+11%3A08%3A17&version=1.0&out_trade_no=shh202308152306&total_amount=1000.00&trade_no=2023081922001436410500638614&auth_app_id=9021000125631019&buyer_logon_id=gdrrob7568%40sandbox.com&point_amount=0.00")

	reqParam, err := client.AsyncRequestGetParamAndCheckSign(bytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(reqParam, "success")
}

func TestAliClient_AliPayTradRefundOnSameTime(t *testing.T) {
	privateKey, err := LoadPrivateKeyWithPath(CLIENT_ENVIRONMENT_SB, "./cert/secret_key.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	publicKey, err := LoadPublicKeyWithPath("./cert/public_key.pem")
	if err != nil {
		fmt.Println(err)
		return
	}
	appId := "123456123456123456"
	outTradeNo := "shh202308152311"
	client := NewAliClient(appId, "JSON", "utf-8", "RSA2", CLIENT_ENVIRONMENT_SB, privateKey, publicKey)

	tradeRefundRequest := TradeRefundRequest{
		OutTradeNo:   outTradeNo,
		RefundAmount: 100,
		RefundReason: "只买一半",
		OutRequestNo: "TK0000002",
		//query_options: []string{"deposit_back_info"},
	}
	//refund, err := client.AliPayTradeRefund("http://47.103.50.40:8080/check_sign", tradeRefundRequest)
	//bytes, _ := json.Marshal(refund)
	//fmt.Println(string(bytes), err)

	refundFunc := func(wg *sync.WaitGroup) {
		refund, err := client.AliPayTradeRefund("", "http://47.103.50.40:8080/check_sign", tradeRefundRequest)
		bytes, _ := json.Marshal(refund)
		fmt.Println(string(bytes), err)
		wg.Done()
	}

	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go refundFunc(&wg)
	}

	wg.Wait()
}
