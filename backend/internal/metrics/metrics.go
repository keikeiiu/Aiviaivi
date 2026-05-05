package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	WSConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "ws_connections_active",
			Help: "Number of active WebSocket connections.",
		},
	)

	DanmakuSentTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "danmaku_sent_total",
			Help: "Total number of danmaku messages sent.",
		},
	)

	VideosUploadedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "videos_uploaded_total",
			Help: "Total number of videos uploaded.",
		},
	)

	VideosTranscodedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "videos_transcoded_total",
			Help: "Total number of videos successfully transcoded.",
		},
	)
)

// IncWS increments the active WebSocket connection count.
func IncWS() { WSConnectionsActive.Inc() }

// DecWS decrements the active WebSocket connection count.
func DecWS() { WSConnectionsActive.Dec() }

// IncDanmaku increments the danmaku sent counter.
func IncDanmaku() { DanmakuSentTotal.Inc() }

// IncUpload increments the video upload counter.
func IncUpload() { VideosUploadedTotal.Inc() }

// IncTranscoded increments the successful transcode counter.
func IncTranscoded() { VideosTranscodedTotal.Inc() }
