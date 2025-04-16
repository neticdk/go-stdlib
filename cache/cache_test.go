package cache

import (
	"reflect"
	"testing"
	"time"
)

func withTestClock(c *mockClock) Option {
	return func(cache *cache) {
		cache.clock = c
	}
}

func TestCache_GetAllItems(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		items       map[string]Item
		want        map[string]Item
	}{
		{
			name:        "No items",
			currentTime: 9,
			items:       map[string]Item{},
			want:        map[string]Item{},
		},
		{
			name:        "With items",
			currentTime: 9,
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
				"key2": {Value: "value2", Expiration: time.Unix(15, 0)},
			},
			want: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
				"key2": {Value: "value2", Expiration: time.Unix(15, 0)},
			},
		},
		{
			name:        "With expired items",
			currentTime: 11,
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
				"key2": {Value: "value2", Expiration: time.Unix(5, 0)},
			},
			want: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
				"key2": {Value: "value2", Expiration: time.Unix(5, 0)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCache(withTestClock(newMockClock(tt.currentTime)), WithItems(tt.items))
			got := c.GetAllItems()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_GetItems(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		items       map[string]Item
		want        map[string]Item
	}{
		{
			name:        "No items",
			currentTime: 9,
			items:       map[string]Item{},
			want:        map[string]Item{},
		},
		{
			name:        "With items",
			currentTime: 9,
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
				"key2": {Value: "value2", Expiration: time.Unix(15, 0)},
			},
			want: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
				"key2": {Value: "value2", Expiration: time.Unix(15, 0)},
			},
		},
		{
			name:        "With expired items",
			currentTime: 11,
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
				"key2": {Value: "value2", Expiration: time.Unix(15, 0)},
			},
			want: map[string]Item{
				"key2": {Value: "value2", Expiration: time.Unix(15, 0)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCache(withTestClock(newMockClock(tt.currentTime)), WithItems(tt.items))
			got := c.GetItems()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Get(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		key         string
		items       map[string]Item
		want        any
		wantFound   bool
	}{
		{
			name:        "Key exists and not expired",
			currentTime: 9,
			key:         "key1",
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want:      "value1",
			wantFound: true,
		},
		{
			name:        "Key exists but expired",
			currentTime: 11,
			key:         "key1",
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want:      nil,
			wantFound: false,
		},
		{
			name:        "Key does not exist",
			currentTime: 9,
			key:         "key2",
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want:      nil,
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCache(withTestClock(newMockClock(tt.currentTime)), WithItems(tt.items))
			got, gotFound := c.Get(tt.key)
			if got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
			if gotFound != tt.wantFound {
				t.Errorf("Get() found = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func TestCache_GetTTL(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		key         string
		items       map[string]Item
		want        time.Duration
		wantFound   bool
	}{
		{
			name:        "Key exists and not expired",
			currentTime: 9,
			key:         "key1",
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want:      1 * time.Second,
			wantFound: true,
		},
		{
			name:        "Key exists but expired",
			currentTime: 11,
			key:         "key1",
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want:      0,
			wantFound: false,
		},
		{
			name:        "Key does not exist",
			currentTime: 9,
			key:         "key2",
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want:      0,
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCache(withTestClock(newMockClock(tt.currentTime)), WithItems(tt.items))
			got, gotFound := c.GetTTL(tt.key)
			if got != tt.want {
				t.Errorf("GetTTL() = %v, want %v", got, tt.want)
			}
			if gotFound != tt.wantFound {
				t.Errorf("GetTTL() found = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func TestCache_Delete(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		key         string
		items       map[string]Item
		want        map[string]Item
	}{
		{
			name:        "Key exists and not expired",
			currentTime: 9,
			key:         "key1",
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want: map[string]Item{},
		},
		{
			name:        "Key exists but expired",
			currentTime: 11,
			key:         "key1",
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want: map[string]Item{},
		},
		{
			name:        "Key does not exist",
			currentTime: 9,
			key:         "key2",
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCache(withTestClock(newMockClock(tt.currentTime)), WithItems(tt.items))
			c.Delete(tt.key)
			got := c.GetAllItems()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Clear(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		items       map[string]Item
	}{
		{
			name:        "Existing items",
			currentTime: 9,
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
				"key2": {Value: "value2", Expiration: time.Unix(5, 0)},
			},
		},
		{
			name:        "No items",
			currentTime: 11,
			items:       map[string]Item{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCache(withTestClock(newMockClock(tt.currentTime)), WithItems(tt.items))
			c.Clear()
			got := c.GetAllItems()
			if len(got) != 0 {
				t.Errorf("Clear() = %v, want %v", got, map[string]Item{})
			}
		})
	}
}

func TestCache_DeleteExpired(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		items       map[string]Item
		want        map[string]Item
	}{
		{
			name:        "With both expired and non-expired items",
			currentTime: 10,
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(8, 0)},
				"key2": {Value: "value2", Expiration: time.Unix(9, 0)},
				"key3": {Value: "value3", Expiration: time.Unix(10, 0)},
			},
			want: map[string]Item{
				"key3": {Value: "value3", Expiration: time.Unix(10, 0)},
			},
		},
		{
			name:        "All expired",
			currentTime: 11,
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(8, 0)},
				"key2": {Value: "value2", Expiration: time.Unix(9, 0)},
				"key3": {Value: "value3", Expiration: time.Unix(10, 0)},
			},
			want: map[string]Item{},
		},
		{
			name:        "No expirations",
			currentTime: 5,
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(8, 0)},
				"key2": {Value: "value2", Expiration: time.Unix(9, 0)},
				"key3": {Value: "value3", Expiration: time.Unix(10, 0)},
			},
			want: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(8, 0)},
				"key2": {Value: "value2", Expiration: time.Unix(9, 0)},
				"key3": {Value: "value3", Expiration: time.Unix(10, 0)},
			},
		},
		{
			name:        "No items",
			currentTime: 11,
			items:       map[string]Item{},
			want:        map[string]Item{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCache(withTestClock(newMockClock(tt.currentTime)), WithItems(tt.items))
			c.DeleteExpired()
			got := c.GetAllItems()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Renew(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		opts        []Option
		key         string
		items       map[string]Item
		want        time.Duration
		wantFound   bool
	}{
		{
			name:        "Key exists and not expired",
			currentTime: 9,
			key:         "key1",
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want:      1 * time.Second,
			wantFound: true,
		},
		{
			name:        "Key exists but expired",
			currentTime: 11,
			key:         "key1",
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want:      1 * time.Second,
			wantFound: true,
		},
		{
			name:        "Key does not exist",
			currentTime: 9,
			key:         "key2",
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want:      0,
			wantFound: false,
		},
		{
			name:        "With custom cache interval",
			currentTime: 9,
			opts:        []Option{WithCacheInterval(2 * time.Second)},
			key:         "key1",
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want:      2 * time.Second,
			wantFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := []Option{withTestClock(newMockClock(tt.currentTime)), WithItems(tt.items)}
			opts = append(opts, tt.opts...)
			c := NewCache(opts...)
			c.Renew(tt.key)
			got, gotFound := c.GetTTL(tt.key)
			if got != tt.want {
				t.Errorf("Renew() = %v, want %v", got, tt.want)
			}
			if gotFound != tt.wantFound {
				t.Errorf("Renew() found = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func TestCache_RenewInterval(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		opts        []Option
		key         string
		renewalTime time.Duration
		items       map[string]Item
		want        time.Duration
		wantFound   bool
	}{
		{
			name:        "Key exists and not expired",
			currentTime: 9,
			key:         "key1",
			renewalTime: 5 * time.Second,
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want:      5 * time.Second,
			wantFound: true,
		},
		{
			name:        "Key exists but expired",
			currentTime: 11,
			key:         "key1",
			renewalTime: 5 * time.Second,
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want:      5 * time.Second,
			wantFound: true,
		},
		{
			name:        "Key does not exist",
			currentTime: 9,
			key:         "key2",
			renewalTime: 5 * time.Second,
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want:      0,
			wantFound: false,
		},
		{
			name:        "With custom cache interval",
			currentTime: 9,
			opts:        []Option{WithCacheInterval(2 * time.Second)},
			key:         "key1",
			renewalTime: 5 * time.Second,
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want:      5 * time.Second,
			wantFound: true,
		},
		{
			name:        "Negative renewal time",
			currentTime: 9,
			opts:        []Option{WithCacheInterval(2 * time.Second)},
			key:         "key1",
			renewalTime: -5 * time.Second,
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want:      1 * time.Second,
			wantFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := []Option{withTestClock(newMockClock(tt.currentTime)), WithItems(tt.items)}
			opts = append(opts, tt.opts...)
			c := NewCache(opts...)
			c.RenewInterval(tt.key, tt.renewalTime)
			got, gotFound := c.GetTTL(tt.key)
			if got != tt.want {
				t.Errorf("RenewInterval() = %v, want %v", got, tt.want)
			}
			if gotFound != tt.wantFound {
				t.Errorf("RenewInterval() found = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func TestCache_SetInterval(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		opts        []Option
		key         string
		value       any
		interval    time.Duration
		items       map[string]Item
		want        map[string]Item
	}{
		{
			name:        "Key exists and not expired",
			currentTime: 9,
			key:         "key1",
			value:       "newvalue1",
			interval:    5 * time.Second,
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want: map[string]Item{
				"key1": {Value: "newvalue1", Expiration: time.Unix(14, 0)},
			},
		},
		{
			name:        "Key exists but expired",
			currentTime: 9,
			key:         "key1",
			value:       "newvalue1",
			interval:    5 * time.Second,
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(8, 0)},
			},
			want: map[string]Item{
				"key1": {Value: "newvalue1", Expiration: time.Unix(14, 0)},
			},
		},
		{
			name:        "Key does not exists",
			currentTime: 9,
			key:         "key2",
			value:       "value2",
			interval:    5 * time.Second,
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
				"key2": {Value: "value2", Expiration: time.Unix(14, 0)},
			},
		},
		{
			name:        "No items",
			currentTime: 9,
			key:         "key2",
			value:       "value2",
			interval:    5 * time.Second,
			items:       map[string]Item{},
			want: map[string]Item{
				"key2": {Value: "value2", Expiration: time.Unix(14, 0)},
			},
		},
		{
			name:        "Negative interval",
			currentTime: 9,
			key:         "key2",
			value:       "value2",
			interval:    -5 * time.Second,
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
				"key2": {Value: "value2", Expiration: time.Unix(9, 0)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := []Option{withTestClock(newMockClock(tt.currentTime)), WithItems(tt.items)}
			opts = append(opts, tt.opts...)
			c := NewCache(opts...)
			c.SetInterval(tt.key, tt.value, tt.interval)
			got := c.GetAllItems()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Set(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		opts        []Option
		key         string
		value       any
		items       map[string]Item
		want        map[string]Item
	}{
		{
			name:        "Key exists and not expired",
			currentTime: 9,
			key:         "key1",
			value:       "newvalue1",
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want: map[string]Item{
				"key1": {Value: "newvalue1", Expiration: time.Unix(10, 0)},
			},
		},
		{
			name:        "Key exists but expired",
			currentTime: 9,
			key:         "key1",
			value:       "newvalue1",
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(8, 0)},
			},
			want: map[string]Item{
				"key1": {Value: "newvalue1", Expiration: time.Unix(10, 0)},
			},
		},
		{
			name:        "Key does not exists",
			currentTime: 9,
			key:         "key2",
			value:       "value2",
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
				"key2": {Value: "value2", Expiration: time.Unix(10, 0)},
			},
		},
		{
			name:        "No items",
			currentTime: 9,
			key:         "key2",
			value:       "value2",
			items:       map[string]Item{},
			want: map[string]Item{
				"key2": {Value: "value2", Expiration: time.Unix(10, 0)},
			},
		},
		{
			name:        "Custom interval",
			currentTime: 9,
			key:         "key2",
			value:       "value2",
			opts:        []Option{WithCacheInterval(5 * time.Second)},
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
				"key2": {Value: "value2", Expiration: time.Unix(14, 0)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := []Option{withTestClock(newMockClock(tt.currentTime)), WithItems(tt.items)}
			opts = append(opts, tt.opts...)
			c := NewCache(opts...)
			c.Set(tt.key, tt.value)
			got := c.GetAllItems()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_AddInterval(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		opts        []Option
		key         string
		value       any
		interval    time.Duration
		items       map[string]Item
		want        map[string]Item
	}{
		{
			name:        "Key exists and not expired",
			currentTime: 9,
			key:         "key1",
			value:       "newvalue1",
			interval:    5 * time.Second,
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
		},
		{
			name:        "Key exists but allow overrides",
			currentTime: 9,
			key:         "key1",
			value:       "newvalue1",
			interval:    5 * time.Second,
			opts:        []Option{WithAllowOverwrites()},
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want: map[string]Item{
				"key1": {Value: "newvalue1", Expiration: time.Unix(14, 0)},
			},
		},
		{
			name:        "Key exists but expired",
			currentTime: 9,
			key:         "key1",
			value:       "newvalue1",
			interval:    5 * time.Second,
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(8, 0)},
			},
			want: map[string]Item{
				"key1": {Value: "newvalue1", Expiration: time.Unix(14, 0)},
			},
		},
		{
			name:        "Key does not exists",
			currentTime: 9,
			key:         "key2",
			value:       "value2",
			interval:    5 * time.Second,
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
				"key2": {Value: "value2", Expiration: time.Unix(14, 0)},
			},
		},
		{
			name:        "No items",
			currentTime: 9,
			key:         "key2",
			value:       "value2",
			interval:    5 * time.Second,
			items:       map[string]Item{},
			want: map[string]Item{
				"key2": {Value: "value2", Expiration: time.Unix(14, 0)},
			},
		},
		{
			name:        "Negative interval",
			currentTime: 9,
			key:         "key2",
			value:       "value2",
			interval:    -5 * time.Second,
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
				"key2": {Value: "value2", Expiration: time.Unix(9, 0)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := []Option{withTestClock(newMockClock(tt.currentTime)), WithItems(tt.items)}
			opts = append(opts, tt.opts...)
			c := NewCache(opts...)
			c.AddInterval(tt.key, tt.value, tt.interval)
			got := c.GetAllItems()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddInterval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Add(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		opts        []Option
		key         string
		value       any
		items       map[string]Item
		want        map[string]Item
	}{
		{
			name:        "Key exists and not expired",
			currentTime: 9,
			key:         "key1",
			value:       "newvalue1",
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
		},
		{
			name:        "Key exists but allow overrides",
			currentTime: 9,
			key:         "key1",
			value:       "newvalue1",
			opts:        []Option{WithAllowOverwrites()},
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want: map[string]Item{
				"key1": {Value: "newvalue1", Expiration: time.Unix(10, 0)},
			},
		},
		{
			name:        "Key exists but expired",
			currentTime: 9,
			key:         "key1",
			value:       "newvalue1",
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(8, 0)},
			},
			want: map[string]Item{
				"key1": {Value: "newvalue1", Expiration: time.Unix(10, 0)},
			},
		},
		{
			name:        "Key does not exists",
			currentTime: 9,
			key:         "key2",
			value:       "value2",
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
				"key2": {Value: "value2", Expiration: time.Unix(10, 0)},
			},
		},
		{
			name:        "No items",
			currentTime: 9,
			key:         "key2",
			value:       "value2",
			items:       map[string]Item{},
			want: map[string]Item{
				"key2": {Value: "value2", Expiration: time.Unix(10, 0)},
			},
		},
		{
			name:        "Custom interval",
			currentTime: 9,
			key:         "key2",
			value:       "value2",
			opts:        []Option{WithCacheInterval(5 * time.Second)},
			items: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
			},
			want: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
				"key2": {Value: "value2", Expiration: time.Unix(14, 0)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := []Option{withTestClock(newMockClock(tt.currentTime)), WithItems(tt.items)}
			opts = append(opts, tt.opts...)
			c := NewCache(opts...)
			c.Add(tt.key, tt.value)
			got := c.GetAllItems()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_WithGarbageCollectorInterval(t *testing.T) {
	tests := []struct {
		name     string
		interval time.Duration
		want     time.Duration
	}{
		{
			name:     "Set garbage collector interval",
			interval: 5 * time.Second,
			want:     5 * time.Second,
		},
		{
			name:     "Set garbage collector interval to zero",
			interval: 0,
			want:     2 * time.Second,
		},
		{
			name:     "Set garbage collector interval to negative",
			interval: -5 * time.Second,
			want:     2 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCache(WithGarbageCollectorInterval(tt.interval))
			ca, ok := c.(*cache)
			if !ok {
				t.Fatalf("Expected cache type, got %T", c)
			}
			if tt.interval <= 0 {
				if ca.garbageCollector != nil {
					t.Errorf("WithGarbageCollectorInterval() = %v, want %v", ca.garbageCollector, nil)
				}
			} else if ca.garbageCollectorInterval != tt.want {
				t.Errorf("WithGarbageCollectorInterval() = %v, want %v", ca.garbageCollectorInterval, tt.want)
			}
		})
	}
}

func TestCache_WithMap(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		inputMap    map[string]any
		want        map[string]Item
	}{
		{
			name:        "Set map with items",
			currentTime: 9,
			inputMap: map[string]any{
				"key1": "value1",

				"key2": "value2",
			},
			want: map[string]Item{
				"key1": {Value: "value1", Expiration: time.Unix(10, 0)},
				"key2": {Value: "value2", Expiration: time.Unix(10, 0)},
			},
		},
		{
			name:        "Set map with empty items",
			currentTime: 9,
			inputMap:    map[string]any{},
			want:        map[string]Item{},
		},
		{
			name:        "Set map with nil items",
			currentTime: 9,
			inputMap:    nil,
			want:        map[string]Item{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCache(withTestClock(newMockClock(tt.currentTime)), WithMap(tt.inputMap))
			got := c.GetAllItems()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Stop(t *testing.T) {
	tests := []struct {
		name                string
		opts                []Option
		expectedInitialized bool
	}{
		{
			name:                "Default options",
			opts:                []Option{},
			expectedInitialized: true,
		},
		{
			name:                "No garbage collector",
			opts:                []Option{WithGarbageCollectorInterval(0)},
			expectedInitialized: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCache(tt.opts...)

			gc := c.GetGarbageCollector()
			if !tt.expectedInitialized && gc != nil {
				t.Fatalf("Expected garbage collector to be nil")
			} else if !tt.expectedInitialized {
				return
			}

			if gc == nil {
				t.Fatalf("Expected garbage collector to be initialized")
			}

			if !gc.active {
				t.Errorf("Expected garbage collector to be active")
			}

			c.Stop()
			gc = c.GetGarbageCollector()
			if gc == nil {
				t.Fatalf("Expected garbage collector to be initialized")
			}
			if gc.active {
				t.Errorf("Expected garbage collector to be inactive")
			}
		})
	}
}

func TestCache_GetGarbageCollector(t *testing.T) {
	tests := []struct {
		name        string
		opts        []Option
		expectedNil bool
	}{
		{
			name:        "Default options",
			opts:        []Option{},
			expectedNil: false,
		},
		{
			name:        "No garbage collector",
			opts:        []Option{WithGarbageCollectorInterval(0)},
			expectedNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCache(tt.opts...)

			gc := c.GetGarbageCollector()
			if (gc == nil) != tt.expectedNil {
				t.Fatalf("Expected garbage collector to be nil: %v", tt.expectedNil)
			}
		})
	}
}

func TestCache_SetGarbageCollector(t *testing.T) {
	tests := []struct {
		name string
		gc   *garbageCollector
	}{
		{
			name: "Set garbage collector",
			gc:   &garbageCollector{active: true},
		},
		{
			name: "Set nil garbage collector",
			gc:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCache(WithCacheInterval(0))
			c.SetGarbageCollector(tt.gc)

			gc := c.GetGarbageCollector()
			if tt.gc != gc {
				t.Fatalf("Expected garbage collector to be set")
			}
		})
	}
}

func TestCache_GetClock(t *testing.T) {
	tests := []struct {
		name  string
		clock *mockClock
	}{
		{
			name:  "Get clock",
			clock: &mockClock{},
		},
		{
			name:  "Get nil clock",
			clock: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCache(WithCacheInterval(0), withTestClock(tt.clock))

			got := c.GetClock()
			if tt.clock != got {
				t.Fatalf("Expected clock to be set")
			}
		})
	}
}

func TestCache_SetClock(t *testing.T) {
	tests := []struct {
		name  string
		clock *mockClock
	}{
		{
			name:  "Set clock",
			clock: &mockClock{},
		},
		{
			name:  "Set nil clock",
			clock: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCache(WithCacheInterval(0))
			c.SetClock(tt.clock)

			got := c.GetClock()
			if tt.clock != got {
				t.Fatalf("Expected clock to be set")
			}
		})
	}
}
