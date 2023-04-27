package filelock

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestFileLock_Lock(t *testing.T) {
	getwd, _ := os.Getwd()
	lock := NewLock(filepath.Join(getwd))
	err := lock.Lock()
	assert.Nil(t, err)
	defer lock.Unlock()
}

