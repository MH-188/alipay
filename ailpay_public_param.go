package alipay

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-querystring/query"
	"net/url"
	"sort"
	"strings"
)

/*
 * 错误处理
 */
type ErrorParam struct {
	Msg     string `json:"msg"`
	Code    string `json:"code"`
	SubCode string `json:"sub_code,omitempty"`
	SubMsg  string `json:"sub_msg,omitempty"`
}

func (e ErrorParam) GetBadResponseDesc() string {
	return e.SubMsg
}

type RBuffer []byte

// IsBadResponse 响应code校验
func (rb RBuffer) IsBadResponse() bool {
	// "code":"10000"
	str := bytes2string(rb) //可以更换成共享内存
	index := strings.Index(str, "code")
	code := str[index+7 : index+12]
	if code != RESPONSE_SUCCESS {
		return true
	}
	return false
}

// GetSyncRespSignData 同步响应待验证字串
func (rb RBuffer) GetSyncRespSignData() (string, error) {
	str := bytes2string(rb)
	indexR := strings.Index(str, "_response")
	if indexR < 0 {
		err := errors.New("数据错误，找不到_response结构")
		return "", err
	}

	indexC := strings.Index(str, `,"alipay_cert_sn":`)
	if indexC > 0 && indexC > indexR {
		return str[indexR+11 : indexC], nil
	}

	indexS := strings.Index(str, `,"sign":`)
	if indexS < 0 {
		err := errors.New("数据错误，找不到sign结构")
		return "", err
	}

	return str[indexR+11 : indexS], nil
}

// GetSyncRespSign 获取响应中的sign
func (rb RBuffer) GetSyncRespSign() (string, error) {
	str := bytes2string(rb)
	indexBegin := strings.Index(str, `,"sign":`)
	if indexBegin < 0 {
		err := errors.New("数据错误，找不到sign结构")
		return "", err
	}

	str2 := str[indexBegin+9:]
	indexEnd := strings.Index(str2, `"`)
	if indexEnd < 0 {
		err := errors.New("数据错误，sign数据结构异常")
		return "", err
	}
	return str2[:indexEnd], nil
}

