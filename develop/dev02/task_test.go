package main

import "testing"

func TestUnpack(t *testing.T) {
	type args struct {
		s string
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
				s: "a4bc2d5e",
			},
			want:    "aaaabccddddde",
			wantErr: false,
		},
		{
			name: "#2",
			args: args{
				s: "abcd",
			},
			want:    "abcd",
			wantErr: false,
		},
		{
			name: "#3",
			args: args{
				s: "45",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "#4",
			args: args{
				s: "",
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "#5",
			args: args{
				s: `qwe\4\5`,
			},
			want:    "qwe45",
			wantErr: false,
		},
		{
			name: "#6",
			args: args{
				s: `qwe\45`,
			},
			want:    "qwe44444",
			wantErr: false,
		},
		{
			name: "#7",
			args: args{
				s: `qwe\\5`,
			},
			want:    `qwe\\\\\`,
			wantErr: false,
		},
		{
			name: "#8",
			args: args{
				s: `\`,
			},
			want:    ``,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unpack(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unpack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Unpack() = %v, want %v", got, tt.want)
			}
		})
	}
}
