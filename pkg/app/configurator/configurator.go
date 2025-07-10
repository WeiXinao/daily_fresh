package configurator

import (
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/WeiXinao/daily_fresh/pkg/app/configurator/subscriber"
	"github.com/WeiXinao/daily_fresh/pkg/log"
)

var (
	errEmptyConfig            = errors.New("empty config value")
	errMissingUnmarshalerType = errors.New("missing unmarshaler type")
)

// Configurator is the interface for configuration center.
type Configurator[T any] interface {
	// GetConfig returns the subscription value.
	GetConfig() (T, error)
	GetConfigString() (string, error)
	// AddListener adds a listener to the subscriber.
	AddListener(listener func(key, raw string, data T))
}

type (
	// Config is the configuration for Configurator.
	Config struct {
		// Type is the value type, yaml, json or toml.
		Type string `json:",default=yaml,options=[yaml,json,toml]"`
		// Log is the flag to control logging.
		Log bool `json:",default=true"`
	}

	configCenter[T any] struct {
		conf        Config
		unmarshaler LoaderFn
		subscriber  subscriber.Subscriber
		listeners   []func(key, raw string, data T)
		lock        sync.Mutex
		snapshot    atomic.Value
	}

	value[T any] struct {
		key 				string
		data        string
		marshalData T
		err         error
	}
)

func must(err error) {
	if err == nil {
		return
	}

	msg := fmt.Sprintf("%+v\n\n%s", err.Error(), debug.Stack())
	log.Error(msg)
	panic(msg)
}

// Configurator is the interface for configuration center.
var _ Configurator[any] = (*configCenter[any])(nil)

// MustNewConfigCenter returns a Configurator, exits on errors.
func MustNewConfigCenter[T any](c Config, subscriber subscriber.Subscriber) Configurator[T] {
	cc, err := NewConfigCenter[T](c, subscriber)
	must(err)
	return cc
}

// NewConfigCenter returns a Configurator.
func NewConfigCenter[T any](c Config, subscriber subscriber.Subscriber) (Configurator[T], error) {
	unmarshaler, ok := Unmarshaler(strings.ToLower(c.Type))
	if !ok {
		return nil, fmt.Errorf("unknown format: %s", c.Type)
	}

	cc := &configCenter[T]{
		conf:        c,
		unmarshaler: unmarshaler,
		subscriber:  subscriber,
	}

	s, err := cc.subscriber.Value(); 
	if err != nil {
		if cc.conf.Log {
			log.Errorf("ConfigCenter loads configuration, error: %v", err)
		}
		return nil, err
	}

	if err := cc.loadConfig(cc.subscriber.Key(),s); err != nil {
		return nil, err
	}

	if err := cc.subscriber.AddListener(cc.onChange); err != nil {
		return nil, err
	}

	if _, err := cc.GetConfig(); err != nil {
		return nil, err
	}

	return cc, nil
}

// AddListener adds listener to s.
func (c *configCenter[T]) AddListener(listener func(key, row string, data T)) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.listeners = append(c.listeners, listener)
}

// GetConfig return structured config.
func (c *configCenter[T]) GetConfig() (T, error) {
	v := c.value()
	if v == nil || len(v.data) == 0 {
		var empty T
		return empty, errEmptyConfig
	}

	return v.marshalData, v.err
}

func (c *configCenter[T]) GetConfigString() (string, error) {
	v := c.value()
	if v == nil || len(v.data) == 0 {
		return "", errEmptyConfig
	}

	return v.data, v.err
}

// Value returns the subscription value.
func (c *configCenter[T]) Value() string {
	v := c.value()
	if v == nil {
		return ""
	}
	return v.data
}

func (c *configCenter[T]) loadConfig(key, data string) error {
	if c.conf.Log {
		log.Infof("ConfigCenter loads changed configuration, key [%s], data [%s]", key, data)
	}

	c.snapshot.Store(c.genValue(key, data))
	return nil
}

func (c *configCenter[T]) onChange(key, data string) {
	if err := c.loadConfig(key, data); err != nil {
		return
	}

	c.lock.Lock()
	listeners := make([]func(key, raw string, data T), len(c.listeners))
	copy(listeners, c.listeners)
	c.lock.Unlock()

	for _, l := range listeners {
		go l(key, c.value().data, c.value().marshalData)
	}
}

func (c *configCenter[T]) value() *value[T] {
	content := c.snapshot.Load()
	if content == nil {
		return nil
	}
	return content.(*value[T])
}

// deref dereferences a type, if pointer type, returns its element type.
func deref(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t
}

func (c *configCenter[T]) genValue(key, data string) *value[T] {
	v := &value[T]{
		key: key,
		data: data,
	}
	if len(data) == 0 {
		return v
	}

	t := reflect.TypeOf(v.marshalData)
	// if the type is nil, it means that the user has not set the type of the configuration.
	if t == nil {
		v.err = errMissingUnmarshalerType
		return v
	}

	t = deref(t)
	switch t.Kind() {
	case reflect.Struct, reflect.Array, reflect.Slice:
		if err := c.unmarshaler([]byte(data), &v.marshalData); err != nil {
			v.err = err
			if c.conf.Log {
				log.Errorf("ConfigCenter unmarshal configuration failed, err: %+v, content [%s]",
					err.Error(), data)
			}
		}
	case reflect.String:
		if str, ok := any(data).(T); ok {
			v.marshalData = str
		} else {
			v.err = errMissingUnmarshalerType
		}
	default:
		if c.conf.Log {
			log.Errorf("ConfigCenter unmarshal configuration missing unmarshaler for type: %s, content [%s]",
				t.Kind(), data)
		}
		v.err = errMissingUnmarshalerType
	}

	return v
}