// GetAsyncRequestSignStr 获取异步通知的待验证参数字符串 和 签名
func (rb RBuffer) GetAsyncRequestSignStr() (string, string, error) {
	asyncRequest := make(map[string]string)
	requestStr := string(rb)
	paramStrs := strings.Split(requestStr, "&")
	for i := 0; i < len(paramStrs); i++ {
		kvStrs := strings.Split(paramStrs[i], "=")
		if len(kvStrs) != 2 {
			err := errors.New("参数异常，无法解析：" + requestStr)
			return "", "", err
		}

		valDecode, err := url.QueryUnescape(kvStrs[1])
		if err != nil {
			return "", "", err
		}

		asyncRequest[kvStrs[0]] = valDecode

	}

	keys := make([]string, 0, len(asyncRequest))
	for k, _ := range asyncRequest {
		if k == "sign" || k == "sign_type" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	feildStrs := make([]string, 0, len(keys))
	for i := 0; i < len(keys); i++ {
		str := fmt.Sprintf("%s=%s", keys[i], asyncRequest[keys[i]])
		feildStrs = append(feildStrs, str)
	}

	sign, ok := asyncRequest["sign"]
	if !ok {
		err := errors.New("sign类型异常，请检查响应参数")
		return "", "", err
	}
	return strings.Join(feildStrs, "&"), sign, nil
}

/*
 * 公共参数
 */

// 公共请求参数
type CommonRequestParam struct {
	AppId        string `json:"app_id" url:"app_id"`                                     // 必选	最大长度32	支付宝分配给开发者的应用ID 2014072300007148
	Method       string `json:"method" url:"method"`                                     // 必选	128 		接口名称 alipay.trade.page.pay
	Format       string `json:"format,omitempty" url:"format"`                           // 可选	40 			仅支持JSON JSON
	Charset      string `json:"charset" url:"charset"`                                   // 必选	10 			请求使用的编码格式，如utf-8,gbk,gb2312等
	SignType     string `json:"sign_type" url:"sign_type"`                               // 必选	10 			商户生成签名字符串所使用的签名算法类型，目前支持RSA2和RSA，推荐使用RSA2
	Sign         string `json:"sign,omitempty" url:"sign"`                               // 必选	344 		商户请求参数的签名串
	Timestamp    string `json:"timestamp" url:"timestamp"`                               // 必选	19 			发送请求的时间，格式"yyyy-MM-dd HH:mm:ss" 2014-07-24 03:07:50
	Version      string `json:"version" url:"version"`                                   // 必选 3 			调用的接口版本，固定为：1.0
	NotifyUrl    string `json:"notify_url,omitempty" url:"notify_url,omitempty"`         // 可选	256	 		支付宝服务器主动通知商户服务器里指定的页面http/https路径。
	AppAuthToken string `json:"app_auth_token,omitempty" url:"app_auth_token,omitempty"` // 可选	40 			详见应用授权概述
	BizContent   string `json:"biz_content" url:"biz_content"`                           // 必选	无长度限制 	请求参数的集合，最大长度不限，除公共参数外所有请求参数都必须放在这个参数中传递，具体参照各产品快速接入文档
	//AppCertSn        string `json:"app_cert_sn" url:"app_cert_sn,omitempty"`                 // 可选				具体参照各产品快速接入文档
	//AlipayRootCertSn string `json:"alipay_root_cert_sn" url:"alipay_root_cert_sn,omitempty"` // 可选				具体参照各产品快速接入文档
}

// 公共响应参数
type CommonResponseParam struct {
	Code    string `json:"code"`     //必选	网关返回码,详见文档 40004
	Msg     string `json:"msg"`      //必选	网关返回码描述,详见文档 Business Failed
	SubCode string `json:"sub_code"` //可选  业务返回码，参见具体的API接口文档 ACQ.TRADE_HAS_SUCCESS
	SubMsg  string `json:"sub_msg"`  //可选  业务返回码描述，参见具体的API接口文档 交易已被支付
	Sign    string `json:"sign"`     //必选	签名,详见文档 DZXh8eeTuAHoYE3w1J+POiPhfDxOYBfUNn1lkeT/V7P4zJdyojWEa6IZs6Hz0yDW5Cp/viufUb5I0/V5WENS3OYR8zRedqo6D+fUTdLHdc+EFyCkiQhBxIzgngPdPdfp1PIS7BdhhzrsZHbRqb7o4k3Dxc+AAnFauu4V6Zdwczo=
}

// struct 转 map
func (crp *CommonRequestParam) structToMapJson() (map[string]interface{}, error) {
	data, err := json.Marshal(crp)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var paramMap map[string]interface{}
	err = json.Unmarshal(data, &paramMap)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return paramMap, err
}

// 生成参数拼接字符串: 对参数的key排序，使用key=value&key1=value1格式拼接参数
func (crp *CommonRequestParam) genParamString(paramMap map[string]interface{}) (string, error) {
	paramKeys := make([]string, 0, len(paramMap))
	for k, _ := range paramMap {
		paramKeys = append(paramKeys, k)
	}
	// 对key排序
	sort.Strings(paramKeys)

	// 拼接字串用于签名: 此处拼接无sign字段
	var paramBuffer bytes.Buffer
	for i := 0; i < len(paramKeys); i++ {
		paramBuffer.WriteString(paramKeys[i])
		paramBuffer.WriteString("=")
		paramBuffer.WriteString(fmt.Sprintf("%v", paramMap[paramKeys[i]]))
		if i != len(paramKeys)-1 {
			paramBuffer.WriteString("&")
		}
	}
	paramBuffer.Truncate(paramBuffer.Len())
	return paramBuffer.String(), nil
}

// 设置BigContent参数
func (crp *CommonRequestParam) SetRequestBizContent(req IRequester) (string, error) {
	bytes, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	crp.BizContent = string(bytes)
	return crp.BizContent, nil
}

// 设置Sign
func (crp *CommonRequestParam) SetRequestSign(privateKey *rsa.PrivateKey) (string, error) {
	if len(crp.BizContent) < 1 {
		err := errors.New("big_content为空，请先传入参数")
		return "", err
	}
	//先使用简单方法实现structToMap 并不是效率最高的
	paramMap, err := crp.structToMapJson()
	if err != nil {
		return "", err
	}

	paramStr, err := crp.genParamString(paramMap)
	//fmt.Println("用于生成签名的参数字串：", paramStr)

	// 生成签名
	signed := Rsa2PrivateSign(paramStr, privateKey, crypto.SHA256)
	//fmt.Println("签名：", signed)
	crp.Sign = signed
	return signed, nil
}

// 生成请求参数字符串: 该方法的执行必须在SetRequestSign方法后
func (crp *CommonRequestParam) GenRequestParamStr() (string, error) {
	if len(crp.Sign) < 1 {
		err := errors.New("sign参数为空，请先生成sign")
		return "", err
	}
	paramMap, err := crp.structToMapJson()
	if err != nil {
		return "", err
	}

	_, err = crp.genParamString(paramMap)
	if err != nil {
		return "", err
	}

	//fmt.Println("用于发起请求的参数字串：", paramStr)

	paramUrlEncode, err := query.Values(crp)
	if err != nil {
		return "", err
	}

	// 发起请求前要对内容进行url encode
	return paramUrlEncode.Encode(), err
}
