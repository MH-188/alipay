package alipay

import (
	"fmt"
	"testing"
)

func TestGetAsyncRequestSignStr(t *testing.T) {
	rbuff := RBuffer("gmt_create=2023-08-19+11%3A08%3A14&charset=utf-8&seller_email=cyfegy7191%40sandbox.com&subject=shh202308152306&sign=TccV1D1aPzlPS6VuhLAVCEHRUNjPrKH3vW%2FlbF71G4bt4r9bX%2FReWO99tz2KYxlw5Wp6qcbT4i2A%2FtXwJLn5eSU2N5x8Gx2JJqLltLqS%2BIgEWASbEii43GEl3ZQC2oxH5DfCSnrqln5aNoFD7XZJzKneTvNmbkbwSNa3wDSTxL%2FpoOepixF9vYKnVquhntujoIhoYISWk4h6CFCdRE9GLj4geT6spEPJ8A2Sb%2FiStSr4U%2BWOeuowAEpNjEeWBQYdVX%2BXZZ3%2Bg5fbRUB5CVqdNObUZJj%2BD2cVCplyz8JSv3HK9Z78xHaA0vfLK7pveBnb3OwcTa7vGJr19ldCWH3CNQ%3D%3D&buyer_id=2088722008336414&invoice_amount=1000.00&notify_id=2023081901222110816136410500676618&fund_bill_list=%5B%7B%22amount%22%3A%221000.00%22%2C%22fundChannel%22%3A%22ALIPAYACCOUNT%22%7D%5D&notify_type=trade_status_sync&trade_status=TRADE_SUCCESS&receipt_amount=1000.00&app_id=9021000125631019&buyer_pay_amount=1000.00&sign_type=RSA2&seller_id=2088721008358437&gmt_payment=2023-08-19+11%3A08%3A15&notify_time=2023-08-19+11%3A08%3A17&version=1.0&out_trade_no=shh202308152306&total_amount=1000.00&trade_no=2023081922001436410500638614&auth_app_id=9021000125631019&buyer_logon_id=gdrrob7568%40sandbox.com&point_amount=0.00")
	str, s, err := rbuff.GetAsyncRequestSignStr()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(str, s)
}
