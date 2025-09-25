package main

import (
	"fmt"
	"net"
	"os"

	"github.com/pirasl/payment-service/internal/data"
	payments "github.com/pirasl/payment-service/proto"
	"google.golang.org/grpc"
)

type PaymentServer struct {
	payments.UnimplementedPaymentServiceServer

	models data.Models
}

func (s *service) gRPCListen() {

	gRPCPort := getOptionalStringEnv("SERVICE_GRPC_PORT", "50001")
	appEnv := getOptionalStringEnv("APP_ENV", "developement")

	var listenAddr string

	if appEnv == "production" {
		listenAddr = fmt.Sprintf(":%s", gRPCPort)
	} else {

		listenAddr = fmt.Sprintf("localhost:%s", gRPCPort)
	}

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		s.logger.Error("failed to listen for gRPC", "err", err)
		os.Exit(1)
	}

	grpc := grpc.NewServer()

	payments.RegisterPaymentServiceServer(grpc, &PaymentServer{models: *s.models})
	s.logger.Info("gRPC server started", "port:", gRPCPort)

	if err := grpc.Serve(lis); err != nil {
		s.logger.Error("failed to listen to gRPC", "err: ", err)
		os.Exit(1)
	}

}
