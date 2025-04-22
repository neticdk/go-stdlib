package inmem

import (
	"context"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/neticdk/go-stdlib/assert"
)

func newTestGarbageCollector(interval time.Duration, c *mockClock) GarbageCollector[string, any] {
	gc := NewGarbageCollector[string, any](interval)
	gc.clock = c
	return gc
}

func newTestSafeMapCache(items map[string]item[any], currentTime int64) *safeMapCache[string, any] {
	mc := newMockClock(currentTime)
	c := &safeMapCache[string, any]{
		items: items,
		clock: mc,
	}
	gc := newTestGarbageCollector(5*time.Minute, mc)
	if gc != nil {
		c.garbageCollector = gc
		go c.garbageCollector.Start(c)
	}
	runtime.AddCleanup(c, stopGarbageCollector[string, any], c.garbageCollector)
	return c
}

func newMockContext(calls int) context.Context {
	return &mockContext{calls: calls, called: 0}
}

type mockContext struct {
	calls  int
	called int
}

func (m *mockContext) Err() error {
	m.called++
	if m.called >= m.calls {
		return context.Canceled
	}
	return nil
}
func (m *mockContext) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}
func (m *mockContext) Done() <-chan struct{} {
	return nil
}
func (m *mockContext) Value(key any) any {
	return nil
}

func TestCache_Get(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		key         string
		ctx         context.Context
		items       map[string]item[any]
		want        any
		wantErr     bool
	}{
		{
			name:        "Key exists and not expired",
			currentTime: 9,
			key:         "key1",
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want:    "value1",
			wantErr: false,
		},
		{
			name:        "Key exists but expired",
			currentTime: 11,
			key:         "key1",
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:        "Key does not exist",
			currentTime: 9,
			key:         "key2",
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:        "Context error",
			currentTime: 9,
			key:         "key1",
			ctx:         newMockContext(1),
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:        "Context error after locking",
			currentTime: 9,
			key:         "key1",
			ctx:         newMockContext(2),
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestSafeMapCache(tt.items, tt.currentTime)
			ctx := context.TODO()
			if tt.ctx != nil {
				ctx = tt.ctx
			}
			got, err := c.Get(ctx, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() expected error %t, got %v", tt.wantErr, err)
			} else if err != nil {
				return
			}
			if got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestCache_Delete(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		key         string
		ctx         context.Context
		items       map[string]item[any]
		want        map[string]item[any]
		wantErr     bool
	}{
		{
			name:        "Key exists and not expired",
			currentTime: 9,
			key:         "key1",
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want:    map[string]item[any]{},
			wantErr: false,
		},
		{
			name:        "Key exists but expired",
			currentTime: 11,
			key:         "key1",
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want:    map[string]item[any]{},
			wantErr: false,
		},
		{
			name:        "Key does not exist",
			currentTime: 9,
			key:         "key2",
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			wantErr: false,
		},
		{
			name:        "Context error",
			currentTime: 9,
			key:         "key1",
			ctx:         newMockContext(1),
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want:    map[string]item[any]{},
			wantErr: true,
		},
		{
			name:        "Context error after locking",
			currentTime: 9,
			key:         "key1",
			ctx:         newMockContext(2),
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want:    map[string]item[any]{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestSafeMapCache(tt.items, tt.currentTime)
			ctx := context.TODO()
			if tt.ctx != nil {
				ctx = tt.ctx
			}
			err := c.Delete(ctx, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() expected error %t, got %v", tt.wantErr, err)
			} else if err != nil {
				return
			}
			got := c.items
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
		ctx         context.Context
		items       map[string]item[any]
		wantErr     bool
	}{
		{
			name:        "Existing items",
			currentTime: 9,
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(5, 0)},
			},
			wantErr: false,
		},
		{
			name:        "No items",
			currentTime: 11,
			items:       map[string]item[any]{},
			wantErr:     false,
		},
		{
			name:        "Context error",
			currentTime: 9,
			ctx:         newMockContext(1),
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(5, 0)},
			},
			wantErr: true,
		},
		{
			name:        "Context error after locking",
			currentTime: 9,
			ctx:         newMockContext(2),
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(5, 0)},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestSafeMapCache(tt.items, tt.currentTime)
			ctx := context.TODO()
			if tt.ctx != nil {
				ctx = tt.ctx
			}
			err := c.Clear(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Clear() expected error %t, got %v", tt.wantErr, err)
			} else if err != nil {
				return
			}
			got := c.items
			if len(got) != 0 {
				t.Errorf("Clear() = %v, want %v", got, map[string]item[any]{})
			}
		})
	}
}

