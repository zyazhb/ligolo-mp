package tun

import (
	"fmt"

	"gvisor.dev/gvisor/pkg/tcpip/stack"
)

func New(tunName string) (stack.LinkEndpoint, int, error) {
	return nil, 0, fmt.Errorf("not implmeneted for this architecture")
}
