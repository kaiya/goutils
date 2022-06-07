package hotpatch

import "testing"

func Test_Hotpatch(t *testing.T) {
	patch := ApplyFunc(fakeFunc, func() string {
		return "should invoke this"
	})
	defer patch.Reset()
	output := fakeFunc()
	t.Logf("got output: %s", output)
}

func fakeFunc() string {
	return "fakeFunc"
}
