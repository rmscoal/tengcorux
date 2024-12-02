// This package is just a wrapper for go-resty/resty/v2. It is unnecessary to
// reinvent the wheel just because you want to add some little feature. The extra
// feature this package has is the tracing ability.
//
// When opted, it will start a new span before requests and captures attributes.
// Then, after request it will end the span marking the process has finished as
// well as capturing attributes for the span.

package rest
