# **onqlave-challenge**

Thanks for taking the time to solve this challenge. 

# Requirement

A piece of design is handed over as following. You should firstly understand what the design is and then develop a REDIS implementation of the pub/sub according to the design. The aim of design to keep it flexible enough to use the same skeleton for different underlying infrastrcutures  

```go
import (
	"context"
	"errors"
	"time"
)

type DomainEvent interface {
	Topic() string
	Metadata() map[string]string
}

type DomainCommand interface {
	Binding() string
	Metadata() map[string]string
}

type HandlerFactory interface {
	WithREDIS(configs REDISHandlerFactoryConfigs) WithREDISHandler
}

type WithREDISHandler interface {
	Subscriber(subscriberName string) RecieveAndCreateHandler
	Publisher() RecieveAndCreateHandler
}

type ReceiveHandlerFunction func(ctx context.Context, message *HandlerMessageEnvelope, eventStream chan<- DomainEvent) error
type ReceiveCommandFunction func() DomainCommand
type ReceiveEventFunction func() DomainEvent
type FilterFunction func() string

type RecieveAndCreateHandler interface {
	RecieveCommand(binding ReceiveCommandFunction, filter FilterFunction, handler ReceiveHandlerFunction) RecieveAndCreateHandler
	RecieveEvent(binding ReceiveEventFunction, filter FilterFunction, handler ReceiveHandlerFunction) RecieveAndCreateHandler
	Create() Handler
}

type HandlerMessageEnvelope struct {
	ID         string
	Data       []byte
	Timestamp  time.Time
	Attributes map[string]string
}

type Configurations struct {
}

const DEFAULT_CONCURRENCY int = 1

type Handler interface {
	Start(ctx context.Context) error
	PublishCommand(ctx context.Context, command DomainCommand) error
	PublishEvent(ctx context.Context, event DomainEvent) error
	Stop() error
}

type REDISHandlerFactoryConfigs struct {
	//include whatever else is required here for redis
	AckDeadline    int
	MaxConcurrency int
}

type handlerFactory struct {
}

func NewHandlerFactory() HandlerFactory {
	return &handlerFactory{}
}

type handlerGroup struct {
	Receiver ReceiveHandlerFunction
	Filter   FilterFunction
}

type redisHandlerFactory struct {
	configs        REDISHandlerFactoryConfigs
	subscriptions  map[string]handlerGroup
	topics         map[string]any //replace any with whatever type required for redis
	subscriberName string
}

type redisHandler struct {
	factory        *redisHandlerFactory
	client         *any
	cancelFunction context.CancelFunc
}

func (s *handlerFactory) WithREDIS(configs REDISHandlerFactoryConfigs) WithREDISHandler {
	if configs.MaxConcurrency == 0 {
		configs.MaxConcurrency = DEFAULT_CONCURRENCY
	}
	return &redisHandlerFactory{configs: configs}
}

func (s *redisHandlerFactory) Subscriber(subscriberName string) RecieveAndCreateHandler {
	//add logic here!
	return s
}

func (s *redisHandlerFactory) Publisher() RecieveAndCreateHandler {
	//add logic here!
	return s
}

func (s *redisHandlerFactory) RecieveCommand(binding ReceiveCommandFunction, filter FilterFunction, handler ReceiveHandlerFunction) RecieveAndCreateHandler {
	//add logic here!
	return s
}

func (s *redisHandlerFactory) RecieveEvent(binding ReceiveEventFunction, filter FilterFunction, handler ReceiveHandlerFunction) RecieveAndCreateHandler {
	//add logic here!
	return s
}

func (s *redisHandlerFactory) Create() Handler {
	return &redisHandler{factory: s}
}

func (s *redisHandler) Start(ctx context.Context) error {
	return errors.New("not implemented")
}

func (s *redisHandler) PublishCommand(ctx context.Context, command DomainCommand) error {
	return errors.New("not implemented")
}

func (s *redisHandler) PublishEvent(ctx context.Context, event DomainEvent) error {
	return errors.New("not implemented")
}

func (s *redisHandler) Stop() error {
	return errors.New("not implemented")
}

type TestCommand struct {
}

func (cmd *TestCommand) Binding() string {
	return "TestCommand"
}

func (cmd *TestCommand) Metadata() map[string]string {
	return nil
}

//usage 
func main() {
    //to receive and react to commands
	command := TestCommand{}
	ctx := context.Background()
	factory := NewHandlerFactory()
	handler := factory.
		WithREDIS(REDISHandlerFactoryConfigs{ /* add whatever configuration is required for REDIS */ }).
		Subscriber("service").
		RecieveCommand(func() DomainCommand { return &command }, func() string { return "" }, func(ctx context.Context, message *HandlerMessageEnvelope, eventStream chan<- DomainEvent) error {
			//we handle the commands here and do whatever is required with the command here
            return nil
		}).Create()
	err := handler.Start(ctx)
	if err != nil {
		panic("error in starting the command handler")
	}
}

```

# Outcome

What expteced as the outcome of this challenge is:

- Firstly explain what problem the design is trying to address
- Branch out from main
- Setup proper go module initialisation
- Structure your code based on what you believe would be the best practice - please indicate the reason
- Any improvement in design is appreciated - please indicate the reason
- Provide a package which can be used by other packages & modules in another golang application (can be imported by another module to be consumed)
- Implement enough test (to achieve at least 65% code coverage)


# Clarification

You can reach out directly on my [personal slack](https://join.slack.com/t/mposeidon/shared_invite/zt-1cfu52y47-koR3gkiXi3yTk_11sizNgQ) in case you need further clarification

## Note

Considering design is already provider, I would not expect the REDIS implementation to take you more than 3-4 hours. In case you need more time or there is something stopping you from progressing please reachout as above.

## License
Even though is it MIT licence, please keep it confidential!

[MIT License](./LICENSE)
