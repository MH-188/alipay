package alipay

import (
	"encoding/json"
	"fmt"
)

/*
 * 换取授权令牌
 * 换取应用授权令牌
 * alipay.open.auth.token.app
 */

// 换取应用授权令牌 请求 code的值在商家授权后支付宝回调参数中获取，回调示例：http://example.com/doc/toAuthPage.html?app_id=2021000000000318&source=alipay_app_auth&application_type=TINYAPP,WEBAPP&app_auth_code=P1798b23682e34d96859afa000000003
// 用户需要先从回调参数中获取到 app_auth_code 后再调取本接口换取 app_auth_token
type OpenAuthTokenAppRequest struct {
	GrantType    string `json:"grant_type"`              //必选	20 授权方式。支持： authorization_code：使用应用授权码换取应用授权令牌app_auth_token。 refresh_token：使用app_refresh_token刷新获取新授权令牌。 authorization_code或者refresh_token
	Code         string `json:"code,omitempty"`          //可选	40 应用授权码，传入应用授权后得到的 app_auth_code。 说明： grant_type 为 authorization_code 时，本参数必填； grant_type 为 refresh_token 时不填。 1cc19911172e4f8aaa509c8fb5d12F56
	RefreshToken string `json:"refresh_token,omitempty"` //可选	40 刷新令牌，上次换取访问令牌时得到。本参数在 grant_type 为 authorization_code 时不填；为 refresh_token 时必填，且该值来源于此接口的返回值 app_refresh_token（即至少需要通过 grant_type=authorization_code 调用此接口一次才能获取）。 201509BBdcba1e3347de4e75ba3fed2c9abebE36
}

func (o *OpenAuthTokenAppRequest) HttpMethod() string {
	return "POST"
}

func (o *OpenAuthTokenAppRequest) ApiParamMethod() string {
	return "alipay.open.auth.token.app"
}

// 从响应数据中生成结构体
func (o *OpenAuthTokenAppRequest) GenResponse(data []byte) (IResponse, error) {
	type tempResponse struct {
		OpenAuthTokenAppResponse OpenAuthTokenAppResponse `json:"alipay_open_auth_token_app_response"`
	}
	var tResp tempResponse
	err := json.Unmarshal(data, &tResp)
	if err != nil {
		return nil, err
	}

	return tResp.OpenAuthTokenAppResponse, nil
}

func (o *OpenAuthTokenAppRequest) DoValidate() error {
	if o.GrantType == AUTH_TOKEN_APP_SUTHORIZATION_CODE && len(o.Code) < 1 {
		err := fmt.Errorf("换取应用令牌方式为 %s 时，code不能为空", AUTH_TOKEN_APP_SUTHORIZATION_CODE)
		return err
	}
	if o.GrantType == AUTH_TOKEN_APP_RFRESH_TOKEN && len(o.RefreshToken) < 1 {
		err := fmt.Errorf("换取应用令牌方式为 %s 时，refresh_token 不能为空", AUTH_TOKEN_APP_RFRESH_TOKEN)
		return err
	}
	return nil
}

// 换取应用授权令牌 响应
type OpenAuthTokenAppResponse struct {
	UserId          string `json:"user_id"`           //必选	16 授权商户的user_id 2088102150527498
	AuthAppId       string `json:"auth_app_id"`       //必选	20 授权商户的appid 2013121100055554
	AppAuthToken    string `json:"app_auth_token"`    //必选	40 应用授权令牌 201509BBeff9351ad1874306903e96b91d248A36
	AppRefreshToken string `json:"app_refresh_token"` //必选	40 刷新令牌 201509BBdcba1e3347de4e75ba3fed2c9abebE36
	ReExpiresIn     string `json:"re_expires_in"`     //必选	16 刷新令牌的有效时间（从接口调用时间作为起始时间），单位到秒 123456
	//expires_in	string	//必选	16 该字段已作废，应用令牌长期有效，接入方不需要消费该字段 123456
}

func (t OpenAuthTokenAppResponse) ApiParamMethod() string {
	return "alipay.open.auth.token.app"
}
