//go:build linux

package tunlink

import (
	"net"

	"github.com/vishvananda/netlink"
)

func New() (int, string, error) {
	la := netlink.NewLinkAttrs()
	link := &netlink.Tuntap{
		LinkAttrs: la,
		Mode:      netlink.TUNTAP_MODE_TUN,
	}
	if err := netlink.LinkAdd(link); err != nil {
		return 0, "", err
	}

	if err := netlink.LinkSetUp(link); err != nil {
		return 0, "", err
	}

	return link.Index, link.Name, nil
}

func Remove(ID int) error {
	link, err := netlink.LinkByIndex(ID)
	if err != nil {
		return err
	}

	if err := netlink.LinkDel(link); err != nil {
		return err
	}

	return nil
}

func AddRoute(ID int, CIDR *net.IPNet, metric int) error {
	route := &netlink.Route{
		LinkIndex: ID,
		Dst:       CIDR,
		Priority:  metric,
	}
	if err := netlink.RouteAdd(route); err != nil {
		return err
	}

	return nil
}

func RemoveAllRoutes(ID int) error {
	link, err := netlink.LinkByIndex(ID)
	if err != nil {
		return err
	}

	currentRoutes, err := netlink.RouteList(link, netlink.FAMILY_ALL)
	if err != nil {
		return err
	}

	for _, route := range currentRoutes {
		netlink.RouteDel(&route)
	}

	return nil
}

func GetName(ID int) (string, error) {
	link, err := netlink.LinkByIndex(ID)
	if err != nil {
		return "", err
	}
	return link.Attrs().Name, nil
}

func GetRoute(address net.IP) ([]netlink.Route, error) {
	routes, err := netlink.RouteGet(address)
	if err != nil {
		return nil, err
	}

	return routes, nil
}
