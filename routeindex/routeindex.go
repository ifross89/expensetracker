package routeindex

import (
	"net/http"
	"strings"
	"fmt"
)

type RouteInfo struct {
	TemplateNames []string
	RouteString   string
	Name          string
	http.Handler
}

type Interface interface {
	ByName(string) (RouteInfo, error)
	MustByName(string) RouteInfo
	URL(string, ...string) (string, error)
}

type memoryIndex map[string]RouteInfo

func CreateMemoryIndex(routes ...RouteInfo) Interface {
	index := make(memoryIndex)
	for _, r := range routes {
		index[r.Name] = r
	}

	return index
}

func (i memoryIndex) ByName(name string) (RouteInfo, error) {
	r, ok :=  i[name]
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

	if len(pairs) % 2 != 0 {
		return "", fmt.Errorf("number of pairs must be even, got %d", len(pairs))
	}

	out := ri.RouteString

	for i := 0; i < len(pairs); i += 2 {
		temp := strings.Replace(out, ":" + pairs[i], pairs[i+1], 1)
		if temp == out {
			return "", fmt.Errorf("Unable to find paramerter %s in route string (%s)", pairs[i], ri.RouteString)
		}

		out = temp
	}

	return out, nil
}
