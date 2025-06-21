package api

import (
	"sync"
	"time"
)

type Entry struct {
	Time  time.Time
	Value interface{}
	Err   error
}

var mu = sync.Mutex{}
var data = map[string]Entry{}

func LazyCache[T any](key string, t time.Duration, callback func() (T, error)) (T, bool, error) {
	mu.Lock()
	defer mu.Unlock()

	entry, found := data[key]

	if found {
		if time.Since(entry.Time) > t {
			now := time.Now()
			data[key] = Entry{
				Time:  now,
				Value: entry.Value,
			}
			go func() {
				new_value, err := callback()

				mu.Lock()
				defer mu.Unlock()

				if err != nil {
					data[key] = Entry{
						Time:  entry.Time,
						Value: entry.Value,
						Err:   err,
					}
				} else {
					data[key] = Entry{
						Time:  data[key].Time,
						Value: new_value,
					}
				}
			}()
		}

		err := entry.Err
		if entry.Err != nil {
			data[key] = Entry{
				Time:  entry.Time,
				Value: entry.Value,
			}
		}
		return entry.Value.(T), true, err
	}

	new_value, err := callback()

	if err != nil {
		return new_value, false, err
	}
	data[key] = Entry{
		Time:  time.Now(),
		Value: new_value,
	}

	return new_value, true, nil
}
