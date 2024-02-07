package utils

import "unsafe"

// ToString convert b to string without copy
// !!!注意!!! ToString转换后,原字节数组与转换后的字符串避免进行写修改,否则会导致内存崩溃
func ToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// ToBytes convert s to []byte without copy
// !!!注意!!! ToBytes转换后,原数据和转换后的数组避免进行写操作,否则会导致内存崩溃
func ToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
