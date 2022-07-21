package model

type DomainEvent interface {
	Topic() string
	Metadata() map[string]string
}

type DomainCommand interface {
	Binding() string
	Metadata() map[string]string
}
