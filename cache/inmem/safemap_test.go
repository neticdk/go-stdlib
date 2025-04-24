package inmem

import (
	"context"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/neticdk/go-stdlib/assert"
	"github.com/neticdk/go-stdlib/cache"
)

func newTestGarbageCollector(interval time.Duration, c *mockClock) GarbageCollector[string, any] {
	gc := NewGarbageCollector[string, any](interval)
	gc.clock = c
	return gc
}

func newTestSafeMapCache(items map[string]*item[any], currentTime int64, opts ...SafeMapCacheOption[string, any]) *safeMapCache[string, any] {
	mc := newMockClock(currentTime)
	c := &safeMapCache[string, any]{
		items:             items,
		clock:             mc,
		count:             atomic.Int64{},
		defaultTTL:        defaultTTL,
		defaultGCInterval: defaultGCInterval,
		maxSize:           maxSize,
	}
	for _, opt := range opts {
		opt(c)
	}
	c.count.Store(int64(len(items)))
	gc := newTestGarbageCollector(5*time.Minute, mc)
	if gc != nil {
		c.garbageCollector = gc
		go c.garbageCollector.Start(c)
	}
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

func TestSafeMapCache_NewSafeMap(t *testing.T) {
	opts := []SafeMapCacheOption[string, any]{
		WithDefaultTTL[string, any](10 * time.Second),
		WithMaxSize[string, any](100),
		WithGCInterval[string, any](4 * time.Second),
	}

	c := NewSafeMap(opts...)
	assert.NotNil(t, c, "NewSafeMap() should not return nil")
	assert.Equal(t, c.Len(), 0, "NewSafeMap() should initialize with length 0")
	assert.Equal(t, c.defaultTTL, 10*time.Second, "NewSafeMap() should set default TTL to 10 seconds")
	assert.Equal(t, c.maxSize, 100, "NewSafeMap() should set max size to 100")
	assert.Equal(t, c.defaultGCInterval, 5*time.Second, "NewSafeMap() should set GC interval to 5 seconds")
}

func TestSafeMapCache_Get(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		key         string
		ctx         context.Context
		items       map[string]*item[any]
		want        any
		wantErr     bool
	}{
		{
			name:        "Key exists and not expired",
			currentTime: 9,
			key:         "key1",
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want:    "value1",
			wantErr: false,
		},
		{
			name:        "Key exists but expired",
			currentTime: 11,
			key:         "key1",
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:        "Key does not exist",
			currentTime: 9,
			key:         "key2",
			items: map[string]*item[any]{
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
			items: map[string]*item[any]{
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
			items: map[string]*item[any]{
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
			if tt.wantErr {
				assert.Error(t, err, "Get() expected error, got %v", err)
				return
			} else {
				assert.NoError(t, err, "Get() expected no error, got %v", err)
			}
			assert.Equal(t, got, tt.want, "Get() = %v, want %v", got, tt.want)
		})
	}
}

func TestSafeMapCache_Delete(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		key         string
		ctx         context.Context
		items       map[string]*item[any]
		want        map[string]*item[any]
		missing     bool
		wantErr     bool
	}{
		{
			name:        "Key exists and not expired",
			currentTime: 9,
			key:         "key1",
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want:    map[string]*item[any]{},
			wantErr: false,
		},
		{
			name:        "Key exists but expired",
			currentTime: 11,
			key:         "key1",
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want:    map[string]*item[any]{},
			wantErr: false,
		},
		{
			name:        "Key does not exist",
			currentTime: 9,
			key:         "key2",
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			missing: true,
			wantErr: false,
		},
		{
			name:        "Context error",
			currentTime: 9,
			key:         "key1",
			ctx:         newMockContext(1),
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want:    map[string]*item[any]{},
			wantErr: true,
		},
		{
			name:        "Context error after locking",
			currentTime: 9,
			key:         "key1",
			ctx:         newMockContext(2),
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want:    map[string]*item[any]{},
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
			if tt.wantErr {
				assert.Error(t, err, "Delete() expected error, got %v", err)
				return
			} else {
				assert.NoError(t, err, "Delete() expected no error, got %v", err)
			}
			_, err = c.Get(ctx, tt.key)
			if tt.missing {
				assert.ErrorIs(t, err, cache.NewErrCacheMiss(), "Get() expected cache miss, got %v", err)
			} else {
				assert.ErrorIs(t, err, cache.NewErrExpired(), "Get() expected expired error, got %v", err)
			}
			assert.NoError(t, c.deleteExpired(ctx))
			newItems := c.items
			assert.Equal(t, newItems, tt.want, "Delete() = %v, want %v", newItems, tt.want)
		})
	}
}

