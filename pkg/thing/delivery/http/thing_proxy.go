package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
)

// ThingProxy proxy a request to the thing service interface
type ThingProxy interface {
	SendCreateThing(id, name, authorization string) (idGenerated string, err error)
	UpdateSchema(ID string, schemaList []entities.Schema) error
	Get(ID string) (thing *entities.Thing, err error)
}

// Proxy proxy a request to the thing service
type Proxy struct {
	url    string
	logger logging.Logger
}

// RequestInfo aims to group all request releated information
type RequestInfo struct {
	method        string
	url           string
	authorization string
	contentType   string
	data          []byte
}

type UpdateSchemaRequest struct {
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

func (p Proxy) mapErrorFromStausCode(code int) error {
	var err error

	if code != http.StatusCreated {
		switch code {
		case http.StatusConflict:
			err = errorConflict{}
		case http.StatusForbidden:
			err = errorForbidden{}
		}
	}
	return err
}

func (p Proxy) sendRequest(info *RequestInfo) (*http.Response, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(info.method, info.url, bytes.NewBuffer(info.data))
	if err != nil {
		p.logger.Error(err)
		return nil, err
	}

	req.Header.Set("Authorization", info.authorization)
	req.Header.Set("Content-Type", info.contentType)

	return client.Do(req)
}

// SendCreateThing proxy the http request to thing service
func (p Proxy) SendCreateThing(id, name, authorization string) (idGenerated string, err error) {
	p.logger.Debug("Proxying request to create thing")
	jsonThing, err := json.Marshal(entities.Thing{ID: id, Name: name})
	if err != nil {
		p.logger.Error(err)
		return "", err
	}

	requestInfo := &RequestInfo{
		"POST",
		p.url + "/things",
		authorization,
		"application/json",
		jsonThing,
	}

	resp, err := p.sendRequest(requestInfo)
	print(resp)
	if err != nil {
		p.logger.Error(err)
		return "", err
	}

	locationHeader := resp.Header.Get("Location")
	fmt.Print(locationHeader)
	thingID := locationHeader[len("/things/"):] // get substring after "/things/"
	return thingID, p.mapErrorFromStausCode(resp.StatusCode)
}

// UpdateSchema receives the thing's ID and schema and send a HTTP request to
// the thing's service in order to update it with the schema.
func (p Proxy) UpdateSchema(ID string, schemaList []entities.Schema) error {
	parsedSchema, err := json.Marshal(schemaList)
	if err != nil {
		p.logger.Error(err)
		return err
	}

	requestInfo := &RequestInfo{
		"PUT",
		p.url + "/things" + ID,
		"authorization",
		"application/json",
		parsedSchema,
	}

	resp, err := p.sendRequest(requestInfo)
	if err != nil {
		p.logger.Error(err)
		return err
	}

	defer resp.Body.Close()

	return p.mapErrorFromStausCode(resp.StatusCode)
}

// Get proxy a request to find a registered thing according to its ID
func (p Proxy) Get(ID string) (*entities.Thing, error) {
	return &entities.Thing{}, nil
}
