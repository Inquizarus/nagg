package events

import (
	"net/http"

	"github.com/inquizarus/gosebus"
	gevent "github.com/inquizarus/gosebus/pkg/event"
	"github.com/inquizarus/nagg/internal/domain"
)

const (
	REQUEST_WAS_HANDLED_EVENT   = "nagg_request_handled"
	UPSTREAM_REQUEST_DONE_EVENT = "nagg_upstream_request_done"
	ROUTE_MATCHED_REQUEST_EVENT = "nagg_route_matched_request"
)

type RequestHandled struct {
	Request        http.Request
	ResponseWriter http.ResponseWriter
	Route          domain.Route
}

type UpstreamRequestDone struct {
	Request  http.Request
	Response http.Response
}

type RouteMatchedRequest struct {
	Request http.Request
	Route   domain.Route
}

func PublishRequestWasHandled(bus gosebus.Bus, responseWriter http.ResponseWriter, request http.Request, route domain.Route) error {
	return bus.Publish(gevent.NewEvent(REQUEST_WAS_HANDLED_EVENT, RequestHandled{
		Request:        request,
		ResponseWriter: responseWriter,
		Route:          route,
	}))
}

func PublishUpstreamRequestWasDone(bus gosebus.Bus, response http.Response, request http.Request) error {
	return bus.Publish(gevent.NewEvent(UPSTREAM_REQUEST_DONE_EVENT, UpstreamRequestDone{
		Request:  request,
		Response: response,
	}))
}

func PublishRouteMatchedRequest(bus gosebus.Bus, route domain.Route, request http.Request) error {
	return bus.Publish(gevent.NewEvent(ROUTE_MATCHED_REQUEST_EVENT, RouteMatchedRequest{
		Request: request,
		Route:   route,
	}))
}
