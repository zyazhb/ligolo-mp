package tunlink

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

func New() (int, string, error) {
	return 0, "", fmt.Errorf("not implmeneted for this architecture")
}

func Remove(ID int) error {
	return fmt.Errorf("not implemented for this architecture")
}

func AddRoute(ID int, CIDR *net.IPNet, metric int) error {
	return fmt.Errorf("not implemented for this architecture")
}

func RemoveAllRoutes(ID int) error {
	return fmt.Errorf("not implemented for this architecture")
}

func GetName(ID int) (string, error) {
	return "", fmt.Errorf("not implemented for this architecture")
}

func GetRoute(address net.IP) ([]netlink.Route, error) {
	return nil, fmt.Errorf("not implemented for this architecture")
}
