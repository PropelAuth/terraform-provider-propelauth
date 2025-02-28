package propelauth

import (
	"testing"
)

func TestContains(t *testing.T) {
	type args struct {
		slice  []string
		target string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test contains",
			args: args{
				slice:  []string{"a", "b", "c"},
				target: "b",
			},
			want: true,
		},
		{
			name: "Test does not contain",
			args: args{
				slice:  []string{"a", "b", "c"},
				target: "d",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.slice, tt.args.target); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPortFromLocalhost(t *testing.T) {
	type args struct {
		inputUrl string
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 int
	}{
		{
			name: "Test localhost URL",
			args: args{
				inputUrl: "http://localhost:3001",
			},
			want:  true,
			want1: 3001,
		},
		{
			name: "Test invalid URL",
			args: args{
				inputUrl: "example.com",
			},
			want:  false,
			want1: 0,
		},
		{
			name: "Test URL without port",
			args: args{
				inputUrl: "https://example.com",
			},
			want:  false,
			want1: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetPortFromLocalhost(tt.args.inputUrl)
			if got != tt.want {
				t.Errorf("GetPortFromLocalhost() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetPortFromLocalhost() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
