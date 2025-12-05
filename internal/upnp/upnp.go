package upnp

// obv fix this
import (
	"errors"
	"net"

	"github.com/huin/goupnp/dcps/ocf/internetgateway2"
	"golang.org/x/sync/errgroup"
)

type RouterClient interface {
	AddPortMapping(
		NewRemoteHost string,
		NewExternalPort uint16,
		NewProtocol string,
		NewInternalPort uint16,
		NewInternalClient string,
		NewEnabled bool,
		NewPortMappingDescription string,
		NewLeaseDuration uint32,
	) (err error)

	DeletePortMapping(NewRemoteHost string, NewExternalPort uint16, NewProtocol string) (err error)

	GetExternalIPAddress() (
		NewExternalIPAddress string,
		err error,
	)

	LocalAddr() net.IP
}

func NewClient() (RouterClient, error) {
	g := errgroup.Group{}

	var ip1Clients []*internetgateway2.WANIPConnection1
	g.Go(func() error {
		var err error
		ip1Clients, _, err = internetgateway2.NewWANIPConnection1Clients()
		return err
	})

	var ip2Clients []*internetgateway2.WANIPConnection2
	g.Go(func() error {
		var err error
		ip2Clients, _, err = internetgateway2.NewWANIPConnection2Clients()
		return err
	})
	var ppp1Clients []*internetgateway2.WANPPPConnection1
	g.Go(func() error {
		var err error
		ppp1Clients, _, err = internetgateway2.NewWANPPPConnection1Clients()
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// Handle multiple client return in the future
	switch {
	case len(ip2Clients) == 1:
		return ip2Clients[0], nil
	case len(ip1Clients) == 1:
		return ip1Clients[0], nil
	case len(ppp1Clients) == 1:
		return ppp1Clients[0], nil
	default:
		return nil, errors.New("multiple or no services found")
	}
}
