package hotpatch

import (
	"fmt"
	"reflect"
	"runtime"
	"syscall"
	"unsafe"
)

type Patches struct {
	originals    map[uintptr][]byte
	values       map[reflect.Value]reflect.Value
	valueHolders map[reflect.Value]reflect.Value
}

func create() *Patches {
	return &Patches{originals: make(map[uintptr][]byte), values: make(map[reflect.Value]reflect.Value), valueHolders: make(map[reflect.Value]reflect.Value)}
}

func NewPatches() *Patches {
	return create()
}

func ApplyFunc(target, double interface{}) *Patches {
	return create().ApplyFunc(target, double)
}

func (p *Patches) ApplyFunc(target, double interface{}) *Patches {
	t := reflect.ValueOf(target)
	d := reflect.ValueOf(double)
	return p.ApplyCore(t, d)
}

func (p *Patches) ApplyCore(target, double reflect.Value) *Patches {
	check(target, double)
	assTarget := *(*uintptr)(getPointer(target))
	if _, ok := p.originals[assTarget]; ok {
		panic("patch has been existed")
	}

	p.valueHolders[double] = double
	original := replace(assTarget, uintptr(getPointer(double)))
	p.originals[assTarget] = original
	return p

}

func (p *Patches) Reset() {
	for target, bytes := range p.originals {
		modifyBinary(target, bytes)
		delete(p.originals, target)
	}

	for target, variable := range p.values {
		target.Elem().Set(variable)
	}
}

func replace(target, double uintptr) []byte {
	code := buildJmpDirective(double)
	bytes := entryAddress(target, len(code))
	original := make([]byte, len(bytes))
	copy(original, bytes)
	modifyBinary(target, code)
	return original
}

func buildJmpDirective(double uintptr) []byte {
	d0 := byte(double)
	d1 := byte(double >> 8)
	d2 := byte(double >> 16)
	d3 := byte(double >> 24)
	d4 := byte(double >> 32)
	d5 := byte(double >> 40)
	d6 := byte(double >> 48)
	d7 := byte(double >> 56)

	return []byte{
		0x48, 0xBA, d0, d1, d2, d3, d4, d5, d6, d7, // MOV rdx, double
		0xFF, 0x22, // JMP [rdx]
	}
}

// 关键函数：重写目标函数
func modifyBinary(target uintptr, bytes []byte) {
	function := entryAddress(target, len(bytes))

	page := entryAddress(pageStart(target), syscall.Getpagesize())
	var err error
	if runtime.GOOS == "darwin" {
		err = syscall.Mprotect(page, syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC)

	} else {
		err = syscall.Mprotect(page, syscall.PROT_READ|syscall.PROT_WRITE)
	}
	if err != nil {
		panic(err)
	}
	copy(function, bytes)

	err = syscall.Mprotect(page, syscall.PROT_READ|syscall.PROT_EXEC)
	if err != nil {
		panic(err)
	}
}

func entryAddress(p uintptr, l int) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: p, Len: l, Cap: l}))
}

func pageStart(ptr uintptr) uintptr {
	return ptr & ^(uintptr(syscall.Getpagesize() - 1))
}

func getPointer(v reflect.Value) unsafe.Pointer {
	return (*funcValue)(unsafe.Pointer(&v)).p
}

type funcValue struct {
	_ uintptr
	p unsafe.Pointer
}

func check(target, double reflect.Value) {
	if target.Kind() != reflect.Func {
		panic("target is not a func")
	}

	if double.Kind() != reflect.Func {
		panic("double is not a func")
	}

	if target.Type() != double.Type() {
		panic(fmt.Sprintf("target type(%s) and double type(%s) are different", target.Type(), double.Type()))
	}
}
