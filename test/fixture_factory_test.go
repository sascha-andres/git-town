//nolint:testpackage
package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixtureFactory(t *testing.T) {
	t.Parallel()
	t.Run("NewFixtureFactory()", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		_, err := NewFixtureFactory(dir)
		assert.Nil(t, err, "creating memoized environment failed")
		memoizedPath := filepath.Join(dir, "memoized")
		_, err = os.Stat(memoizedPath)
		assert.Falsef(t, os.IsNotExist(err), "memoized directory %q not found", memoizedPath)
	})

	t.Run(".CreateFixture()", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		gm, err := NewFixtureFactory(dir)
		assert.Nil(t, err, "creating memoized environment failed")
		result, err := gm.CreateFixture("foo")
		assert.Nil(t, err, "cannot create scenario environment")
		_, err = os.Stat(result.DevRepo.workingDir)
		assert.False(t, os.IsNotExist(err), "scenario environment directory %q not found", result.DevRepo.WorkingDir)
	})
}