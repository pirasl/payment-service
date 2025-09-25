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