func TestSafeMapCache_Clear(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		ctx         context.Context
		items       map[string]*item[any]
		wantErr     bool
	}{
		{
			name:        "Existing items",
			currentTime: 9,
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(5, 0)},
			},
			wantErr: false,
		},
		{
			name:        "No items",
			currentTime: 11,
			items:       map[string]*item[any]{},
			wantErr:     false,
		},
		{
			name:        "Context error",
			currentTime: 9,
			ctx:         newMockContext(1),
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(5, 0)},
			},
			wantErr: true,
		},
		{
			name:        "Context error after locking",
			currentTime: 9,
			ctx:         newMockContext(2),
			items: map[string]*item[any]{
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
			if tt.wantErr {
				assert.Error(t, err, "Clear() expected error, got %v", err)
				return
			} else {
				assert.NoError(t, err, "Clear() expected no error, got %v", err)
			}
			got := c.items
			assert.Equal(t, got, map[string]*item[any]{}, "Clear() = %v, want %v", got, map[string]*item[any]{})
		})
	}
}

func TestSafeMapCache_deleteExpired(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		ctx         context.Context
		items       map[string]*item[any]
		want        map[string]*item[any]
		wantErr     bool
	}{
		{
			name:        "With both expired and non-expired items",
			currentTime: 10,
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(8, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(9, 0)},
				"key3": {value: "value3", expiresAt: time.Unix(10, 0)},
			},
			want: map[string]*item[any]{
				"key3": {value: "value3", expiresAt: time.Unix(10, 0)},
			},
			wantErr: false,
		},
		{
			name:        "All expired",
			currentTime: 11,
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(8, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(9, 0)},
				"key3": {value: "value3", expiresAt: time.Unix(10, 0)},
			},
			want:    map[string]*item[any]{},
			wantErr: false,
		},
		{
			name:        "No expirations",
			currentTime: 5,
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(8, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(9, 0)},
				"key3": {value: "value3", expiresAt: time.Unix(10, 0)},
			},
			want: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(8, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(9, 0)},
				"key3": {value: "value3", expiresAt: time.Unix(10, 0)},
			},
			wantErr: false,
		},
		{
			name:        "No items",
			currentTime: 11,
			items:       map[string]*item[any]{},
			want:        map[string]*item[any]{},
			wantErr:     false,
		},
		{
			name:        "Context error",
			currentTime: 10,
			ctx:         newMockContext(1),
			items: map[string]*item[any]{
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
			items: map[string]*item[any]{
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
			if tt.wantErr {
				assert.Error(t, err, "deleteExpired() expected error, got %v", err)
				return
			} else {
				assert.NoError(t, err, "deleteExpired() expected no error, got %v", err)
			}
			got := c.items
			assert.Equal(t, got, tt.want, "deleteExpired() = %v, want %v", got, tt.want)
		})
	}
}

func TestSafeMapCache_Set(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		key         string
		value       any
		ctx         context.Context
		opts        []SafeMapCacheOption[string, any]
		items       map[string]*item[any]
		want        map[string]*item[any]
		wantErr     bool
	}{
		{
			name:        "Key exists and not expired",
			currentTime: 9,
			key:         "key1",
			value:       "newvalue1",
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(0, 0)},
			},
			want: map[string]*item[any]{
				"key1": {value: "newvalue1", expiresAt: time.Time{}},
			},
			wantErr: false,
		},
		{
			name:        "Key exists but expired",
			currentTime: 9,
			key:         "key1",
			value:       "newvalue1",
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(8, 0)},
			},
			want: map[string]*item[any]{
				"key1": {value: "newvalue1", expiresAt: time.Time{}},
			},
			wantErr: false,
		},
		{
			name:        "Key does not exists",
			currentTime: 9,
			key:         "key2",
			value:       "value2",
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want: map[string]*item[any]{
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
			items:       map[string]*item[any]{},
			want: map[string]*item[any]{
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
			items: map[string]*item[any]{
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
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(0, 0)},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:        "With custom default TTL",
			currentTime: 9,
			key:         "key2",
			value:       "value2",
			opts:        []SafeMapCacheOption[string, any]{WithDefaultTTL[string, any](10 * time.Second)},
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(19, 0)},
			},
			wantErr: false,
		},
		{
			name:        "Exceeding max size",
			currentTime: 9,
			key:         "key2",
			value:       "value2",
			opts:        []SafeMapCacheOption[string, any]{WithMaxSize[string, any](1)},
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestSafeMapCache(tt.items, tt.currentTime, tt.opts...)
			ctx := context.TODO()
			if tt.ctx != nil {
				ctx = tt.ctx
			}
			err := c.Set(ctx, tt.key, tt.value)
			if tt.wantErr {
				assert.Error(t, err, "Set() expected error, got %v", err)
				return
			} else {
				assert.NoError(t, err, "Set() expected no error, got %v", err)
			}
			got := c.items
			assert.Equal(t, got, tt.want, "Set() = %v, want %v", got, tt.want)
		})
	}
}

