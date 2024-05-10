package main

import "testing"

func TestA(t *testing.T) {

	tests := map[string]struct {
		input  string
		result string
	}{
		"one": {
			input:  "https://www.domain.com/",
			result: "https://www.domain.com/cdn-cgi/image/onerror=redirect/http://anything.www.domain.com",
		},
	}

	for name, test := range tests {
		test := test
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, _ := createRedirectURL(test.input)
			if got != test.result {
				t.Fatalf("createRedirectURL(%q) returned %q; expected %q", test.input, got, test.result)
			}
		})
	}
}
