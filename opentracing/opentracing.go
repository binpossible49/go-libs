package opentracing

import (
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

// InitOpentracing initializes opentracing
func InitOpentracing(addr, name string) error {
	cfg := config.Configuration{
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: addr,
			LogSpans:           false,
		},
	}
	_, err := cfg.InitGlobalTracer(
		name,
		config.Logger(jaeger.StdLogger),
		config.Sampler(jaeger.NewConstSampler(true)),
		config.Metrics(metrics.NullFactory),
	)
	if err != nil {
		return err
	}
	return nil
}