func TestSafeMapCache_SetWithTTL(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		key         string
		value       any
		ttl         time.Duration
		ctx         context.Context
		items       map[string]*item[any]
		want        map[string]*item[any]
		wantErr     bool
	}{
		{
			name:        "Key exists and not expired",
			currentTime: 9,
			key:         "key1",
			value:       "newvalue1",
			ttl:         1 * time.Second,
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want: map[string]*item[any]{
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
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(8, 0)},
			},
			want: map[string]*item[any]{
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
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
			},
			want: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(10, 0)},
			},
			wantErr: false,
		},
		{
			name:        "Previous item not expired but new item is",
			currentTime: 9,
			key:         "key1",
			value:       "newvalue1",
			ttl:         -1 * time.Second,
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Time{}},
			},
			want: map[string]*item[any]{
				"key1": {value: "newvalue1", expiresAt: time.Unix(8, 0)},
			},
			wantErr: false,
		},
		{
			name:        "No items",
			currentTime: 9,
			key:         "key2",
			value:       "value2",
			ttl:         1 * time.Second,
			items:       map[string]*item[any]{},
			want: map[string]*item[any]{
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
			items: map[string]*item[any]{
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
			items: map[string]*item[any]{
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
			if tt.wantErr {
				assert.Error(t, err, "Set() expected error, got %v", err)
				return
			} else {
				assert.NoError(t, err, "Set() expected no error, got %v", err)
			}
			got := c.items
			assert.Equal(t, got, tt.want, "Set() = %v, want %v", got, tt.want)
		})
	}
}

func TestSafeMapCache_Len(t *testing.T) {
	tests := []struct {
		name        string
		currentTime int64
		items       map[string]*item[any]
		want        int
		wantAfter   int
	}{
		{
			name:        "Non-empty cache",
			currentTime: 9,
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(11, 0)},
			},
			want:      2,
			wantAfter: 2,
		},
		{
			name:        "Empty cache",
			currentTime: 9,
			items:       map[string]*item[any]{},
			want:        0,
			wantAfter:   0,
		},
		{
			name:        "Cache with expired items",
			currentTime: 11,
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(8, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(9, 0)},
			},
			want:      2,
			wantAfter: 0,
		},
		{
			name:        "Cache with mix of expired and non-expired items",
			currentTime: 9,
			items: map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(8, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(9, 0)},
				"key3": {value: "value3", expiresAt: time.Unix(10, 0)},
			},
			want:      3,
			wantAfter: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestSafeMapCache(tt.items, tt.currentTime)
			got := c.Len()
			if got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
			// Simulate garbage collection
			assert.NoError(t, c.deleteExpired(context.TODO()))
			gotAfter := c.Len()
			assert.Equal(t, gotAfter, tt.wantAfter, "Len() after deleteExpired = %v, want %v", gotAfter, tt.wantAfter)
		})
	}
}

func newMockTimeoutContext(timeout time.Duration, count int) context.Context {
	ctx, cancel := context.WithTimeout(newMockContext(count), timeout)
	defer cancel()
	return ctx
}

