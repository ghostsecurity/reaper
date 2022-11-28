package highlight

import "testing"

func TestHighlighting(t *testing.T) {
	tests := []struct {
		name  string
		theme string
		input string
		want  string
	}{
		{
			name:  "empty",
			theme: "",
			input: "",
			want:  "",
		},
		{
			name:  "headers with dark theme",
			theme: "dark",
			input: `GET / HTTP/1.1
Accept: */*
`,
			want: `<span style="display:flex;"><span><span style="color:#2299cf">GET</span> <span style="color:#2299cf">/</span> <span style="color:#fff;font-weight:bold">HTTP</span>/<span style="color:#2de;font-weight:bold">1.1</span>
</span></span><span style="display:flex;"><span><span style="color:#2299cf">Accept</span>: */*
</span></span><br/><br/>`,
		},
		{
			name:  "headers with ghost theme",
			theme: "ghost",
			input: `GET / HTTP/1.1
Accept: */*
`,
			want: `<span style="display:flex;"><span><span style="color:#6a84fa;font-weight:bold">GET</span> <span style="color:#c3c8df">/</span> <span style="color:#6a84fa">HTTP</span><span style="color:#c3c8df">/</span>1.1
</span></span><span style="display:flex;"><span><span style="color:#6a84fa">Accept</span><span style="color:#c3c8df">:</span> */*
</span></span><br/><br/>`,
		},
		{
			name:  "headers with light theme",
			theme: "light",
			input: `GET / HTTP/1.1
Accept: */*
`,
			want: `<span style="display:flex;"><span><span style="color:#58a1dd">GET</span> <span style="color:#58a1dd">/</span> <span style="color:#ff636f">HTTP</span><span style="color:#ff636f">/</span><span style="color:#a6be9d">1.1</span>
</span></span><span style="display:flex;"><span><span style="color:#58a1dd">Accept</span><span style="color:#ff636f">:</span> <span style="color:#a6be9d">*/*</span>
</span></span><br/><br/>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HTTP(tt.input, tt.theme)
			if got != tt.want {
				t.Errorf("Highlight() = %v, want %v", got, tt.want)
			}
		})
	}
}
