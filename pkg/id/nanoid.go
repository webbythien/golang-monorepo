package id

import (
	"github.com/google/uuid"
	nanoid "github.com/matoous/go-nanoid/v2"
)

const (
	DefaultAlphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
	DefaultLength   = 8 // -> max combinations: 36^8 = 2,821,109,907,456, collision 1% at 2,821,109,907,456 * 0.01 = 28,211,099,074
	UpperAlphabet   = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Base33Alphabet  = "0123456789ABCDEFGHJKLMNPRSTUVWXYZ" // uppercase, without I, O, Q (easy for mistaking)
)

type Generator interface {
	NewID(params ...interface{}) (string, error)
}

type nanoIDConfig struct {
	alphabet string
	length   int
	prefix   string
}

type NanoIDOption func(*nanoIDConfig)

func WithAlphabet(alphabet string) NanoIDOption {
	return func(c *nanoIDConfig) {
		c.alphabet = alphabet
	}
}

func WithUpperAlphabet() NanoIDOption {
	return func(c *nanoIDConfig) {
		c.alphabet = UpperAlphabet
	}
}

func WithBase33Alphabet() NanoIDOption {
	return func(c *nanoIDConfig) {
		c.alphabet = Base33Alphabet
	}
}

func WithLength(length int) NanoIDOption {
	return func(c *nanoIDConfig) {
		c.length = length
	}
}

func WithPrefix(prefix string) NanoIDOption {
	return func(c *nanoIDConfig) {
		c.prefix = prefix
	}
}

type nanoIDGenerator struct {
	nanoIDConfig
}

func NewGenerator(opts ...NanoIDOption) Generator {
	cfg := nanoIDConfig{
		alphabet: DefaultAlphabet,
		length:   DefaultLength,
		prefix:   "",
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	return &nanoIDGenerator{
		nanoIDConfig: cfg,
	}
}

func (s *nanoIDGenerator) NewID(_ ...interface{}) (string, error) {
	targetID, err := nanoid.Generate(s.alphabet, s.length)
	if err != nil {
		return "", err
	}
	if s.prefix != "" {
		return s.prefix + targetID, nil
	}
	return targetID, nil
}

type uuidGenerator struct{}

func NewUUIDGenerator() Generator {
	return &uuidGenerator{}
}

func (s *uuidGenerator) NewID(_ ...interface{}) (string, error) {
	return uuid.NewString(), nil
}
