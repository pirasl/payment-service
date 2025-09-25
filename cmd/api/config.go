package main

import (
	_ "github.com/lib/pq"
)

type serviceConfig struct {
	servicePort int
	gRPCPort    int

	jwtConfig         *jwtConfig
	rateLimiterConfig *rateLimiterConfig
}

type rateLimiterConfig struct {
	enabled bool
	rps     int
	burst   int
}

func newJWTConfig() (*jwtConfig, error) {
	secret, err := getRequiredStringEnv("JWT_SECRET")
	if err != nil {
		return nil, err
	}

	return &jwtConfig{
		secret: *secret,
	}, nil
}

func newRateLimiterConfig() *rateLimiterConfig {
	return &rateLimiterConfig{
		enabled: getOptionalBoolEnv("RATE_LIMITER_ENABLED", true),
		rps:     getOptionalIntEnv("RATE_LIMITER_RPS", 5),
		burst:   getOptionalIntEnv("RATE_LIMITER_BURST", 10),
	}
}

func newServiceConfig() (*serviceConfig, error) {

	servicePort := getOptionalIntEnv("SERVICE_PORT", 8080)
	grpcPort := getOptionalIntEnv("SERVICE_GRPC_PORT", 8080)

	rateLimiterConfig := newRateLimiterConfig()

	jwtConfig, err := newJWTConfig()
	if err != nil {
		return nil, err
	}

	serviceConfig := &serviceConfig{
		servicePort:       servicePort,
		gRPCPort:          grpcPort,
		rateLimiterConfig: rateLimiterConfig,
		jwtConfig:         jwtConfig,
	}

	return serviceConfig, nil
}
