package alipay

import "unsafe"

// 这种方法会导致生成的[]byte cap()方法不可用
func string2bytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

func bytes2string(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
