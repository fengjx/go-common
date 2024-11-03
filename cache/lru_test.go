package cache

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_lruCache_Get(t *testing.T) {
	c := NewLRUCache[string, string](10, 30, func(ctx context.Context, missKeys []string) (map[string]string, error) {
		m := make(map[string]string)
		for _, key := range missKeys {
			m[key] = key
		}
		return m, nil
	})
	result, _ := c.Get(context.Background(), "foo").Result()
	t.Log("result", result)
}

func Test_lruCache_GetMulti(t *testing.T) {
	type user struct {
		ID   int
		Name string
	}

	type testCase[K comparable, V any] struct {
		name string
		args []int
		want *Result[map[K]V]
	}
	tests := []testCase[int, *user]{
		{
			name: "GetMulti",
			args: []int{1, 2},
			want: &Result[map[int]*user]{
				val: map[int]*user{
					1: {ID: 1, Name: "user-1"},
					2: {ID: 2, Name: "user-2"},
				},
			},
		},
		{
			name: "Get",
			args: []int{1},
			want: &Result[map[int]*user]{
				val: map[int]*user{
					1: {ID: 1, Name: "user-1"},
				},
			},
		},
	}
	c := NewLRUCache[int, *user](10, 30, func(ctx context.Context, missKeys []int) (map[int]*user, error) {
		m := make(map[int]*user)
		for _, key := range missKeys {
			m[key] = &user{ID: key, Name: fmt.Sprintf("user-%d", key)}
		}
		t.Logf("fallback: %#v", m)
		return m, nil
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.GetMulti(context.Background(), tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMulti() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lruCache(t *testing.T) {
	type user struct {
		ID   int
		Name string
	}
	c := NewLRUCache[int, *user](10, 30, func(ctx context.Context, missKeys []int) (map[int]*user, error) {
		m := make(map[int]*user)
		for _, key := range missKeys {
			m[key] = &user{ID: key, Name: fmt.Sprintf("user-%d", key)}
		}
		t.Logf("fallback: %v", missKeys)
		return m, nil
	})
	u1, _ := c.Get(context.Background(), 1).Result()
	assert.Equal(t, "user-1", u1.Name)

	m, _ := c.GetMulti(context.Background(), []int{1, 2}).Result()
	assert.Equal(t, "user-1", m[1].Name)
	assert.Equal(t, "user-2", m[2].Name)

	c.Set(context.Background(), 3, &user{ID: 3, Name: "user-3"})
	u3, _ := c.Get(context.Background(), 3).Result()
	assert.Equal(t, "user-3", u3.Name)

	m2, _ := c.GetMulti(context.Background(), []int{3, 4}).Result()
	assert.Equal(t, "user-3", m2[3].Name)
	assert.Equal(t, "user-4", m2[4].Name)

	c.SetMulti(context.Background(), map[int]*user{
		5: {ID: 5, Name: "user-5"},
		6: {ID: 6, Name: "user-6"},
		7: {ID: 7, Name: "user-7"},
	})

	has, _ := c.Has(context.Background(), 1).Result()
	assert.Equal(t, true, has)

	cnt, _ := c.Del(context.Background(), 1, 2, 10).Result()
	assert.Equal(t, 2, cnt)

}

func TestTTL(t *testing.T) {
	type user struct {
		ID   int
		Name string
	}
	c := NewLRUCache[int, *user](10, time.Second*3, func(ctx context.Context, missKeys []int) (map[int]*user, error) {
		m := make(map[int]*user)
		for _, key := range missKeys {
			m[key] = &user{ID: key, Name: fmt.Sprintf("user-%d", key)}
		}
		t.Logf("fallback: %#v", m)
		return m, nil
	})
	u1, _ := c.Get(context.Background(), 1).Result()
	assert.Equal(t, "user-1", u1.Name)
	has, _ := c.Has(context.Background(), 1).Result()
	assert.Equal(t, true, has)

	time.Sleep(time.Second * 4)
	has, _ = c.Has(context.Background(), 1).Result()
	assert.Equal(t, false, has)
	u1, _ = c.Get(context.Background(), 1).Result()
}
