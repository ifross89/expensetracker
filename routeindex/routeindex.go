package routeindex

import (
	"fmt"
	"net/http"
	"strings"
)

// RouteInfo contains information about a given route. The idea of this is to
// be "router agnostic". The aim of the structure is to keep all of the
// information about a route in one place. This can make reverse-lookups of,
// for example, URLs possible.
type RouteInfo struct {
	TemplateNames []string
	RouteString   string
	Name          string
	http.Handler
}

// Interface is the interface that is satisfied by the index of routes.
type Interface interface {
	// ByName retrieves a route given the name of the route. If there
	// is no route with the name then an error is returned.
	ByName(string) (RouteInfo, error)

	// MustByName returns the RouteInfo structure for the given name.
	// If the route does not exist then the method panics. This should
	// only be used to initialize the server, not during general use
	// of the index.
	MustByName(string) RouteInfo

	// URL creates a URL for the named route. If the route string contains
	// route parameters, then this method takes pairs. The first is the name
	// of the parameter and the second is the value to replace it with.
	// If the parameter is not contained within the route string, an
	// error is returned.
	URL(string, ...string) (string, error)
}

type memoryIndex map[string]RouteInfo

// CreateMemoryIndex creates an in-memory route index that satisfies
// Interface.
func CreateMemoryIndex(routes ...RouteInfo) Interface {
	index := make(memoryIndex)
	for _, r := range routes {
		index[r.Name] = r
	}

	return index
}

func (i memoryIndex) ByName(name string) (RouteInfo, error) {
	r, ok := i[name]
	if !ok {
		return RouteInfo{}, fmt.Errorf("ByName: no route named %s", name)
	}

	return r, nil
}

func (i memoryIndex) MustByName(name string) RouteInfo {
	r, err := i.ByName(name)
	if err != nil {
		panic(err.Error())
	}

	return r
}

func (i memoryIndex) URL(name string, pairs ...string) (string, error) {
	ri, ok := i[name]
	if !ok {
		return "", fmt.Errorf("No route named %s", name)
	}

	if len(pairs)%2 != 0 {
		return "", fmt.Errorf("number of pairs must be even, got %d", len(pairs))
	}

	out := ri.RouteString

	for i := 0; i < len(pairs); i += 2 {
		temp := strings.Replace(out, ":"+pairs[i], pairs[i+1], 1)
		if temp == out {
			return "", fmt.Errorf("Unable to find paramerter %s in route string (%s)", pairs[i], ri.RouteString)
		}

		out = temp
	}

	return out, nil
}
