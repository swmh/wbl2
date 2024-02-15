package main

import (
	"reflect"
	"testing"
)

func TestGroupAnagrams(t *testing.T) {
	type args struct {
		words []string
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "#1",
			args: args{
				words: []string{},
			},
			want: map[string][]string{},
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GroupAnagrams(tt.args.words); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupAnagrams() = %v, want %v", got, tt.want)
			}
		})
	}
}
