package tun

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/hashicorp/yamux"
	"github.com/ttpreport/ligolo-mp/internal/netstack"
	"github.com/ttpreport/ligolo-mp/internal/netstack/tunlink"
	"github.com/ttpreport/ligolo-mp/internal/route"
	"github.com/ttpreport/ligolo-mp/pkg/memstore"
	pb "github.com/ttpreport/ligolo-mp/protobuf"
)

type Tun struct {
	ID       int
	Name     string
	Active   bool
	Routes   *memstore.Syncmap[string, *route.Route]
	netstack *netstack.NetStack `json:"-"`
}

func NewTun() (*Tun, error) {
	slog.Debug("creating new tun")
	ret := &Tun{
		Routes: memstore.NewSyncmap[string, *route.Route](),
	}

	return ret, nil
}

func (t *Tun) Start(multiplex *yamux.Session, maxConnections int, maxInFlight int) error {
	if t.Active {
		return nil
	}

	linkID, linkName, err := tunlink.New()
	if err != nil {
		slog.Error("could not create tun link")
		return err
	}
	slog.Debug("tun link created", slog.Any("id", linkID), slog.Any("name", linkName))

	t.ID = linkID
	t.Name = linkName

	ns, err := netstack.NewNetstack(maxConnections, maxInFlight, t.Name)
	if err != nil {
		slog.Error("could not create netstack")
		return err
	}
	slog.Debug("netstack created", slog.Any("netstack", ns))

	t.netstack = ns

	go func() {
		for {
			select {
			case <-t.netstack.ClosePool(): // pool closed, we can't process packets!
				slog.Debug("connection pool closed")
				return
			case relayPacket := <-t.netstack.GetTunConn(): // Process connections/packets
				localRoutes := t.GetLocalRoutes() // a bit dirty, but allows granular localhost routing
				go t.netstack.HandlePacket(relayPacket, multiplex, localRoutes)
			}
		}
	}()

	t.Active = true

	slog.Debug("tun activated")

	if err := t.ApplyRoutes(); err != nil {
		slog.Error("could not apply routes")
	}

	return nil
}

func (t *Tun) Stop() {
	if err := tunlink.Remove(t.ID); err != nil {
		slog.Debug("could not delete link", slog.Any("error", err))
	}
	slog.Debug("tun removed", slog.Any("tun", t))

	if t.netstack != nil {
		if err := t.netstack.Destroy(); err != nil {
			slog.Debug("could not destroy netstack", slog.Any("err", err), slog.Any("netstack", t.netstack))
		}
	}

	t.Active = false
}

func (t *Tun) ApplyRoutes() error {
	if t.Active {
		slog.Debug("applying routes")

		if err := t.removeAllRoutes(); err != nil {
			slog.Warn("could not remove current routes", slog.Any("err", err))
		}

		for _, route := range t.Routes.All() {
			err := tunlink.AddRoute(t.ID, route.Cidr)
			if err != nil {
				slog.Error("could not add route to the system", slog.Any("err", err), slog.Any("route", route))
			}
		}
	}

	return nil
}

func (t *Tun) removeAllRoutes() error {
	err := tunlink.RemoveAllRoutes(t.ID)
	if err != nil {
		slog.Error("could not find link", slog.Any("link_id", t.ID))
		return err
	}

	return nil
}

func (t *Tun) NewRoute(cidr string, isLoopback bool) error {
	slog.Debug("adding route to tun", slog.Any("route", cidr))

	_, dst, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}

	t.Routes.Set(dst.String(), &route.Route{
		Cidr:       dst,
		IsLoopback: isLoopback,
	})

	slog.Debug("route added to tun")

	return nil
}

func (t *Tun) RemoveRoute(cidr string) error {
	slog.Debug("removing route from tun", slog.Any("route", cidr))

	_, dst, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}

	t.Routes.Delete(dst.String())

	if err := t.ApplyRoutes(); err != nil {
		slog.Error("could not apply routes", slog.Any("tun", t), slog.Any("routes", t.Routes.All()))
		return err
	}

	slog.Debug("route removed", slog.Any("routes", cidr))

	return nil
}

func (t *Tun) GetName() (string, error) {
	return tunlink.GetName(t.ID)
}

func (t *Tun) GetRoutes() []route.Route {
	var result []route.Route
	for _, route := range t.Routes.All() {
		result = append(result, *route)
	}

	return result
}

func (t *Tun) GetLocalRoutes() []route.Route {
	var routes []route.Route
	for _, route := range t.Routes.All() {
		if route.IsLoopback {
			routes = append(routes, *route)
		}
	}

	return routes
}

func (t *Tun) String() string {
	return fmt.Sprintf("ID=%d Name=%s with %d routes", t.ID, t.Name, len(t.Routes.All()))
}

func (t *Tun) Proto() *pb.Tun {
	var Routes []*pb.Route
	for _, route := range t.Routes.All() {
		Routes = append(Routes, route.Proto())
	}

	return &pb.Tun{
		Name:   t.Name,
		Routes: Routes,
	}
}

func ProtoToTun(p *pb.Tun) *Tun {
	routes := memstore.NewSyncmap[string, *route.Route]()
	for _, r := range p.Routes {
		routes.Set(r.Cidr, route.ProtoToRoute(r))
	}

	return &Tun{
		Name:   p.Name,
		Routes: routes,
	}
}
