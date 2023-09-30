package middlewares

import (
	"net/http"
	"strings"
)

var DefaultLoader = func(name string, args []interface{}) func(http.Handler) http.Handler {
	switch name {
	case "setHeader":
		headerName := args[0].(string)
		headerValue := args[1].(string)
		target := "request"
		if len(args) > 2 {
			target = args[2].(string)
		}
		return MakeSetHeaderMiddleware(headerName, headerValue, target)
	case "removeHeader":
		headerName := args[0].(string)
		target := "request"
		if len(args) > 2 {
			target = args[2].(string)
		}
		return MakeRemoveHeaderMiddleware(headerName, target)
	case "moveHeader":
		source := args[0].(string)
		destination := args[1].(string)
		target := "request"
		if len(args) > 2 {
			target = args[2].(string)
		}
		return MakeMoveHeaderMiddleware(source, destination, target)
	case "setPath":
		return MakeSetPathMiddleware(args[0].(string))
	case "setPathPrefix":
		return MakeSetRequestPathPrefixMiddleware(args[0].(string))
	case "setRequestParameter":
		return MakeSetRequestParameterMiddleware(args[0].(string), args[1].(string))
	case "removeRequestParameter":
		return MakeRemoveRequestParameterMiddleware(args[0].(string))
	case "dedupeResponseHeaders":
		headers := []string{}
		for _, header := range args {
			headers = append(headers, header.(string))
		}
		return MakeDedupeResponseHeaders(headers...)
	case "redirect":
		return MakeRedirectMiddleware(args[0].(int), args[1].(string))
	case "validateJWTByJWKSURL":
		whitelist := []string{}
		if len(args) > 1 {
			whitelist = strings.Split(args[1].(string), ",")
		}
		return MakeCheckJWTValidityByJWKSURL(args[0].(string), whitelist, nil, nil)
	case "copyJWTClaimToHeader":
		whitelist := []string{}
		if len(args) > 2 {
			whitelist = strings.Split(args[2].(string), ",")
		}
		return MakeJWTCopyClaimToHeaderMiddleware(args[0].(string), args[1].(string), whitelist)
	case "setResponseBody":
		return MakeSetResponseBodyMiddleware(strings.NewReader(args[0].(string)))
	case "setResponseStatusCode":

		return MakeSetResponseStatusCodeMiddleware(int(args[0].(float64)))
	default:
		return nil
	}
}
