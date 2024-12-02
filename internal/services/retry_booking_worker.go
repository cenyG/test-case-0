package services

import (
	"applicationDesignTest/internal/model"
	"applicationDesignTest/internal/usecases"
	"context"
	"errors"
	"log/slog"
	"sync/atomic"
	"time"
)

type RetryConfig struct {
	RetryCount  int           // RetryCount - times to retry
	Delay       time.Duration // Delay - queue delay between executions
	QueueLength int           // QueueLength - max queue length
	TTL         time.Duration // TTL - will retry only if current_time - created_at < TTL
}

type RetryEntity struct {
	order     model.Order
	createdAt time.Time
}

type RetryBookingWorkerService struct {
	useCase     usecases.RetryBookingUseCase
	queue       chan RetryEntity
	retryConfig RetryConfig

	workerRunning atomic.Bool
}

func NewRetryCreateOrderWorker(useCase usecases.RetryBookingUseCase, retryConfig RetryConfig) RetryBookingWorkerService {
	return RetryBookingWorkerService{
		useCase:     useCase,
		queue:       make(chan RetryEntity, retryConfig.QueueLength),
		retryConfig: retryConfig,
	}
}

func (r *RetryBookingWorkerService) Enqueue(order model.Order) {
	if r.workerRunning.Load() {
		r.queue <- RetryEntity{
			order:     order,
			createdAt: time.Now(),
		}
	} else {
		slog.Warn("try to enqueue order: %d but worker is stopped", order.ID)
	}
}

func (r *RetryBookingWorkerService) Start(ctx context.Context) {
	r.workerRunning.Store(true)
	defer r.onWorkerStop()

	for {
		select {
		case <-ctx.Done():
			slog.Error("[worker] stop worker by context: %v", ctx.Err())
			return
		case order := <-r.queue:
			err := r.retryCreateOrder(ctx, order)
			if err != nil {
				slog.Error("[worker] fail retry create order %v: %v", order, err)
			}
		default:
		}

		time.Sleep(r.retryConfig.Delay)
	}
}

func (r *RetryBookingWorkerService) onWorkerStop() {
	r.workerRunning.Store(false)
	for {
		select {
		case <-r.queue:
		default:
			close(r.queue)
			return
		}
	}
}

func (r *RetryBookingWorkerService) retryCreateOrder(ctx context.Context, entity RetryEntity) error {
	for remain := r.retryConfig.RetryCount; remain >= 0; remain-- {
		if time.Since(entity.createdAt) > r.retryConfig.TTL {
			return errors.New("canceled by TTL")
		}

		err := r.useCase.RetryBooking(ctx, entity.order)
		if err == nil {
			return nil
		}

		if remain == 0 {
			return err
		}

		time.Sleep(r.retryConfig.Delay)
	}
	return nil
}
