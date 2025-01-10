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

func TestIsValidUrlWithoutTrailingSlash(t *testing.T) {
	type args struct {
		inputUrl string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Test valid URL",
			args: args{
				inputUrl: "https://example.com",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Test invalid URL",
			args: args{
				inputUrl: "example.com",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Test URL with trailing slash",
			args: args{
				inputUrl: "https://example.com/",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Test localhost URL",
			args: args{
				inputUrl: "http://localhost:3001",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsValidUrlWithoutTrailingSlash(tt.args.inputUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValidUrlWithoutTrailingSlash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsValidUrlWithoutTrailingSlash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidUrl(t *testing.T) {
	type args struct {
		inputUrl string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Test valid URL",
			args: args{
				inputUrl: "https://example.com",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Test invalid URL",
			args: args{
				inputUrl: "example.com",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Test URL with trailing slash",
			args: args{
				inputUrl: "https://example.com/",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Test localhost URL",
			args: args{
				inputUrl: "http://localhost:3001",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsValidUrl(tt.args.inputUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValidUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsValidUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
