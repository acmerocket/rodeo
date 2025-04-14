package main

import "testing"

func Test_matches(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		type_name string
		types     []string
		want      bool
	}{
		{name: "fail", type_name: "app.bsky.feed.post", types: []string{"yoyo"}, want: false},
		{name: "full", type_name: "app.bsky.feed.post", types: []string{"app.bsky.feed.post"}, want: true},
		{name: "partial", type_name: "app.bsky.feed.post", types: []string{"post"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matches(tt.type_name, tt.types)
			if got != tt.want {
				t.Errorf("matches() = %v, want %v", got, tt.want)
			}
		})
	}
}