func TestCache_deleteExpired(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		ctx         context.Context
		items       map[string]item[any]
		want        map[string]item[any]
		wantErr     bool
	}{
		{
			name:        "With both expired and non-expired items",
			currentTime: 10,
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(8, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(9, 0)},
				"key3": {value: "value3", expiresAt: time.Unix(10, 0)},
			},
			want: map[string]item[any]{
				"key3": {value: "value3", expiresAt: time.Unix(10, 0)},
			},
			wantErr: false,
		},
		{
			name:        "All expired",
			currentTime: 11,
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(8, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(9, 0)},
				"key3": {value: "value3", expiresAt: time.Unix(10, 0)},
			},
			want:    map[string]item[any]{},
			wantErr: false,
		},
		{
			name:        "No expirations",
			currentTime: 5,
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(8, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(9, 0)},
				"key3": {value: "value3", expiresAt: time.Unix(10, 0)},
			},
			want: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(8, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(9, 0)},
				"key3": {value: "value3", expiresAt: time.Unix(10, 0)},
			},
			wantErr: false,
		},
		{
			name:        "No items",
			currentTime: 11,
			items:       map[string]item[any]{},
			want:        map[string]item[any]{},
			wantErr:     false,
		},
		{
			name:        "Context error",
			currentTime: 10,
			ctx:         newMockContext(1),
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(8, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(9, 0)},
				"key3": {value: "value3", expiresAt: time.Unix(10, 0)},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:        "Context error after locking",
			currentTime: 10,
			ctx:         newMockContext(2),
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(8, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(9, 0)},
				"key3": {value: "value3", expiresAt: time.Unix(10, 0)},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestSafeMapCache(tt.items, tt.currentTime)
			ctx := context.TODO()
			if tt.ctx != nil {
				ctx = tt.ctx
			}
			err := c.deleteExpired(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("deleteExpired() expected error %t, got %v", tt.wantErr, err)
			} else if err != nil {
				return
			}
			got := c.items
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Set(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		key         string
		value       any
		ctx         context.Context
		items       map[string]item[any]
		want        map[string]item[any]
		wantErr     bool
	}{
		{
			name:        "Key exists and not expired",
			currentTime: 9,
			key:         "key1",
			value:       "newvalue1",
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(0, 0)},
			},
			want: map[string]item[any]{
				"key1": {value: "newvalue1", expiresAt: time.Time{}},
			},
			wantErr: false,
		},
		{
			name:        "Key exists but expired",
			currentTime: 9,
			key:         "key1",
			value:       "newvalue1",
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(8, 0)},
			},
			want: map[string]item[any]{
				"key1": {value: "newvalue1", expiresAt: time.Time{}},
			},
			wantErr: false,
		},
		{
			name:        "Key does not exists",
			currentTime: 9,
			key:         "key2",
			value:       "value2",
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
				"key2": {value: "value2", expiresAt: time.Time{}},
			},
			wantErr: false,
		},
		{
			name:        "No items",
			currentTime: 9,
			key:         "key2",
			value:       "value2",
			items:       map[string]item[any]{},
			want: map[string]item[any]{
				"key2": {value: "value2", expiresAt: time.Time{}},
			},
			wantErr: false,
		},
		{
			name:        "Context error",
			currentTime: 9,
			key:         "key1",
			value:       "newvalue1",
			ctx:         newMockContext(1),
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(0, 0)},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:        "Context error after locking",
			currentTime: 9,
			key:         "key1",
			value:       "newvalue1",
			ctx:         newMockContext(2),
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(0, 0)},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestSafeMapCache(tt.items, tt.currentTime)
			ctx := context.TODO()
			if tt.ctx != nil {
				ctx = tt.ctx
			}
			err := c.Set(ctx, tt.key, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set() expected error %t, got %v", tt.wantErr, err)
			} else if err != nil {
				return
			}
			got := c.items
			assert.Equal(t, got, tt.want, "Set() = %v, want %v", got, tt.want)
		})
	}
}

func TestCache_SetWithTTL(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		key         string
		value       any
		ttl         time.Duration
		ctx         context.Context
		items       map[string]item[any]
		want        map[string]item[any]
		wantErr     bool
	}{
		{
			name:        "Key exists and not expired",
			currentTime: 9,
			key:         "key1",
			value:       "newvalue1",
			ttl:         1 * time.Second,
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want: map[string]item[any]{
				"key1": {value: "newvalue1", expiresAt: time.Unix(10, 0)},
			},
			wantErr: false,
		},
		{
			name:        "Key exists but expired",
			currentTime: 9,
			key:         "key1",
			value:       "newvalue1",
			ttl:         1 * time.Second,
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(8, 0)},
			},
			want: map[string]item[any]{
				"key1": {value: "newvalue1", expiresAt: time.Unix(10, 0)},
			},
			wantErr: false,
		},
		{
			name:        "Key does not exists",
			currentTime: 9,
			key:         "key2",
			value:       "value2",
			ttl:         1 * time.Second,
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(10, 0)},
			},
			wantErr: false,
		},
		{
			name:        "No items",
			currentTime: 9,
			key:         "key2",
			value:       "value2",
			ttl:         1 * time.Second,
			items:       map[string]item[any]{},
			want: map[string]item[any]{
				"key2": {value: "value2", expiresAt: time.Unix(10, 0)},
			},
			wantErr: false,
		},
		{
			name:        "Context error",
			currentTime: 9,
			key:         "key1",
			value:       "newvalue1",
			ttl:         1 * time.Second,
			ctx:         newMockContext(1),
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:        "Context error after locking",
			currentTime: 9,
			key:         "key1",
			value:       "newvalue1",
			ttl:         1 * time.Second,
			ctx:         newMockContext(2),
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestSafeMapCache(tt.items, tt.currentTime)
			ctx := context.TODO()
			if tt.ctx != nil {
				ctx = tt.ctx
			}
			err := c.SetWithTTL(ctx, tt.key, tt.value, tt.ttl)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set() expected error %t, got %v", tt.wantErr, err)
			} else if err != nil {
				return
			}
			got := c.items
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Len(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		items       map[string]item[any]
		want        int
	}{
		{
			name:        "Non-empty cache",
			currentTime: 9,
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(11, 0)},
			},
			want: 2,
		},
		{
			name:        "Empty cache",
			currentTime: 9,
			items:       map[string]item[any]{},
			want:        0,
		},
		{
			name:        "Cache with expired items",
			currentTime: 11,
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(8, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(9, 0)},
			},
			want: 0,
		},
		{
			name:        "Cache with mix of expired and non-expired items",
			currentTime: 9,
			items: map[string]item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(8, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(9, 0)},
				"key3": {value: "value3", expiresAt: time.Unix(10, 0)},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestSafeMapCache(tt.items, tt.currentTime)
			got := c.Len()
			if got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}
