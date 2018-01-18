package haslocation

import (
	"errors"

	"github.com/FreifunkBremen/yanic/output/filter"
	"github.com/FreifunkBremen/yanic/runtime"
)

type haslocation struct {
	has bool
}

func init() {
	filter.Register("has_location", build)
}

func build(v interface{}) (filter.Filter, error) {
	if config, ok := v.(bool); ok {
		return &haslocation{has: config}, nil
	}
	return nil, errors.New("invalid configuration, bool expected")
}

func (h *haslocation) Apply(node *runtime.Node) *runtime.Node {
	if nodeinfo := node.Nodeinfo; nodeinfo != nil {
		if h.has {
			if location := nodeinfo.Location; location != nil {
				return node
			}
		} else {
			if location := nodeinfo.Location; location == nil {
				return node
			}
		}
	} else if !h.has {
		return node
	}
	return nil
}