func TestSafeMapCache_Stop(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		wantErr bool
	}{
		{
			name:    "Stop with no error",
			ctx:     context.TODO(),
			wantErr: false,
		},
		{
			name:    "Stop with initial error",
			ctx:     newMockContext(1),
			wantErr: true,
		},
		{
			name:    "Stop with error after locking",
			ctx:     newMockContext(2),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			items := map[string]*item[any]{
				"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
				"key2": {value: "value2", expiresAt: time.Unix(11, 0)},
			}

			c := newTestSafeMapCache(items, 9)

			err := c.Stop(tt.ctx)
			if tt.wantErr {
				assert.Error(t, err, "Stop() expected error, got %v", err)
				return
			} else {
				assert.NoError(t, err, "Stop() expected no error, got %v", err)
			}
			got := c.items
			assert.Equal(t, got, map[string]*item[any]{}, "Stop() = %v, want %v", got, map[string]*item[any]{})
		})
	}

}

func TestSafeMapCache_ConcurrentAccess(t *testing.T) {
	items := map[string]*item[any]{
		"key1": {value: "value1", expiresAt: time.Unix(10, 0)},
		"key2": {value: "value2", expiresAt: time.Unix(11, 0)},
	}

	c := newTestSafeMapCache(items, 9)
	ctx := context.TODO()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			_, _ = c.Get(ctx, "key1")
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			_ = c.Set(ctx, "key2", "newvalue")
		}
	}()

	wg.Wait()

	assert.Equal(t, c.Len(), 2)
}

func createItemMap(size int) map[string]*item[any] {
	items := make(map[string]*item[any], size)
	for i := range size {
		items[fmt.Sprintf("key%d", i)] = &item[any]{value: fmt.Sprintf("value%d", i), expiresAt: time.Unix(int64(i), 0)}
	}
	return items
}

func BenchmarkSafeMapCache_Get(b *testing.B) {
	sizes := []int{1, 1000, 1000000}

	for _, s := range sizes {
		b.Run(fmt.Sprintf("size%d", s), func(b *testing.B) {
			items := createItemMap(s)
			c := newTestSafeMapCache(items, 9)
			b.ResetTimer()
			for i := range b.N {
				_, _ = c.Get(context.TODO(), fmt.Sprint("key", i%s))
			}
		})
	}
}

func BenchmarkSafeMapCache_Set(b *testing.B) {
	sizes := []int{1, 1000, 1000000}

	for _, s := range sizes {
		b.Run(fmt.Sprintf("size%d", s), func(b *testing.B) {
			items := createItemMap(s)
			c := newTestSafeMapCache(items, 9)
			b.ResetTimer()
			for i := range b.N {
				_ = c.Set(context.TODO(), fmt.Sprint("key", i%s), "newvalue")
			}
		})
	}
}

func BenchmarkSafeMapCache_Clear(b *testing.B) {
	sizes := []int{1, 1000, 1000000}

	for _, s := range sizes {
		b.Run(fmt.Sprintf("size%d", s), func(b *testing.B) {
			items := createItemMap(s)
			c := newTestSafeMapCache(items, 9)
			b.ResetTimer()
			for b.Loop() {
				_ = c.Clear(context.TODO())
				c.items = items
			}
		})
	}
}

func BenchmarkSafeMapCache_Delete(b *testing.B) {
	sizes := []int{1, 1000, 1000000}

	for _, s := range sizes {
		b.Run(fmt.Sprintf("size%d", s), func(b *testing.B) {
			items := createItemMap(s)
			c := newTestSafeMapCache(items, 9)
			b.ResetTimer()
			for i := range b.N {
				_ = c.Delete(context.TODO(), fmt.Sprint("key", i%s))
				c.items = items
			}
		})
	}
}

func BenchmarkSafeMapCache_Len(b *testing.B) {
	sizes := []int{1, 1000, 1000000}

	for _, s := range sizes {
		b.Run(fmt.Sprintf("size%d", s), func(b *testing.B) {
			items := createItemMap(s)
			c := newTestSafeMapCache(items, 9)
			b.ResetTimer()
			for b.Loop() {
				_ = c.Len()
			}
		})
	}
}

func BenchmarkSafeMapCache_deleteExpiredNoDeletions(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000, 100000, 1000000}
	times := []float32{0.0, 0.5, 1.0}

	for _, s := range sizes {
		for _, t := range times {
			time := int64(math.Round(float64(s) * float64(t)))
			b.Run(fmt.Sprintf("size%d-time%d", s, time), func(b *testing.B) {
				items := createItemMap(s)
				c := newTestSafeMapCache(items, time)
				assert.Equal(b, c.Len(), s)
				b.ResetTimer()
				_ = c.deleteExpired(context.TODO())
			})
		}
	}
}
