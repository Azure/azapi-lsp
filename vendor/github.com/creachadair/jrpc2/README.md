# jrpc2

[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=yellow)](https://pkg.go.dev/github.com/creachadair/jrpc2)
[![Go Report Card](https://goreportcard.com/badge/github.com/creachadair/jrpc2)](https://goreportcard.com/report/github.com/creachadair/jrpc2)

This repository provides Go package that implements a [JSON-RPC 2.0][spec] client and server.
There is also a working [example in the Go playground](https://play.golang.org/p/MSClCk55UzF).

## Packages

*  Package [jrpc2](http://godoc.org/github.com/creachadair/jrpc2) implements the base client and server.

*  Package [channel](http://godoc.org/github.com/creachadair/jrpc2/channel) defines the communication channel abstraction used by the server & client.

*  Package [code](http://godoc.org/github.com/creachadair/jrpc2/code) defines standard error codes as defined by the JSON-RPC 2.0 protocol.

*  Package [handler](http://godoc.org/github.com/creachadair/jrpc2/handler) defines support for adapting functions to service methods.

*  Package [jctx](http://godoc.org/github.com/creachadair/jrpc2/jctx) implements an encoder and decoder for request context values, allowing context metadata to be propagated through JSON-RPC requests.

*  Package [jhttp](http://godoc.org/github.com/creachadair/jrpc2/jhttp) allows clients and servers to use HTTP as a transport.

*  Package [metrics](http://godoc.org/github.com/creachadair/jrpc2/metrics) defines a server metrics collector.

*  Package [server](http://godoc.org/github.com/creachadair/jrpc2/server) provides support for running a server to handle multiple connections, and an in-memory implementation for testing.

[spec]: http://www.jsonrpc.org/specification

### Versioning

This module is currently still at v0 and subject to change. To the extent practical, I try to avoid breaking changes to the API, but when I do make a breaking change I will update the minor version. For bug fixes and non-breaking minor changes I update only the patch. Hence, when going from (say) `v0.11.3` to `v0.12.0`, be advised that some API changes may occur.

I am planning to commit to a stable v1 at some point. If you have opinions about what that should mean, please feel free to comment on https://github.com/creachadair/jrpc2/issues/46.

## Implementation Notes

The following describes some of the implementation choices made by this module.

### Batch requests and error reporting

The JSON-RPC 2.0 spec is ambiguous about the semantics of batch requests. Specifically, the definition of notifications says:

> A Notification is a Request object without an "id" member.
> ...
> The Server MUST NOT reply to a Notification, including those that are within a batch request.
>
> Notifications are not confirmable by definition, since they do not have a Response object to be returned. As such, the Client would not be aware of any errors (like e.g. "Invalid params", "Internal error").

This conflicts with the definition of batch requests, which asserts:

> A Response object SHOULD exist for each Request object, except that there SHOULD NOT be any Response objects for notifications.
> ...
> The Response objects being returned from a batch call MAY be returned in any order within the Array.
> ...
> If the batch rpc call itself fails to be recognized as an valid JSON or as an Array with at least one value, the response from the Server MUST be a single Response object.

and includes examples that contain request values with no ID (which are, perforce, notifications) and report errors back to the client. Since order may not be relied upon, and there are no IDs, the client cannot correctly match such responses back to their originating requests.

This implementation resolves the conflict in favour of the notification rules. Specifically:

-  If a batch is empty or not valid JSON, the server reports error -32700 (Invalid JSON) as a single error Response object.

-  Otherwise, parse or validation errors resulting from any batch member without an ID are mapped to error objects with a `null` ID, in the same position in the reply as the corresponding request. Preservation of order is not required by the specification, but it ensures the server has stable behaviour.

Because a server is allowed to reorder the results, a client should not depend on this implementation detail.

### Non-standard server push

The specification defines client and server as follows:

> The Client is defined as the origin of `Request` objects and the handler of `Response` objects.
> The Server is defined as the origin of `Response` objects and the handler of `Request` objects.

Although a client may also be a server, and vice versa, the specification does not require them to do so. The server notification support defined in the `jrpc2` package is thus "non-standard" in that it allows the server to act as a client, and the client as a server, in the narrow context of "push" notifications and server callbacks. Otherwise the feature is not special: Requests sent by `*jrpc2.Server.Notify` and `*jrpc2.Server.Callback`are standard `Request` objects.
