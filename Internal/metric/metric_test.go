package metric_test

import (
	"testing"
	"time"

	"github.com/chadit/Metrics/Internal/metric"
	"github.com/stretchr/testify/assert"
)

func Test_BasicSum(t *testing.T) {
	newCollection := metric.New()
	key := "test"

	newMectrics := []int{1, 2, 3, 4, 5}

	for _, v := range newMectrics {
		newCollection.Add(key, v)
	}

	actual, err := newCollection.Sum(key, 1*time.Minute)

	assert.Nil(t, err)
	assert.Equal(t, 15, actual)
}

func Test_DurationSum(t *testing.T) {
	newCollection := metric.New()
	key := "test"

	newMectrics := []int{1, 2, 3, 4, 5}

	for i, v := range newMectrics {
		dur := time.Duration(i * int(time.Minute))
		eventDate := time.Now().UTC().Add(-dur)
		newCollection.AddByTime(key, eventDate, v)
	}

	actual, err := newCollection.Sum(key, 2*time.Minute)

	assert.Nil(t, err)
	assert.Equal(t, 3, actual)
}

func Test_Purge(t *testing.T) {
	newCollection := metric.New()
	key := "test"

	newMectrics := []int{1, 2, 3, 4, 5}

	for i, v := range newMectrics {
		dur := time.Duration((i * 10) * int(time.Minute))
		eventDate := time.Now().UTC().Add(-dur)
		newCollection.AddByTime(key, eventDate, v)
	}

	newCollection.Purge(time.Duration(25 * int(time.Minute)))

	actual, err := newCollection.Sum(key, 590*time.Minute)

	assert.Nil(t, err)
	assert.Equal(t, 6, actual)
}
