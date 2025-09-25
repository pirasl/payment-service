package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
)

type Task struct {
	Payload []byte
}

type workerPool struct {
	amqpClient  *rabbitMQClient
	workerCount int
	wg          *errgroup.Group
	ctx         context.Context
	cancel      context.CancelFunc
}

func newWorkerPool(numWorkers int, rabbitmq *rabbitMQClient) (*workerPool, error) {

	ctx, cancel := context.WithCancel(context.Background())

	pool := &workerPool{
		workerCount: numWorkers,
		amqpClient:  rabbitmq,
		ctx:         ctx,
		cancel:      cancel,
	}

	pool.wg, pool.ctx = errgroup.WithContext(pool.ctx)

	for index := range pool.workerCount {
		pool.wg.Go(func() error {
			return pool.worker(index + 1)
		})
	}
	return pool, nil
}

func (p *workerPool) worker(workerID int) error {
	retryCount := 0
	maxRetries := 5
	baseDelay := time.Second

	for {
		select {
		case <-p.ctx.Done():
			log.Printf("Worker %d shutting down due to context cancellation", workerID)
			return p.ctx.Err()
		default:
		}

		err := p.consumeMessages(workerID)
		if err == nil {
			retryCount = 0
			continue
		} else if strings.Contains(err.Error(), "message channel closed") {
			return nil
		}

		if !p.isRecoverableError(err) {
			log.Printf("Worker %d encountered non-recoverable error: %v", workerID, err)
			return fmt.Errorf("worker %d failed with non-recoverable error: %w", workerID, err)
		}

		retryCount++
		if retryCount >= maxRetries {
			log.Printf("Worker %d exceeded max retries (%d), giving up", workerID, maxRetries)
			return fmt.Errorf("worker %d exceeded max retries: %w", workerID, err)
		}

		delay := time.Duration(retryCount) * baseDelay
		log.Printf("Worker %d retrying in %v (attempt %d/%d)", workerID, delay, retryCount, maxRetries)

		select {
		case <-time.After(delay):
			continue
		case <-p.ctx.Done():
			return p.ctx.Err()
		}
	}
}

func (p *workerPool) consumeMessages(workerID int) error {
	ch, err := p.amqpClient.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel for worker %d: %w", workerID, err)
	}
	defer ch.Close()

	err = ch.Qos(1, 0, false)
	if err != nil {
		return fmt.Errorf("failed to set QoS for worker %d: %w", workerID, err)
	}

	msgs, err := ch.Consume(
		"stripe_processing",                // queue name
		fmt.Sprintf("worker-%d", workerID), // consumer tag unique
		false,                              // auto-ack
		false,                              // exclusive
		false,                              // no-local
		false,                              // no-wait
		nil,                                // args
	)
	if err != nil {
		return fmt.Errorf("failed to consume messages for worker %d: %w", workerID, err)
	}

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("worker %d: message channel closed", workerID)
			}

			if err := p.processMessage(workerID, msg); err != nil {
				log.Printf("Worker %d failed to process message: %v", workerID, err)
				msg.Nack(false, false)
			} else {
				msg.Ack(false)
			}

		case <-p.ctx.Done():
			log.Printf("Worker %d stopping message consumption", workerID)
			return nil
		}
	}
}

func (p *workerPool) processMessage(workerID int, msg amqp091.Delivery) error {
	log.Printf("Worker %d processing message: %s", workerID, string(msg.Body))

	if len(msg.Body) == 0 {
		return fmt.Errorf("empty message body")
	}

	// Process message here
	// Switch message.Type

	return nil
}

func (p *workerPool) isRecoverableError(err error) bool {

	if err == nil {
		return true
	}

	if _, ok := err.(*amqp091.Error); ok {
		return true
	}

	if _, ok := err.(*net.OpError); ok {
		return true
	}

	if err == context.DeadlineExceeded {
		return true
	}

	if err == context.Canceled {
		return false
	}

	return true
}

func (p *workerPool) Shutdown() error {
	p.cancel()
	return p.wg.Wait()
}
