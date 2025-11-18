package worker

import (
	"context"
	"log"
	"sync"
)

type Worker[T Eventable] struct {
	consumer Consumer
	wg       sync.WaitGroup
}

func NewWorker[T Eventable](brokers []string, topic, groupID string, processor Processor[T]) *Worker[T] {
	consumer := NewConsumer(brokers, topic, groupID, processor)
	return &Worker[T]{
		consumer: consumer,
	}
}

func (w *Worker[T]) Start(ctx context.Context) {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		w.consumer.Start(ctx)
	}()
}

func (w *Worker[T]) Stop() {
	w.wg.Wait()
	log.Println("Worker stopped")
}
