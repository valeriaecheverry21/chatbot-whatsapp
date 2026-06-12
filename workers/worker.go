package workers

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"

	"github.com/user/chatbot-whatsapp/models"
	"github.com/user/chatbot-whatsapp/queue"
	"github.com/user/chatbot-whatsapp/services"
)

type WorkerPool struct {
	concurrency    int
	queueClient    *queue.Client
	whatsAppService *services.WhatsAppService
	messageService *services.MessageService
	log            *zap.Logger
}

func NewWorkerPool(
	concurrency int,
	queueClient *queue.Client,
	whatsAppService *services.WhatsAppService,
	messageService *services.MessageService,
	log *zap.Logger,
) *WorkerPool {
	return &WorkerPool{
		concurrency:     concurrency,
		queueClient:     queueClient,
		whatsAppService: whatsAppService,
		messageService:  messageService,
		log:             log,
	}
}

func (wp *WorkerPool) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	for i := 0; i < wp.concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			wp.runWorker(ctx, id)
		}(i)
	}

	wp.log.Info("worker pool started",
		zap.Int("concurrency", wp.concurrency),
	)

	select {
	case sig := <-sigChan:
		wp.log.Info("shutdown signal received", zap.String("signal", sig.String()))
	case <-ctx.Done():
	}

	cancel()
	wg.Wait()
	wp.log.Info("all workers stopped")
}

func (wp *WorkerPool) runWorker(ctx context.Context, id int) {
	log := wp.log.With(zap.Int("worker_id", id))
	log.Info("worker started")

	for {
		select {
		case <-ctx.Done():
			log.Info("worker stopping")
			return
		default:
		}

		msg, err := wp.queueClient.Consume(ctx)
		if err != nil {
			continue
		}

		log.Info("processing message",
			zap.String("phone", msg.Phone),
			zap.String("type", string(msg.MessageType)),
		)

		if err := wp.processMessage(ctx, msg); err != nil {
			log.Error("failed to process message, retrying",
				zap.String("phone", msg.Phone),
				zap.Error(err),
			)

			if retryErr := wp.queueClient.PublishRetry(ctx, msg); retryErr != nil {
				log.Error("failed to publish retry", zap.Error(retryErr))
			}
		}
	}
}

func (wp *WorkerPool) processMessage(ctx context.Context, msg *queue.QueueMessage) error {
	var record *models.Message
	var err error

	switch msg.MessageType {
	case models.MessageTypeText:
		record, err = wp.whatsAppService.SendText(ctx, services.SendTextInput{
			Phone:   msg.Phone,
			Message: msg.Body,
			Preview: true,
		})

	case models.MessageTypeImage:
		record, err = wp.whatsAppService.SendImage(ctx, services.SendImageInput{
			Phone:   msg.Phone,
			URL:     msg.MediaURL,
			Caption: msg.Body,
		})

	case models.MessageTypeTemplate:
		if msg.Template == nil {
			return nil
		}

		record, err = wp.whatsAppService.SendTemplate(ctx, services.SendTemplateInput{
			Phone:    msg.Phone,
			Name:     msg.Template.Name,
			Language: msg.Template.Language,
			Params:   msg.Template.Params,
		})

	default:
		record, err = wp.whatsAppService.SendText(ctx, services.SendTextInput{
			Phone:   msg.Phone,
			Message: msg.Body,
			Preview: true,
		})
	}

	if err != nil {
		return err
	}

	if record != nil && msg.CustomerID > 0 {
		record.CustomerID = msg.CustomerID
		if saveErr := wp.messageService.Save(record); saveErr != nil {
			wp.log.Error("failed to save outbound message",
				zap.String("phone", msg.Phone),
				zap.Error(saveErr),
			)
		}
	}

	return nil
}
