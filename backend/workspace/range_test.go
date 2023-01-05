package workspace

import "testing"

func TestRange_Match(t *testing.T) {
	tests := []struct {
		name string
		r    PortList
		port int
		want bool
	}{
		{
			name: "empty",
			r:    PortList{},
			port: 123,
			want: true,
		},
		{
			name: "single match",
			r:    PortList{80},
			port: 80,
			want: true,
		},
		{
			name: "single no match",
			r:    PortList{80},
			port: 443,
			want: false,
		},
		{
			name: "list match",
			r:    PortList{80, 443},
			port: 80,
			want: true,
		},
		{
			name: "list match upper edge",
			r:    PortList{80, 443},
			port: 443,
			want: true,
		},
		{
			name: "list no match (below)",
			r:    PortList{80, 443},
			port: 1,
			want: false,
		},
		{
			name: "list no match (above)",
			r:    PortList{80, 443},
			port: 444,
			want: false,
		},
		{
			name: "list no match (between)",
			r:    PortList{80, 443},
			port: 110,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Match(tt.port); got != tt.want {
				t.Errorf("PortList.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
