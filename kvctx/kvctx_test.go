package kvctx

import (
	"context"
	"testing"
)

func Test_valueCtx(t *testing.T) {

	tests := []struct {
		name  string
		key   string
		value float64
		want  float64
	}{
		{
			"pass-float",
			"str-key",
			10.0,
			10.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := WithContext(context.Background())
			Add(ctx, tt.key, tt.value)
			v := Get(ctx, tt.key)
			value, ok := v.(float64)
			if !ok {
				t.Errorf("got value from context, result: %v, want %v", value, tt.want)
			}
		})
	}
}
