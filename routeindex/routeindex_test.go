package routeindex

import (
	"testing"
)

func TestURLCreation(t *testing.T) {
	index := CreateMemoryIndex(
		[]RouteInfo{
			RouteInfo{
				RouteString:  "/testroute1",
				Name: "test1",
			},
			RouteInfo{
				RouteString: "/testroute2/:category/:id",
				Name: "test2",
			},
			RouteInfo{
				RouteString: "/testroute3/:a",
				Name: "test3",
			},
		}...)

	tests := []struct{
			route string
			params []string
			shouldError bool
			expected string}{
				{route:"test1",params: []string{},shouldError: false, expected: "/testroute1"},
				{route: "test2", params: []string{"category", "baseball", "id", "23"}, shouldError: false, expected: "/testroute2/baseball/23"},
				{route: "test3", params: []string{"a", "b"}, shouldError: false, expected: "/testroute3/b"},
			}

	for i, test := range tests {
		url, err := index.URL(test.route, test.params...)
		if test.shouldError && err == nil {
			t.Fatalf("Expected error for test#%d, got no error", i)
			return
		}

		if err != nil {
			t.Fatalf("Expected no error for test#%d but got err=%v", i, err)
			return
		}

		if test.expected != url {
			t.Fatalf("Fail for test#%d. Expected %s, got %s", i, test.expected, url)
			return
		}
	}
}
