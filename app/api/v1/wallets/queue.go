package apiv1wallets

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"go.uber.org/fx"
	"os"
)

const queueName = "quikgaming"

type Queue struct {
	amqpConnection *amqp.Connection
	walletsRepo IWalletsRepository
	walletsCache IWalletsCache
}

func NewQueue(lc fx.Lifecycle, walletsRepo IWalletsRepository, walletsCache IWalletsCache) IQueue {
	if conn, err := amqp.Dial(os.Getenv("QUIKWALLET_RABBITMQ_DSN")); err != nil {
		panic(err)
	} else {
		lc.Append(
			fx.Hook{
				OnStop: func(ctx context.Context) error {
					return conn.Close()
				},
			},
		)

		return &Queue{
			amqpConnection: conn,
			walletsRepo: walletsRepo,
			walletsCache: walletsCache,
		}
	}
}

type NewWalletLogJob struct {
	WalletUUID string `json:"wallet_uuid"`
}

func (queue *Queue) PublishNewWalletLogJob(walletUUID uuid.UUID) error {
	if ch, err := queue.amqpConnection.Channel(); err != nil {
		return err
	} else {
		defer ch.Close()
		if q, err := getQueue(ch); err != nil {
			return err
		} else if jobJson, err := json.Marshal(&NewWalletLogJob{WalletUUID: walletUUID.String()}); err != nil {
			return err
		} else if err := ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        jobJson,
			},
		); err != nil {
			return err
		} else {
			return nil
		}
	}
}

func (queue *Queue) ConsumeNewWalletLogJob() error {
	if ch, err := queue.amqpConnection.Channel(); err != nil {
		return err
	} else {
		defer ch.Close()
		if q, err := getQueue(ch); err != nil {
			return err
		} else if msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		); err != nil{
			return err
		} else {
			for d := range msgs {
				log.Info().Msgf("Received a message: %s", string(d.Body))

				var job NewWalletLogJob
				if err := json.Unmarshal(d.Body, &job); err != nil {
					return err
				} else if walletUUID, err := uuid.Parse(job.WalletUUID); err != nil {
					return err
				} else if err := queue.walletsRepo.ProcessWalletLogs(walletUUID); err != nil {
					return err
				} else if err := queue.walletsCache.InvalidateWalletBalance(walletUUID); err != nil {
					return err
				} else {
					continue
				}
			}
		}
		return nil
	}
}

func getQueue(channel *amqp.Channel) (amqp.Queue, error) {
	return channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
}
