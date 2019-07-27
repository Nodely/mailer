package processor

import (
	"encoding/json"
	"errors"
)

// Processor struct
type Processor struct {
	Resolver TemplateResolver
}

type dataInput struct {
	TemplateID string `json:"template"`
	Locale     string `json:"locale"`
}

// New init processor instance
func New() *Processor {
	return &Processor{
		Resolver: &defaultResolver{},
	}
}

// Do prepares template for sending including data validation
func (a *Processor) Do(payload []byte) error {
	if payload == nil {
		return errors.New("payload is nil")
	}

	// parse input
	var input dataInput
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}

	// verify locale
	a.Resolver.Locate("some", "en_US")

	return nil
}
