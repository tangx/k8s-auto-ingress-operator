package controllers

import (
	"fmt"
	"sync"

	v1 "github.com/tangx/k8s-auto-ingress-operator/api/v1"
)

type AutoIngressContainer struct {
	mu  sync.Mutex
	set map[string]v1.AutoIngress
}

func NewAutoIngressContainer() *AutoIngressContainer {
	return &AutoIngressContainer{
		set: make(map[string]v1.AutoIngress),
	}
}

func (c *AutoIngressContainer) Add(ing v1.AutoIngress) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.key(ing)
	c.set[key] = ing

}

func (c *AutoIngressContainer) Remove(ing v1.AutoIngress) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.key(ing)
	delete(c.set, key)
}

func (c *AutoIngressContainer) List() []v1.AutoIngress {
	list := make([]v1.AutoIngress, 0)
	for _, v := range c.set {
		list = append(list, v)
	}

	return list
}

func (c *AutoIngressContainer) key(ing v1.AutoIngress) string {
	return fmt.Sprintf("%s-%s", ing.Namespace, ing.Name)
}

var autoIngressSet *AutoIngressContainer

func init() {
	autoIngressSet = NewAutoIngressContainer()
}
