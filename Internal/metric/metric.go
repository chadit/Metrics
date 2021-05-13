package metric

import (
	"errors"
	"sync"
	"time"
)

// ErrNotFound is returned when the key is not found in the cache.
var ErrNotFound = errors.New("Metric key not found")

// Request represents the structure for a new request
type Request struct {
	Value int `json:"value"`
}

// Response represents the structure for a new response
type Response struct {
	Value int `json:"value"`
}

type metricCollection struct {
	cache map[string]map[time.Time]int
	lock  sync.Mutex
}

// Collection represents the rcv for metrics handling.
type Collection metricCollection

// New creates a new metric collection.
func New() *Collection {
	return &Collection{
		cache: make(map[string]map[time.Time]int),
	}
}

// AddByTime will add a new metric entry.
func (c *Collection) AddByTime(metricKey string, newDate time.Time, data int) error {

	// if key exist already
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, exist := c.cache[metricKey]; exist {
		c.cache[metricKey][newDate] = data
		return nil
	}

	c.cache[metricKey] = map[time.Time]int{newDate: data}
	return nil
}

// Add will add a new metric entry.
func (c *Collection) Add(metricKey string, data int) error {
	newDate := time.Now().UTC()
	return c.AddByTime(metricKey, newDate, data)
}

// Sum returns a sum of metrics over x duration.
func (c *Collection) Sum(metricKey string, duration time.Duration) (int, error) {
	var sum int

	c.lock.Lock()
	defer c.lock.Unlock()

	if _, exist := c.cache[metricKey]; !exist {
		return sum, ErrNotFound
	}

	endTime := time.Now().UTC()
	startTime := endTime.Add(-duration)

	for k, v := range c.cache[metricKey] {
		if k.After(startTime) && k.Before(endTime) {
			sum += v
		}
	}

	return sum, nil
}

// Purge removes metrics over x duration.
func (c *Collection) Purge(duration time.Duration) error {

	c.lock.Lock()
	defer c.lock.Unlock()

	startTime := time.Now().UTC().Add(-duration)

	for metricKey := range c.cache {
		for k := range c.cache[metricKey] {
			if k.Before(startTime) {
				delete(c.cache[metricKey], k)
			}
		}
	}

	return nil
}
