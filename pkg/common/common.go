package common

import (
	"math/rand"
	"net/http"
	"os"
	"time"
)

// https://pkg.go.dev/net/http#pkg-overview
// Clients and Transports are safe for concurrent use by multiple goroutines
// and for efficiency should only be created once and re-used.

const ZMQ = true

var ExpirationTTLms = os.Getenv("EXPIRATION_TTL")

func init() {
	rand.Seed(time.Now().UnixNano())
}

var MyName = os.Getenv("APP_NAME_NO_UNDERSCORES")
var MyRawName = os.Getenv("APP_RAW_NAME_NO_UNDERSCORES")

var HTTPClient = &http.Client{
	Transport: &http.Transport{MaxConnsPerHost: 100, MaxIdleConnsPerHost: 100, MaxIdleConns: 100},
	Timeout:   60 * time.Second,
}
