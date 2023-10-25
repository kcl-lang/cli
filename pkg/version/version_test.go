package version

import "testing"

func TestGetVersionString(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "test get version in string",
			want: VersionTypeLatest.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetVersionString(); got != tt.want {
				t.Errorf(" GetVersionString() = %v, want %v", got, tt.want)
			}
		})
	}
}
