package callback

import (
	"encoding/json"
	"errors"
	"fmt"
	"shopping_bot/pkg/logger"
	"strings"
)

// CallbackService - интерфейс для всех типов callback-ов
type CallbackService interface {
	Type() string             // Тип callback (префикс)
	Validate() error          // Валидация данных
}

// Callback - обёртка для работы с callback предоставляющая методы Pack Unpack
type Callback struct {
	TypeStr string `json:"type"`
	Data    []byte `json:"data"`
}

// pack - упаковывает callback в строку для Telegram
func (c *Callback) pack() (string, error) {
	combined := fmt.Sprintf("%s_%s", c.TypeStr, string(c.Data))

	if len(combined) > 64 {
		logger.Sugar.Errorw("combined callback", "combined", combined)
		return "", errors.New("callback data exceeds 64 bytes limit")
	}

	return combined, nil
}

// unpack - распаковывает строку из Telegram в Callback
func unpack(callbackStr string) (*Callback, error) {
	parts := strings.SplitN(callbackStr, "_", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid callback format")
	}

	return &Callback{
		TypeStr: parts[0],
		Data:    []byte(parts[1]),
	}, nil
}

// Registry - реестр для регистрации типов callback
type Registry struct {
	types map[string]func() CallbackService
}

// NewRegistry создает новый реестр
func NewRegistry() *Registry {
	return &Registry{
		types: make(map[string]func() CallbackService),
	}
}

// Register регистрирует новый тип callback
func (r *Registry) Register(ctor func() CallbackService) {
	cb := ctor()
	r.types[cb.Type()] = ctor
}

// parse преобразует строку callback в конкретный тип
func (r *Registry) parse(callbackStr string) (CallbackService, error) {
	cb, err := unpack(callbackStr)
	if err != nil {
		return nil, err
	}

	ctor, exists := r.types[cb.TypeStr]
	if !exists {
		return nil, fmt.Errorf("unknown callback type: %s", cb.TypeStr)
	}

	instance := ctor()
	dataBytes := []byte(cb.Data)
	if err := json.Unmarshal(dataBytes, instance); err != nil {
		return nil, fmt.Errorf("unmarshal failed: %w", err)
	}

	if err := instance.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return instance, nil
}

// Pack Callback and serialize to string
func PackCallback(data CallbackService) (string, error) {
	serilized, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("marshal failed %w", err)
	}

	cb := &Callback{
		TypeStr: data.Type(),
		Data:    serilized,
	}
	return cb.pack()
}

// Parse Callback and deserilize to 
func ParseCallback[T any](registry *Registry, data string) (T, error) {
	var zero T
	callbackData, err := registry.parse(data)
	if err != nil {
		return zero, fmt.Errorf("parse callback: %w", err)
	}

	result, ok := callbackData.(T)
	if !ok {
		return zero, fmt.Errorf("unexpected callback type: %T", callbackData)
	}

	return result, nil
}
