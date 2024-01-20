package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var startTime = time.Now()

var (
	uptimeMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "app_uptime_seconds",
			Help: "Current instance Application uptime in seconds.",
		},
		[]string{"service"},
	)
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Current instance Total number of HTTP requests.",
		},
		[]string{"method", "path", "name"},
	)
	DrinkConsumedCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "drink_consumed_total",
			Help: "Current instance Total number of consumed action per drink.",
		},
		[]string{"name"},
	)
	IngredientsConsumedCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ingredients_consumed_total",
			Help: "Current instance Total number of each ingredient consumed.",
		},
		[]string{"name"},
	)
	MoneyEarnedCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "money_earned_total",
			Help: "Current instance Total money earned.",
		},
	)
	DrinksActiveCounter = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "drink_active_total",
			Help: "Current instance Number of activated drinks in coffeemachine.",
		},
	)
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		method := c.Request.Method
		name := c.Param("name")

		requestCounter.WithLabelValues(method, path, name).Inc()
		c.Next()
	}
}

func RegisterRoutesStatistics(router *gin.RouterGroup) {
	router.GET("/", getStatistics)
	go func() {
		for {
			updateUptimeMetric("coffeemachine")
			time.Sleep(time.Second)
		}
	}()
}

// getStatistics godoc
// @Summary Statistics endpoint (metrics)
// @Description Statistics endpoint (metrics)
// @Produce plain
// @Success 200
// @Failure 500
// @Router /statistics [get]
// @Security BearerAuth
func getStatistics(c *gin.Context) {
	if code, err := parseToken(c); err != nil {
		c.String(code, "")
		return
	}
	promhttp.HandlerFor(customRegistry(), promhttp.HandlerOpts{}).ServeHTTP(c.Writer, c.Request)
}

func updateUptimeMetric(serviceName string) {
	uptime := time.Since(startTime).Seconds()
	uptimeMetric.WithLabelValues(serviceName).Set(uptime)
}

func customRegistry() *prometheus.Registry {
	reg := prometheus.NewRegistry()
	reg.MustRegister(uptimeMetric)
	reg.MustRegister(requestCounter)
	reg.MustRegister(DrinkConsumedCounter)
	reg.MustRegister(IngredientsConsumedCounter)
	reg.MustRegister(MoneyEarnedCounter)
	reg.MustRegister(DrinksActiveCounter)
	return reg
}
