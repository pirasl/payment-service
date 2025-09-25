package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type jwtConfig struct {
	secret string
}

func (s *service) serve() error {
	var addr string
	env := getOptionalStringEnv("APP_ENV", "development")
	if env == "development" {
		addr = fmt.Sprintf("localhost:%d", s.config.servicePort)

	} else {
		addr = fmt.Sprintf(":%d", s.config.servicePort)
	}

	server := &http.Server{
		Addr:         addr,
		Handler:      s.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sign := <-quit

		s.logger.Info("caught signal", "signal", sign.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		s.logger.Info("completing background tasks", "addr", server.Addr)

		s.logger.Info("closing RabbitMQ connection")
		s.rabbitmqClient.channel.Close()
		s.rabbitmqClient.conn.Close()

		s.logger.Info("closing workerpool")
		s.workerPool.Shutdown()

		s.wg.Wait()
		shutdownError <- nil
	}()

	s.logger.Info("starting server", "addr", server.Addr)

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	s.logger.Info("stopped server", "addr", server.Addr)

	return nil
}
