package callback

import (
	"errors"
	"fmt"
	"strings"
)

// CallbackData - интерфейс для всех типов callback-ов
type CallbackData interface {
	Type() string              // Тип callback (префикс)
	Validate() error           // Валидация данных
	Marshal() ([]byte, error)  // Сериализация в bytes
	Unmarshal([]byte) error    // Десериализация из bytes
}

// Callback - обёртка для работы с callback предоставляющая методы Pack Unpack
type Callback struct {
	TypeStr string `json:"type"`
	Data    []byte `json:"data"`
}

// Pack - упаковывает callback в строку для Telegram
func (c *Callback) pack() (string, error) {
	combined := fmt.Sprintf("%s_%s", c.TypeStr, string(c.Data))
	
	if len(combined) > 64 {
		return "", errors.New("callback data exceeds 64 bytes limit")
	}
	
	return combined, nil
}

func PackCallback(data CallbackData) (string, error) {
	serilized, err := data.Marshal()
	if err != nil {
		return "", fmt.Errorf("marshal failed %w", err)
	}

	cb := &Callback{
		TypeStr: data.Type(),
		Data: serilized,
	}
	return cb.pack()
}

// Unpack - распаковывает строку из Telegram в Callback
func Unpack(callbackStr string) (*Callback, error) {
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
	types map[string]func() CallbackData
}

// NewRegistry создает новый реестр
func NewRegistry() *Registry {
	return &Registry{
		types: make(map[string]func() CallbackData),
	}
}

// Register регистрирует новый тип callback
func (r *Registry) Register(ctor func() CallbackData) {
	cb := ctor()
	r.types[cb.Type()] = ctor
}

// Parse преобразует строку callback в конкретный тип
func (r *Registry) Parse(callbackStr string) (CallbackData, error) {
	cb, err := Unpack(callbackStr)
	if err != nil {
		return nil, err
	}
	
	ctor, exists := r.types[cb.TypeStr]
	if !exists {
		return nil, fmt.Errorf("unknown callback type: %s", cb.TypeStr)
	}
	
	instance := ctor()
	if err := instance.Unmarshal(cb.Data); err != nil {
		return nil, fmt.Errorf("unmarshal failed: %w", err)
	}
	
	if err := instance.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	
	return instance, nil
}