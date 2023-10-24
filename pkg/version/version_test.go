package version

import "testing"

func TestGetVersionInStr(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "test get version in string",
			want: "0.7.0-alpha.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetVersionInStr(); got != tt.want {
				t.Errorf("GetVersionInStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
