package alipay

type IRequester interface {
	HttpMethod() string                    //http方法
	ApiParamMethod() string                //请求的接口名称
	GenResponse([]byte) (IResponse, error) //获取响应数据
	DoValidate() error                     //检查请求体数据
}

// 返回结果校验
type IRespVerify interface {
	IsBadResponse() bool //响应code校验
	//VerifySign() error   //响应签名校验
}

type IResponse interface {
	ApiParamMethod() string
}
