package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/CESARBR/knot-babeltower/pkg/entities"
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

type errorConflict struct{ error }
type errorForbidden struct{ error }

func (err errorForbidden) Error() string {
	return "Error forbidden"
}

func (err errorConflict) Error() string {
	return "Error conflict"
}

// NewThingProxy creates a proxy to the thing service
func NewThingProxy(logger logging.Logger, hostname string, port uint16) Proxy {
	url := fmt.Sprintf("http://%s:%d", hostname, port)

	logger.Debug("Proxy setup to " + url)
	return Proxy{url, logger}
}

// SendCreateThing proxy the http request to thing service
func (p Proxy) SendCreateThing(id, name, authorization string) (idGenerated string, err error) {
	var resp *http.Response
	var req *http.Request
	var uuid, locationHeader string

	p.logger.Debug("Proxying request to create thing")
	client := &http.Client{Timeout: 10 * time.Second}
	jsonThing, err := json.Marshal(entities.Thing{ID: id, Name: name})
	if err != nil {
		p.logger.Error(err)
		goto done
	}

	p.logger.Debug(string(jsonThing))

	req, err = http.NewRequest("POST", p.url+"/things", bytes.NewBuffer(jsonThing))
	if err != nil {
		p.logger.Error(err)
		goto done
	}

	p.logger.Debug("Authorization header:", authorization)
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		p.logger.Error(err)
		goto done
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		p.logger.Errorf("Status not created: %d", resp.StatusCode)
		switch resp.StatusCode {
		case http.StatusConflict:
			err = errorConflict{}
		case http.StatusForbidden:
			err = errorForbidden{}
		}
		goto done
	}

	locationHeader = resp.Header.Get("Location")
	uuid = locationHeader[len("/things/"):] // get substring after "/things/"
	idGenerated = uuid
done:
	return idGenerated, err
}
