package jaegerTracer

import (
	"hexa_micro/pkg/shortenservice/container/logger"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

const (
	SERVICE_NAME_SHORTEN_URL string = "ShortenURL Service"
)

// func newTracer(serviceName string) (opentracing.Tracer, io.Closer) {
// 	// use default port (6831) and default max packet length (65000)
// 	sender, err := jaeger.NewUDPTransport("", 0)
// 	if err != nil {
// 		logger.Log.Fatalf("%+v", err)
// 	}

// 	metrics := jaeger.NewNullMetrics()
// 	reportMetrics := jaeger.ReporterOptions.Metrics(metrics)

// 	reportLogger := jaeger.ReporterOptions.Logger(jaeger.StdLogger)

// 	reporter := jaeger.NewRemoteReporter(sender, reportMetrics, reportLogger)
// 	sampler := jaeger.NewConstSampler(true)

// 	//create tracer Jaeger tracer
// 	tracer, closer := jaeger.NewTracer(serviceName, sampler, reporter)

// 	opentracing.SetGlobalTracer(tracer)

// 	return tracer, closer
// }

// initJaeger returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func NewTracer() (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		ServiceName: SERVICE_NAME_SHORTEN_URL,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
		Disabled: true,
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		logger.Log.Fatalf("Cannot init Jaeger: %+v", err)
	}

	opentracing.SetGlobalTracer(tracer)

	return tracer, closer
}
