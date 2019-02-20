package config

import (
	"time"
)

// GatewayConfig defines the configuration for the gateway
type GatewayConfig struct {
	// Port is the port used by the gateway service
	Port int `mapstructure:"port" validate:"required"`

	// Hostname is the hostname used by the gateway service
	Hostname string `mapstructure:"hostname" validate:"hostname"`

	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body.
	//
	// Because ReadTimeout does not let Handlers make per-request
	// decisions on each request body's acceptable deadline or
	// upload rate, most users will prefer to use
	// ReadHeaderTimeout. It is valid to use them both.
	ReadTimeout time.Duration `mapstructure:"read_timeout"`

	// ReadHeaderTimeout is the amount of time allowed to read
	// request headers. The connection's read deadline is reset
	// after reading the headers and the Handler can decide what
	// is considered too slow for the body.
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read. Like ReadTimeout, it does not
	// let Handlers make decisions on a per-request basis.
	WriteTimeout time.Duration `mapstructure:"write_timeout"`

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alives are enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, ReadHeaderTimeout is used.
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`

	// ShutdownTimeout defines the time (in seconds) in which the server would shutdown
	// On the os.Interrupt event
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`

	TLSCertPath string `mapstructure:"tls_cert_path"`

	// Router defines the router configuration
	Router *RouterConfig `mapstructure:"router"`
}

// RouterConfig contains information about the router used in the gateway
type RouterConfig struct {
	// Name gets the router by it's registered name
	Name string `mapstructure:"name"`

	// DefaultMiddlewares are the middlewares used as default for each endpoint
	// without middlewares set from the
	DefaultMiddlewares []string `mapstructure:"default_middlewares"`

	// Prefix is the url prefix for the API
	Prefix string `mapstructure:"prefix"`

	// DefaultPagination defines default ListPagination for the gateway
	DefaultPagination *Pagination `mapstructure:"default_pagination"`
}
