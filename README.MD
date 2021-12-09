# CGS (create-go-service) 🏗️

## Purpose

The purpose of _CGS_ is to decrease the time taken to get a Go service shipped by letting engineers focus on the 
business logic. This is achieved by removing the need to rewrite the boilerplate from service to service.

_CGS_ provides instrumented, healthcheck-enabled dependencies with sane defaults but with the ability to overwrite them as 
needed.

Dependencies currently supported:

- Loading of Config
- MySQL (no instrumentation yet)
- Kafka Publisher
- Kafka Subscriber
- Redis
- Router
- HTTP Server

[Want to add more?](#extending-functionality)

## How to use?

Creating a service with zero dependencies is as simple
```go
app, err := cgs.New()
if err != nil {
	return err
}

err = app.Run(context.Background())
if err != nil {
	return err
}
```

But that's not really useful so instead, lets create an application which has a server and a router
```go
app, err := cgs.New(
	cgs.WithRouter(),
	cgs.WithServer(server.WithAddr(":3030")), 
)
if err != nil {
    return err
}

err = app.Run(ctx, app.Server().Start)
if err != nil {
    return err
}
```

This will provide a server running on `:3030` with the following endpoints:

- `/live` - indicates that this instance of the application should be destroyed or restarted. A failed liveness check 
indicates that this instance is unhealthy, not some upstream dependency. Every liveness check is also included as a 
readiness check.
- `/ready` - indicates that this instance of the application is currently unable to serve requests because of an 
upstream or some transient failure. If a readiness check fails, this instance should no longer receiver requests, but 
should not be restarted or destroyed.
- `/metrics` -  provides application metrics for Prometheus to scrape.

How about adding our application-specific routes? We can add them onto the provided router
```go
app, err := cgs.New(
	cgs.WithRouter(),
	cgs.WithServer(server.WithAddr(":3030")), 
)
if err != nil {
    return err
}

app.Add(cgs.WithRouter(
	router.WithRoute(router.Route{
		Path: "/v1/books",
		HandlerFuncs: map[string]http.HandlerFunc{
			http.MethodGet: h.List,
		},
	}),
	router.WithRoute(router.Route{
		Path: "/v1/book",
		HandlerFuncs: map[string]http.HandlerFunc{
			http.MethodGet: h.Get,
			http.MethodPost: h.Create,
		},
	}),
))

err = app.Run(ctx, app.Server().Start)
if err != nil {
    return err
}
```

So in addition to the routes provided by the default router, we will also have the following available

- `/v1/book` which accepts either a `POST` or `GET` request.
- `/v1/books` which accepts just `GET` requests.

This can then be imported in the app with the following

## Extending functionality

Functionality can be added in two ways:

1) You're adding [a new dependency](#adding-a-new-dependency) (i.e. AWS SQS).
2) You're [extending the options](#extending-options) available to an existing dependency (i.e. Make use of sql.DB's `SetConnMaxIdleTime`).

### Adding a new dependency

Adding a dependency is as simple as adding the type to `cgs.Application` and create an WithDependencyName Option 
which will instantiate the desired dependency.

### Extending options

Find the directory of the Options you wish to extend and add a sane default to the type(like what's been done for SQL's 
Connection Max Lifetime). Then provide an Option function which will allow that config to be set.

