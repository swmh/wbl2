package main

import (
	"testing"
)

func TestGetNumSuffix(t *testing.T) {
	type args struct {
		x string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 int
	}{
		{"#1", args{"10K"}, "10", 1000},
		{"#2", args{"20KB"}, "20", 1024},
		{"#3", args{"3M"}, "3", 1_000_000},
		{"#4", args{"4MB"}, "4", 1_048_576},
		{"#5", args{"5G"}, "5", 1_000_000_000},
		{"#6", args{"6GB"}, "6", 1_073_741_824},
		{"#7", args{"100"}, "", 0},
		{"#7", args{""}, "", 0},
		{"#7", args{"K"}, "", 1000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetNumSuffix(tt.args.x)
			if got != tt.want {
				t.Errorf("GetNumSuffix() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetNumSuffix() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSortByMonth(t *testing.T) {
	type args struct {
		a string
		b string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"#1", args{"jan", "feb"}, -1},
		{"#2", args{"feb", "jan"}, 1},
		{"#3", args{"feb", "feb"}, 0},
		{"#3", args{"invalid", "feb"}, -1},
		{"#3", args{"feb", "invalid"}, 1},
		{"#3", args{"", "feb"}, -1},
		{"#3", args{"feb", ""}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SortByMonth(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("SortByMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}
