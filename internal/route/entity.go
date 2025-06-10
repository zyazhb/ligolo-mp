package route

import (
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/ttpreport/ligolo-mp/v2/cmd/client/tui/utils"
	pb "github.com/ttpreport/ligolo-mp/v2/protobuf"
)

type Route struct {
	ID         uuid.UUID
	Cidr       *net.IPNet
	IsLoopback bool
}

func NewRoute(cidr string, isLoopback bool) (*Route, error) {
	_, dst, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	return &Route{
		ID:         uuid.New(),
		Cidr:       dst,
		IsLoopback: isLoopback,
	}, nil
}

func (route *Route) String() string {
	return fmt.Sprintf("ID=%s, CIDR=%s IsLoopback=%s", route.ID.String(), route.Cidr.String(), utils.HumanBool(route.IsLoopback))
}

func (route *Route) Proto() *pb.Route {
	return &pb.Route{
		ID:         route.ID.String(),
		Cidr:       route.Cidr.String(),
		IsLoopback: route.IsLoopback,
	}
}

func ProtoToRoute(p *pb.Route) *Route {
	_, cidr, _ := net.ParseCIDR(p.Cidr)

	return &Route{
		ID:         uuid.MustParse(p.ID),
		Cidr:       cidr,
		IsLoopback: p.IsLoopback,
	}
}
