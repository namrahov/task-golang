package kafka

import (
	"context"
	"encoding/json"
	"task-golang/config"
	"task-golang/model"
	"time"

	log "github.com/sirupsen/logrus"
)

func RunStatusConsumer(ctx context.Context, c *config.StatusConsumer) error {
	defer c.Reader.Close()
	for {
		m, err := c.Reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil // shutting down
			}
			log.WithError(err).Warn("kafka fetch error")
			time.Sleep(1 * time.Second)
			continue
		}

		// key/value are []byte; your Java listener used StringDeserializer
		msg := string(m.Value)
		log.Infof("Processing status change message: %s (partition=%d offset=%d key=%s)",
			msg, m.Partition, m.Offset, string(m.Key))

		var dto model.TestDto
		if err := json.Unmarshal(m.Value, &dto); err != nil {
			log.WithError(err).Error("Error parsing message JSON")
			// decide: skip, DLQ, or retry. Here we just skip & commit to avoid re-reading bad message.
		} else {
			// TODO: call your service with dto (e.g., taskService.HandleStatusChange(ctx, dto))
			log.Infof("Parsed DTO: %+v", dto)
		}

		// Commit only after successful processing (at-least-once)
		if err := c.Reader.CommitMessages(ctx, m); err != nil {
			log.WithError(err).Warn("commit error")
		}
	}
}

func RunStatusConsumer2(ctx context.Context, c *config.StatusConsumer) error {
	defer c.Reader.Close()
	for {
		m, err := c.Reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil // shutting down
			}
			log.WithError(err).Warn("kafka fetch error")
			time.Sleep(1 * time.Second)
			continue
		}

		// key/value are []byte; your Java listener used StringDeserializer
		msg := string(m.Value)
		log.Infof("Processing status change message: %s (partition=%d offset=%d key=%s)",
			msg, m.Partition, m.Offset, string(m.Key))

		var dto model.TestDto
		if err := json.Unmarshal(m.Value, &dto); err != nil {
			log.WithError(err).Error("Error parsing message JSON")
			// decide: skip, DLQ, or retry. Here we just skip & commit to avoid re-reading bad message.
		} else {
			// TODO: call your service with dto (e.g., taskService.HandleStatusChange(ctx, dto))
			log.Infof("Parsed DTO: %+v", dto)
		}

		// Commit only after successful processing (at-least-once)
		if err := c.Reader.CommitMessages(ctx, m); err != nil {
			log.WithError(err).Warn("commit error")
		}
	}
}
