package schema

import (
	"github.com/open-webtech/go-xaddy-config/schema/nodes"
)

type Builder struct {
	nodes.NodesContainer
}

func NewBuilder() *Builder {
	return &Builder{}
}
