package user

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUserFromUser(t *testing.T) {
	t.Run("test-business-units", func(t *testing.T) {
		departmentTableTest(t, map[string][]string{
			"engineering": []string{
				"software developer",
				"senior software developer",
				"lead software developer",
				"director of engineering",
				"lead tester",
			},
			"executive": []string{
				"cto",
				"cfo",
				"CEO",
			},
			"operations": []string{
				"mid-level dba",
				"senior SRE",
				"lead system administrator",
				"storage administrator",
			},
		})
	})
}

func departmentTableTest(t *testing.T, m map[string][]string) {
	for k, vals := range m {
		for _, v := range vals {
			t.Logf("asserting %s is the department for title %s", k, v)
			assertTitleDepartment(t, v, k)
		}
	}
}

func assertTitleDepartment(t *testing.T, title, department string) {
	oldU := &User{Title: title}
	newU, err := NewUserFromUser(oldU)
	require.NoError(t, err)
	t.Logf("asserting %s == %s", department, newU.BusinessUnit)
	assert.Equal(t, department, newU.BusinessUnit, fmt.Sprintf("title %s != department %s", title, department))
}
