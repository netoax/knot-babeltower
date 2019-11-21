package network

import (
	"fmt"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
)

// ThingProxy proxy a request to the thing service interface
type ThingProxy interface {
	SendCreateThing(id, name, authorization string) (idGenerated string, err error)
}

// Proxy proxy a request to the thing service
type Proxy struct {
	url    string
	logger logging.Logger
}

// NewThingProxy creates a proxy to the thing service
func NewThingProxy(logger logging.Logger, hostname string, port uint16) Proxy {
	url := fmt.Sprintf("http://%s:%d", hostname, port)

	logger.Debug("Proxy setup to " + url)
	return Proxy{url, logger}
}

// SendCreateThing proxy the http request to thing service
func (p Proxy) SendCreateThing(id, name, authorization string) (idGenerated string, err error) {
	p.logger.Debug("Proxying request to create thing")
	return idGenerated, nil
}
