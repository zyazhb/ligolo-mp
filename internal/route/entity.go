package route

import (
	"net"

	pb "github.com/ttpreport/ligolo-mp/protobuf"
)

type Route struct {
	Cidr       *net.IPNet
	IsLoopback bool
}

func (route *Route) Proto() *pb.Route {
	return &pb.Route{
		Cidr:       route.Cidr.String(),
		IsLoopback: route.IsLoopback,
	}
}

func ProtoToRoute(p *pb.Route) *Route {
	_, cidr, _ := net.ParseCIDR(p.Cidr)

	return &Route{
		Cidr:       cidr,
		IsLoopback: p.IsLoopback,
	}
}
