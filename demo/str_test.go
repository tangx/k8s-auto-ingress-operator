package demo

import (
	"fmt"
	"strings"
	"testing"
)

func Test_StrHasPrefix(t *testing.T) {
	fmt.Println(has2("srv-123"))
	fmt.Println(has2("web-23"))
	fmt.Println(has2("istio-mamanger"))

}

func has(name string) bool {
	if strings.HasPrefix(name, "web-") || strings.HasPrefix(name, "srv-") {
		return true
	}

	return false
}

func has2(name string) bool {
	if !strings.HasPrefix(name, "web-") && !strings.HasPrefix(name, "srv-") {
		return false
	}

	return true
}
