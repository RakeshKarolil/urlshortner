package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Counter for number of shortened URLs
	URLsShortened = promauto.NewCounter(prometheus.CounterOpts{
		Name: "urlshortener_urls_shortened_total",
		Help: "Total number of URLs shortened",
	})

	// Counter for number of redirect requests
	Redirects = promauto.NewCounter(prometheus.CounterOpts{
		Name: "urlshortener_redirects_total",
		Help: "Total number of URL redirects",
	})
)
