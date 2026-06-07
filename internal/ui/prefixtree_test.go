package ui

import (
	"reflect"
	"testing"
)

func TestBuildPrefixTree(t *testing.T) {
	keys := []string{
		"users/0001",
		"users/0002",
		"logs/2026-06-03/0001",
		"meta/created_at",
	}
	got := buildPrefixTree(keys)

	want := map[string][]string{
		"":                 {"logs/", "meta/", "users/"},
		"users/":           {"users/0001", "users/0002"},
		"logs/":            {"logs/2026-06-03/"},
		"logs/2026-06-03/": {"logs/2026-06-03/0001"},
		"meta/":            {"meta/created_at"},
	}
	for parent, expect := range want {
		if !reflect.DeepEqual(got[parent], expect) {
			t.Errorf("children of %q = %v, want %v", parent, got[parent], expect)
		}
	}
}

func TestTreeLabel(t *testing.T) {
	cases := []struct{ id, want string }{
		{"", ""},
		{"users/", "users/"},
		{"users/0001", "0001"},
		{"logs/2026-06-03/", "2026-06-03/"},
		{"logs/2026-06-03/0001", "0001"},
	}
	for _, c := range cases {
		if got := treeLabel(c.id); got != c.want {
			t.Errorf("treeLabel(%q) = %q, want %q", c.id, got, c.want)
		}
	}
}
