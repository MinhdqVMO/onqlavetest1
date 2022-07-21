package handler

import (
	"context"
	"github.com/onqlavetest1/model"
	"github.com/onqlavetest1/ulti"
)

type Handler interface {
	Start(ctx context.Context) error
	Stop() error
	PublishCommand(ctx context.Context, command model.DomainCommand) error
	PublishEvent(ctx context.Context, event model.DomainEvent) error
}

type CreateHandler interface {
	Create() Handler
}

type RecieveAndCreateHandler interface {
	CreateHandler
	RecieveCommand(binding ulti.ReceiveCommandFunction, filter ulti.FilterFunction, handler ulti.ReceiveHandlerFunction) RecieveAndCreateHandler
	RecieveEvent(binding ulti.ReceiveEventFunction, filter ulti.FilterFunction, handler ulti.ReceiveHandlerFunction) RecieveAndCreateHandler
}
