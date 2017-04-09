package handlers

import (
	"reflect"
	"testing"

	"github.com/rebeccaskinner/agile17-sample/user"
	"github.com/stretchr/testify/assert"
)

func TestAllUsers(t *testing.T) {
	t.Parallel()
	t.Run("when empty", func(t *testing.T) {
		t.Parallel()
		d := &Datastore{}
		u := d.allUsers()
		assert.Equal(t, 0, len(u))
	})
	t.Run("when-single-key", func(t *testing.T) {
		t.Parallel()
		t.Run("empty-list", func(t *testing.T) {
			t.Parallel()
			var m map[string][]*user.User
			d := &Datastore{oldData: m}
			u := d.allUsers()
			assert.Equal(t, 0, len(u))
		})
		t.Run("single-element-list", func(t *testing.T) {
			t.Parallel()
			l := []*user.User{user.RandomUser()}
			m := map[string][]*user.User{"foo": l}
			d := &Datastore{oldData: m}
			a := d.allUsers()
			assert.True(t, reflect.DeepEqual(a, l))
		})
		t.Run("multi-element-list", func(t *testing.T) {
			t.Parallel()
			l := []*user.User{user.RandomUser(), user.RandomUser(), user.RandomUser()}
			m := map[string][]*user.User{"foo": l}
			d := &Datastore{oldData: m}
			a := d.allUsers()
			assert.True(t, reflect.DeepEqual(a, l))
		})
	})
	t.Run("when-many-keys", func(t *testing.T) {
		t.Parallel()
		l1 := nRandUsers(3)
		l2 := nRandUsers(2)
		l3 := nRandUsers(8)
		m := map[string][]*user.User{
			"l1": l1,
			"l2": l2,
			"l3": l3,
		}
		d := &Datastore{oldData: m}
		a := d.allUsers()
		assert.True(t, reflect.DeepEqual(a, append(append(l1, l2...), l3...)))
	})
}

func nRandUsers(count int) []*user.User {
	u := make([]*user.User, count)
	for i := 0; i < count; i++ {
		u[i] = user.RandomUser()
	}
	return u
}
