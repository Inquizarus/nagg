# NAGG (Not Another Go Gateway)

This is an experimental API Gateway that probably shouldn't have been created. But here it is.

## Usage

Either run the standalone binary or import and register NAGG in your own application.

## Configuration

### Environment variables

|Name|Description|Default|
|---|---|---|
|NAGG_PORT|Defines which port to start built in server on|8080|
|NAGG_LOG_LEVEL|Determines which log level should be used|info|
|NAGG_HTTP_BASE_PATH|Determines which base path that NAGG will listen for requests on|/|
|NAGG_CONFIG|Gateway configuration in JSON format (has precedence) |-|
|NAGG_CONFIG_FILE_PATH|Gateway configuration in JSON format|-|

See `examples/gateway.json` for an example. 

#### Expansion
When configuration is loaded from `NAGG_CONFIG`, standard Golang environment expansion will be applied to expand any variable statements in the string.

Given that there is an environment variable named `FOO` with the value `BAR` then the following string in `NAGG_CONFIG`:

```
'{"foo":"$FOO"}'
```

Would result in the following string to be parsed as a config.

```
'{"foo":"bar"}'
```

### Gateway
Global space for configuration in json structure and have two attributes that can be configured.

**Middlewares**, global middlewares to apply on each request and response.
**Routes**, routes to handle requests for in the gateway.

### Middleware
Middlewares has three attributes.

**Name**, which determines which middleware is loaded.
**Phase**, either pre or post which determines if it's applied before upstream request is done or after. Pre are mainly for request handling while post are mainly for response handling.
**Args**, List of args that are passed to the middleware factory function.

#### Default middlewares
|Name|Description|Args|
|---|---|---|
|setHeader|Sets header on request or response|key,value|
|removeHeader|Removes header from request or response|key|
|moveHeader|Moves a headers value to another and removes original header|sourceKey,destinationKey|
|setPath|Sets an absolute path for upstream request|path|
|setPathPrefix|Adds a prefix to the incoming requests path|prefix|
|setRequestParameter|Adds a query parameter to the upstream request|key,value|
|removeRequestParameter|Removes a query parameter from the upstream request|key|
|dedupeResponseHeaders|Ensures that specified header only has one value, value with index 0 will be used|keys...|
|redirect|Redirects request|statusCode,destination|

### Route

Routes has four attributes.

**Name**, plainly an identifier for the route.
**Predicates**, object (described more in detail further down) that determines which requests will use this route.
**Address**, which upstream URL to use for requests. Can use the `env:<environment variable key>` to load ad-hoc from environment when route is triggered.
**Middlewares**, list of middlewares to apply to route.

### Predicates
Currently predicates only support a static path.
