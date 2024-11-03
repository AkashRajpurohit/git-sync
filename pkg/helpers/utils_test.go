package helpers

import (
	"testing"
)

func TestIsIncludedInListFunction(t *testing.T) {
	tests := []struct {
		list  []string
		value string
		want  bool
	}{
		{[]string{"foo", "bar"}, "baz", false},
		{[]string{"foo", "bar"}, "foo", true},
		{[]string{"foo", "ba*"}, "bar", true}, // partial match should return true
		{[]string{}, "anything", false},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			got := IsIncludedInList(tt.list, tt.value)
			if got != tt.want {
				t.Errorf("IsIncludedInList(%v, %q) = %t; want %t", tt.list, tt.value, got, tt.want)
			}
		})
	}
}
