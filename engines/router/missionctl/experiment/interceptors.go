package experiment

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/gojek/turing/engines/experiment/runner"
	"github.com/gojek/turing/engines/router/missionctl/instrumentation/metrics"
)

type ctxKey string

const (
	startTimeKey ctxKey = "startTimeKey"
)

// MetricsInterceptor is the structural interceptor used for capturing metrics
// from experiment runs
type MetricsInterceptor struct{}

// NewMetricsInterceptor is a creator for a MetricsInterceptor
func NewMetricsInterceptor() runner.Interceptor {
	return &MetricsInterceptor{}
}

// BeforeDispatch associates the start time to the context
func (i *MetricsInterceptor) BeforeDispatch(
	ctx context.Context,
) context.Context {
	return context.WithValue(ctx, startTimeKey, time.Now())
}

// AfterCompletion logs the time taken for the component to process the request,
// to the metrics collector
func (i *MetricsInterceptor) AfterCompletion(
	ctx context.Context,
	err error,
) {
	labels := make(map[string]string)
	labels["status"] = metrics.GetStatusString(err == nil)
	fmt.Println("context value for experiment name", ctx.Value(runner.ExperimentNameKey))
	if experimentName, ok := ctx.Value(runner.ExperimentNameKey).(string); ok {
		print(">>> " + experimentName)
		labels["experiment_name"] = experimentName
	}
	labels["dummy"] = ""
	if rand.Float64() > 0.15 {
		labels["dummy"] = "greater"
	}
	if _, ok := ctx.Value("NOT_EXISTS").(string); ok {
		labels["not_exists"] = "dummy"
	}

	// Get start time
	if startTime, ok := ctx.Value(startTimeKey).(time.Time); ok {
		// Measure the time taken for the experiment run
		metrics.Glob().MeasureDurationMsSince(
			metrics.ExperimentEngineRequestMs,
			startTime,
			labels,
		)
	}
}
