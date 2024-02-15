package main

import (
	"net/url"
	"reflect"
	"testing"
)

const first = `<a href="/programs">Программы</a></li><li><a href="/games">`

func TestGetLinks(t *testing.T) {
	type args struct {
		r []byte
	}
	tests := []struct {
		name string
		args args
		want []string
	}{

		{
			name: "#1",
			args: args{
				r: []byte(first),
			},
			want: []string{"https://habr.com/programs", "https://habr.com/games"},
		},
	}

	l := "https://habr.com/"
	u, err := url.Parse(l)
	if err != nil {
		t.Error(err)
	}

	cfg := Config{
		Url:  u,
		Link: l,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cfg.GetLinks(tt.args.r)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPathFromURL(t *testing.T) {
	type args struct {
		u string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "#1",
			args: args{
				u: "https://habr.com/ru/companies/ruvds/articles/346640/index.html",
			},
			want:    "habr.com/ru/companies/ruvds/articles/346640/index.html",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPathFromURL(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPathFromURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetPathFromURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
