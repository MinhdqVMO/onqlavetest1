
# onqlave-challenge






## Requirement

- Firstly explain what problem the design is trying to address
- Branch out from main - test*
- Setup proper go module initialisation
- Structure your code based on what you believe would be the best practice - please indicate the reason
- Any improvement in design is appreciated - please indicate the reason
- Provide a package which can be used by other packages & modules in another golang application (can be imported by another module to be consumed)
- Implement enough test (to achieve at least 65% code coverage)


### 


#### The Problem

Based on the interface and the description, I figure you want a dynamic and flexible design that easily to scale or implement new technology.

- The main idea of the design based on interfaces implementation and try to be as abstract as possible.

- You want to implement a main factory interface that produces smaller interfaces like (Redis factory, ... ) - which in this case is a redis pub sub manager and factory.

- You also want to customize the payload (as long as they are the implementation of the DomainEvent or DomainCommand) and be able to change the logic function.

The design is "plug and play".

- The user doesn't need to know the code or the technology that we implemented

For example if users want a redis pubsub manager ,they just need to import the package, define the payload , set up the handle function.

They also can define their own handler or even a different database/cache technology base on the existed interfaces.


### Changes

I decided to make some change to the design. If you're disagree with any changes please tell me.

First i want only one instance of main factory is created. Because its purpose is to create smaller factory so i think we only need one of it - which is more memory efficient, easier to manage and debug


```go
var factory *handlerFactory

type handlerFactory struct {
}

type HandlerFactory interface {
	WithREDIS(configs config.REDISHandlerFactoryConfigs) WithREDISHandler
}

func InitHandlerFactory() {
	if factory != nil {
		return
	}
	factory = &handlerFactory{}
}

func GetFactory() *handlerFactory {
	return factory
}

```

I also added a CreateHandler interface


```go

type CreateHandler interface {
	Create() Handler
}

type RecieveAndCreateHandler interface {
	CreateHandler
	RecieveCommand(binding ulti.ReceiveCommandFunction, filter ulti.FilterFunction, handler ulti.ReceiveHandlerFunction) RecieveAndCreateHandler
	RecieveEvent(binding ulti.ReceiveEventFunction, filter ulti.FilterFunction, handler ulti.ReceiveHandlerFunction) RecieveAndCreateHandler
}

```

for the Publisher method because it make more sense to me that Publisher don't need fucntion RecieveCommand and  RecieveEvent


```go

type WithREDISHandler interface {
	Subscriber(subscriberName string) handler.RecieveAndCreateHandler
	Publisher() handler.CreateHandler
}

```

And I changed the the RedisHandlerFactory and redisHandler struct so it behave like like a factory and a manager at the same time


```go

type RedisHandlerFactory struct {
	Configs        config.REDISHandlerFactoryConfigs
	subscribers    map[string]*redisHandler
	subscriberLock *sync.Mutex
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

```

### The Structure

```bash

#Project Structure

├── cmd
│   └── pubsub
│       └── main_test.go
├── config
│   └── config.go
├── configs
├── go.mod
├── go.sum
├── LICENSE
├── model
│   ├── command.go
│   └── envelope.go
├── persistence
│   ├── factory.go
│   ├── handler
│   │   └── handler.go
│   ├── redisFactory
│   │   ├── publish_test.go
│   │   ├── redisFactory.go
│   │   └── subscribe_test.go
│   └── redisHandler.go
├── README.md
├── test
│   ├── subscribe_test.go
│   └── test.go
└── ulti
    └── alias.go



```

The structure i came up with based on the Abstract Factory Pattern and SOLID principle. I splitted the interface and struct base on the layer it's working on and its purpose.

The upper layer interface import lower layer interface with abstract interface on the upper layer and more specific inteface at the lower layer

For example:

- The factory interface - which is abstract - is used to create redisHandler that implemented by RedisHandlerFactory - which is more specific.

- The handler interface is used in specific factory so it should be at the same layer as the factory

This structure also very clear to read and easy to debug. It's also straightforward when implementing different technologies.

I also decided to not use internal package so the user can easily imports and customizes the code.


### Testing

```bash

    //main.test
     cd cmd/pubsub/
     go test

     //publish_test 

    publish_test.go

    //subscribe_test

    subscribe_test.go
    

```


## Contact

If you have any question or something you want to discuss about please contact me through my gmail: minhdq@vmodev.com

or my github account https://github.com/MinhdqVMO.



I apology for any grammar mistakes or misspelled words and also sorry for taking so long :(