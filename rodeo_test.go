package main

import (
	"strings"
	"testing"
)

func params(param string) map[string]string {
	return parse_params(strings.Split(param, " "))
}

func Test_matches(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		type_name string
		params    string
		want      bool
	}{
		{name: "no match", type_name: "app.bsky.feed.post", params: "yoyo", want: false},
		{name: "full", type_name: "app.bsky.feed.post", params: "app.bsky.feed.post", want: true},
		{name: "partial", type_name: "app.bsky.feed.post", params: "post", want: true},
		{name: "list", type_name: "app.bsky.feed.post", params: "like post", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matches(tt.type_name, params(tt.params))
			if got != tt.want {
				t.Errorf("matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resolve_template(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		type_name string
		params    string
		expected  string
	}{
		{name: "test1", type_name: "app.bsky.feed.post", params: "like post", expected: "app.bsky.feed.post"},
		{name: "test2", type_name: "app.bsky.feed.post", params: "", expected: "app.bsky.feed.post"},
		{name: "test3", type_name: "app.bsky.feed.post", params: "like", expected: ""},
		{name: "test4", type_name: "app.bsky.feed.post", params: "post=TEST", expected: "TEST"},
		{name: "test5", type_name: "app.bsky.feed.post", params: "feed.post=test.md", expected: "test.md"},
		{name: "test6", type_name: "app.bsky.feed.post", params: "app.bsky.feed.like=default", expected: ""},
		{name: "test7", type_name: "app.bsky.feed.like", params: "", expected: "app.bsky.feed.like"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolve_template(tt.type_name, params(tt.params))
			if got != tt.expected {
				t.Errorf("resolve_template() = %v, want %v", got, tt.expected)
			}
		})
	}
}
