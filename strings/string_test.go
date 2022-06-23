package strings

import (
	"bytes"
	"testing"
)

// 测试强转换功能
func TestBytes2String(t *testing.T) {
	x := []byte("Hello Gopher!")
	y := Bytes2String(x)
	z := string(x)

	if y != z {
		t.Fail()
	}
}

// 测试强转换功能
func TestString2Bytes(t *testing.T) {
	x := "Hello Gopher!"
	y := String2Bytes(x)
	z := []byte(x)

	if !bytes.Equal(y, z) {
		t.Fail()
	}
}

// 测试标准转换string()性能
func Benchmark_NormalBytes2String(b *testing.B) {
	x := []byte("Hello Gopher! Hello Gopher! Hello Gopher!")
	for i := 0; i < b.N; i++ {
		_ = string(x)
	}
}

// 测试强转换[]byte到string性能
func Benchmark_Byte2String(b *testing.B) {
	x := []byte("Hello Gopher! Hello Gopher! Hello Gopher!")
	for i := 0; i < b.N; i++ {
		_ = Bytes2String(x)
	}
}

// 测试标准转换[]byte性能
func Benchmark_NormalString2Bytes(b *testing.B) {
	x := "Hello Gopher! Hello Gopher! Hello Gopher!"
	for i := 0; i < b.N; i++ {
		_ = []byte(x)
	}
}

// 测试强转换string到[]byte性能
func Benchmark_String2Bytes(b *testing.B) {
	x := "Hello Gopher! Hello Gopher! Hello Gopher!"
	for i := 0; i < b.N; i++ {
		_ = String2Bytes(x)
	}
}
