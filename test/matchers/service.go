package matchers

import (
	"fmt"
	"github.com/sortednet/statuschecker/internal/store"
)

type ServiceMatcher struct {
	Name string
	Url  string
}

func (m *ServiceMatcher) Matches(x interface{}) bool {
	service, ok := x.(*store.Service)
	if !ok {
		fmt.Println("Could not convert type to Service")
		return false
	}

	if service.Name != m.Name {
		fmt.Printf("%s != %s \n", service.Name, m.Name)
		return false
	}

	if service.Url != m.Url {
		fmt.Printf("%s != %s \n", service.Url, m.Url)
		return false
	}

	return false
}
