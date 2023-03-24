# NAGG (Not Another Go Gateway)

This is an experimental API Gateway that probably shouldn't have been created. But here it is.


## Configuration

See `examples/gateway.json` for an example. 

Configuration can be loaded by either setting the environment variable  `NAGG_CONFIG_PATH` to a json file or one can set (has precedence) the environment variable `NAGG_CONFIG` with a JSON string that will be loaded. In the latter case standard Golang expansion `${...}` and `$...` will be expanded with availabe environment variables.


### Gateway
Global space for configuration in json structure and have two attributes that can be configured.

**Middlewares**, global middlewares to apply on each request and response.
**Routes**, routes to handle requests for in the gateway.

### Middleware
Middlewares has three attributes.

**Name**, which determines which middleware is loaded.
**Phase**, either pre or post which determines if it's applied before upstream request is done or after. Pre are mainly for request handling while post are mainly for response handling.
**Args**, List of args that are passed to the middleware factory function.


### Route

Routes has four attributes.

**Name**, plainly an identifier for the route.
**Predicates**, object (described more in detail further down) that determines which requests will use this route.
**Address**, which upstream URL to use for requests. Can use the `env:<environment variable key>` to load ad-hoc from environment when route is triggered.
**Middlewares**, list of middlewares to apply to route.

### Predicates
Currently predicates only support a static path.
