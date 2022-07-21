package ulti

import (
	"context"
	model2 "github.com/onqlavetest1/model"
)

type ReceiveHandlerFunction func(ctx context.Context, message *model2.MessageEnvelope, eventStream chan<- model2.DomainEvent) error
type ReceiveCommandFunction func() model2.DomainCommand
type ReceiveEventFunction func() model2.DomainEvent
type FilterFunction func() string
