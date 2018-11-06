package main

// Action is interface for the actions
type Action interface {
	Name() string
	Do() (error, error)
}
