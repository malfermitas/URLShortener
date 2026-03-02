package metrics

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// Exposed Prometheus metrics (exported for use from other packages)
var (
	HttpRequestsTotal          *prometheus.CounterVec
	HttpRequestDurationSeconds *prometheus.HistogramVec

	UrlsCreatedTotal prometheus.Counter

	CacheHitsTotal   prometheus.Counter
	CacheMissesTotal prometheus.Counter
	CacheSetsTotal   prometheus.Counter

	RedirectsTotal prometheus.Counter

	UrlRedirectsTotal *prometheus.CounterVec

	AnalyticsQueriesTotal prometheus.Counter
	HitsTotal             prometheus.Counter

	WebUIPageviewsTotal prometheus.Counter
)

// InitMetrics initializes and registers all Prometheus metrics.
func InitMetrics() {
	HttpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total HTTP requests",
	}, []string{"path", "method", "status"})

	HttpRequestDurationSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"path", "method"})

	UrlsCreatedTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "urls_created_total",
		Help: "Total URLs created",
	})

	CacheHitsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cache_hits_total",
		Help: "Cache hits",
	})
	CacheMissesTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cache_misses_total",
		Help: "Cache misses",
	})
	CacheSetsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cache_sets_total",
		Help: "Cache set operations",
	})

	RedirectsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "redirects_total",
		Help: "Total redirects",
	})

	UrlRedirectsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "url_redirects_total",
		Help: "Total redirects per shortened URL",
	}, []string{"short_key"})

	AnalyticsQueriesTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "analytics_queries_total",
		Help: "Total analytics queries",
	})
	HitsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "hits_total",
		Help: "Total URL hits recorded",
	})

	WebUIPageviewsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "webui_pageviews_total",
		Help: "Total Web UI page views",
	})

	// Register all metrics with the default registry
	prometheus.MustRegister(
		HttpRequestsTotal,
		HttpRequestDurationSeconds,
		UrlsCreatedTotal,
		CacheHitsTotal,
		CacheMissesTotal,
		CacheSetsTotal,
		RedirectsTotal,
		UrlRedirectsTotal,
		AnalyticsQueriesTotal,
		HitsTotal,
		WebUIPageviewsTotal,
	)
}

// MetricsMiddleware returns a Gin middleware that records HTTP metrics.
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method
		start := time.Now()
		c.Next()
		status := strconv.Itoa(c.Writer.Status())
		if HttpRequestsTotal != nil {
			HttpRequestsTotal.WithLabelValues(path, method, status).Inc()
		}
		if HttpRequestDurationSeconds != nil {
			HttpRequestDurationSeconds.WithLabelValues(path, method).Observe(time.Since(start).Seconds())
		}
	}
}
