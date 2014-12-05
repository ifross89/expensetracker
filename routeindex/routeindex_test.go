package routeindex

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var index = CreateMemoryIndex(
	[]RouteInfo{
		RouteInfo{
			RouteString: "/testroute1",
			Name:        "test1",
		},
		RouteInfo{
			RouteString: "/testroute2/:category/:id",
			Name:        "test2",
		},
		RouteInfo{
			RouteString: "/testroute3/:a",
			Name:        "test3",
		},
	}...)

func TestRetrieveRoute(t *testing.T) {
	Convey("Given an index with some starting values", t, func() {

		Convey("When attempting to retrieve a route with the same name as one added", func() {
			r := index.MustByName("test1")
			Convey("The correct route should be returned", func() {
				So(r.RouteString, ShouldEqual, "/testroute1")
			})
		})

		Convey("When attempting to retrieve a route not present in the index", func() {
			_, err := index.ByName("NotPresent")
			Convey("An appropriate error should be returned", func() {
				So(err, ShouldNotEqual, nil)
			})
		})

		Convey("When attempting to retrieve a route not present in the index with MustByName", func() {
			Convey("The index should panic", func() {
				So(func() { index.MustByName("NotPresent") }, ShouldPanic)
			})
		})
	})
}

func TestURLCreationFailures(t *testing.T) {
	Convey("Given an index initialized with various routes", t, func() {
		Convey("The index should return an error when an invalid route name is supplied", func() {
			_, err := index.URL("NotPresent")
			So(err, ShouldNotEqual, nil)
		})

		Convey("The index should return an error when the arguments supplied are not in pairs", func() {
			_, err := index.URL("test1", "oddParam")
			So(err, ShouldNotEqual, nil)
		})

		Convey("The index should return an error with invalid parameters", func() {
			_, err := index.URL("test2", "not_present", "wont_be_inserted")
			So(err, ShouldNotEqual, nil)
		})
	})
}

func TestURLCreation(t *testing.T) {

	tests := []struct {
		route       string
		params      []string
		shouldError bool
		expected    string
	}{
		{route: "test1", params: []string{}, shouldError: false, expected: "/testroute1"},
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
