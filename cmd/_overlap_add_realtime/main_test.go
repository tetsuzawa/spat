package main

import (
	"reflect"
	"testing"
)

func Test_shift(t *testing.T) {
	type args struct {
		x []float32
		n int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				x: []float32{1., 2., 3., 4., 5., 6., 7., 8., 9., 10.},
				n: 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shift(tt.args.x, tt.args.n)
			want := []float32{4., 5., 6., 7., 8., 9., 10., 0., 0., 0.}
			if !reflect.DeepEqual(tt.args.x, want) {
				t.Errorf("failed got:%v, want:%v", tt.args.x, want)
			}
		})
	}
}
