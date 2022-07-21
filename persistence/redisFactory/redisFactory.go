package redisFactory

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/onqlavetest1/config"
	model2 "github.com/onqlavetest1/model"
	"github.com/onqlavetest1/persistence/handler"
	"github.com/onqlavetest1/ulti"
	"github.com/pborman/uuid"
	"sync"
	"time"
)

var (
	ERR_CONNECT_SERVER    = errors.New("Can't connect to server")
	ERR_INITIALIZE_CLIENT = errors.New("Can't initialize client")
	ERR_DATA_INVALID      = errors.New("Data is invalid")
	ERR_PUBLISH_GENERAL   = errors.New("Can't publish data")
	ERR_SUB_GENERAL       = errors.New("Can't sub to topic")
	ERR_SUB_WAITING       = errors.New("SubClient isn't running")
)

type handlerGroup struct {
	Receiver ulti.ReceiveHandlerFunction
	Filter   ulti.FilterFunction
}

type RedisHandlerFactory struct {
	Configs        config.REDISHandlerFactoryConfigs
	subscribers    map[string]*redisHandler
	subscriberLock *sync.Mutex
}

func NewRedisHandlerFactory(c config.REDISHandlerFactoryConfigs) RedisHandlerFactory {
	return RedisHandlerFactory{
		Configs:        config.REDISHandlerFactoryConfigs{},
		subscribers:    map[string]*redisHandler{},
		subscriberLock: &sync.Mutex{},
	}
}

type redisHandler struct {
	factory        *RedisHandlerFactory
	subscriptions  map[string]handlerGroup
	topics         map[string]any
	client         *redis.Client
	running        bool
	mu             *sync.Mutex
	cancelFunction context.CancelFunc
}

func (r redisHandler) PublishCommand(ctx context.Context, command model2.DomainCommand) error {
	if r.client == nil {
		return ERR_INITIALIZE_CLIENT
	}

	_, err := r.client.Ping(ctx).Result()

	if err != nil {
		return ERR_CONNECT_SERVER
	}
	payload, err := json.Marshal(command.Metadata())

	if err != nil {
		return ERR_DATA_INVALID
	}

	rsult := r.client.Publish(ctx, command.Binding(), payload)

	if rsult.Err() != nil {
		return ERR_PUBLISH_GENERAL
	}

	return nil
}

func (r redisHandler) PublishEvent(ctx context.Context, event model2.DomainEvent) error {
	if r.client == nil {
		return ERR_INITIALIZE_CLIENT
	}

	_, err := r.client.Ping(ctx).Result()

	if err != nil {
		return ERR_CONNECT_SERVER
	}
	payload, err := json.Marshal(event.Metadata())

	if err != nil {
		return ERR_DATA_INVALID
	}

	rsult := r.client.Publish(ctx, event.Topic(), payload)

	if rsult.Err() != nil {
		return ERR_PUBLISH_GENERAL
	}

	return nil
}

func (r RedisHandlerFactory) GetSubcriber(subscriberName string) *redisHandler {
	if rhdl, ok := r.subscribers[subscriberName]; ok {
		return rhdl
	}
	return nil
}

func (r RedisHandlerFactory) Unsubcriber(subscriberName string) error {
	if rhdl, ok := r.subscribers[subscriberName]; ok {
		if rhdl.running {
			err := rhdl.Stop()
			if err != nil {
				return err
			}
		}
		r.subscriberLock.Lock()
		delete(r.subscribers, subscriberName)
		r.subscriberLock.Unlock()
	}
	return nil
}

func (r RedisHandlerFactory) Subscriber(subscriberName string) handler.RecieveAndCreateHandler {
	if r.GetSubcriber(subscriberName) != nil {
		return r.GetSubcriber(subscriberName)
	}

	handler := &redisHandler{
		factory:        &r,
		subscriptions:  map[string]handlerGroup{},
		topics:         map[string]any{},
		client:         nil,
		cancelFunction: nil,
		mu:             &sync.Mutex{},
	}

	r.subscriberLock.Lock()
	r.subscribers[subscriberName] = handler
	r.subscriberLock.Unlock()

	return handler
}

func (r RedisHandlerFactory) Publisher() handler.CreateHandler {
	return redisHandler{
		factory:        &r,
		subscriptions:  nil,
		topics:         map[string]any{},
		client:         nil,
		cancelFunction: nil,
	}
}

func (r redisHandler) Create() handler.Handler {
	r.client = redis.NewClient(&redis.Options{
		Network:      r.factory.Configs.Network,
		Addr:         r.factory.Configs.Addr,
		Username:     r.factory.Configs.Username,
		Password:     r.factory.Configs.Password,
		DB:           r.factory.Configs.DB,
		MaxRetries:   r.factory.Configs.MaxRetries,
		DialTimeout:  time.Duration(r.factory.Configs.AckDeadline * 1000),
		ReadTimeout:  time.Duration(r.factory.Configs.AckDeadline * 1000),
		WriteTimeout: time.Duration(r.factory.Configs.AckDeadline * 1000),
	})
	return &r
}

func (r *redisHandler) RecieveCommand(binding ulti.ReceiveCommandFunction, filter ulti.FilterFunction, handler ulti.ReceiveHandlerFunction) handler.RecieveAndCreateHandler {
	topic := binding().Binding()
	topicHandler := handlerGroup{
		Receiver: handler,
		Filter:   filter,
	}
	r.subscriptions[topic] = topicHandler
	return r
}

func (r *redisHandler) RecieveEvent(binding ulti.ReceiveEventFunction, filter ulti.FilterFunction, handler ulti.ReceiveHandlerFunction) handler.RecieveAndCreateHandler {
	topic := binding().Topic()
	topicHandler := handlerGroup{
		Receiver: handler,
		Filter:   filter,
	}
	r.subscriptions[topic] = topicHandler
	return r
}

func (r *redisHandler) Start(ctx context.Context) error {

	fmt.Println("Starting")
	r.mu.Lock()
	ctx, r.cancelFunction = context.WithCancel(ctx)

	if r.client == nil {
		return ERR_INITIALIZE_CLIENT
	}

	_, err := r.client.Ping(ctx).Result()

	if err != nil {
		r.mu.Unlock()
		return ERR_CONNECT_SERVER
	}

	errChan := make(chan error, 0)

	defer close(errChan)
	defer r.client.Close()

	r.running = true
	for topic, handler := range r.subscriptions {
		sub := r.client.Subscribe(ctx, topic)
		defer sub.Close()
		for i := 0; i <= r.factory.Configs.MaxConcurrency; i++ {
			go func() {
				for {
					msg, err := sub.ReceiveMessage(ctx)
					if err != nil {
						errChan <- ERR_SUB_GENERAL
						return
					}
					handler.Filter()
					err = handler.Receiver(ctx, &model2.MessageEnvelope{
						ID:         uuid.New(),
						Data:       []byte(msg.Payload),
						Timestamp:  time.Now(),
						Attributes: nil,
					}, make(chan model2.DomainEvent))
					if err != nil {
						errChan <- err
						return
					}
				}
			}()
		}
	}
	r.mu.Unlock()
	defer fmt.Println("Shutting down")
	for {
		select {
		case <-ctx.Done():
			r.running = false
			return nil
		case err = <-errChan:
			if err != nil {
				r.running = false
				return err
			}
		}
	}
}

func (r *redisHandler) Stop() error {
	fmt.Println("Canceling")
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.running {
		r.cancelFunction()
		return nil
	}

	return ERR_SUB_WAITING
}
