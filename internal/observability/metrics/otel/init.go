package otel

import (
	"context"
	"crypto/tls"
	"errors"
	"os"
	"strconv"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"google.golang.org/grpc/credentials"
)

var (
	mp         *sdkmetric.MeterProvider
	once       sync.Once
	shutdownFn = func(context.Context) error { return nil }
)

func initMeter(ctx context.Context, serviceName string) error {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		return errors.New("OTEL_EXPORTER_OTLP_ENDPOINT not set")
	}

	opts := []otlpmetricgrpc.Option{otlpmetricgrpc.WithEndpoint(endpoint)}
	insecure, _ := strconv.ParseBool(os.Getenv("OTEL_EXPORTER_OTLP_INSECURE"))
	if insecure {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	} else {
		creds := credentials.NewTLS(
			&tls.Config{MinVersion: tls.VersionTLS12},
		)
		opts = append(opts, otlpmetricgrpc.WithTLSCredentials(creds))
	}

	exp, err := otlpmetricgrpc.New(ctx, opts...)
	if err != nil {
		return err
	}

	version := os.Getenv("SERVICE_VERSION")
	if version == "" {
		version = "0.0.1"
	}

	res, err := sdkresource.New(ctx,
		sdkresource.WithFromEnv(),
		sdkresource.WithProcess(),
		sdkresource.WithTelemetrySDK(),
		sdkresource.WithHost(),
		sdkresource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(version),
		),
	)
	if err != nil {
		return err
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp,
			sdkmetric.WithInterval(10*time.Second),
			sdkmetric.WithTimeout(5*time.Second),
		)),
	)
	otel.SetMeterProvider(provider)

	mp = provider

	shutdownFn = func(ctx context.Context) error {
		return mp.Shutdown(ctx)
	}

	return nil
}
