package storage

import (
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/heathcliff26/netrouse/pkg/server/storage/file"
	"github.com/heathcliff26/netrouse/pkg/server/storage/testsuite"
	"github.com/heathcliff26/netrouse/pkg/server/storage/types"
	"github.com/heathcliff26/netrouse/pkg/server/storage/valkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockBackend struct {
	mock.Mock
	types.StorageBackend
}

func (m *MockBackend) AddHost(host types.Host) error {
	args := m.Called(host)
	return args.Error(0)
}

func (m *MockBackend) RemoveHost(mac string) error {
	args := m.Called(mac)
	return args.Error(0)
}

func (m *MockBackend) GetHost(mac string) (types.Host, error) {
	args := m.Called(mac)
	return args.Get(0).(types.Host), args.Error(1)
}

func (m *MockBackend) GetHosts() ([]types.Host, error) {
	args := m.Called()
	return args.Get(0).([]types.Host), args.Error(1)
}

func (m *MockBackend) Readonly() (bool, error) {
	args := m.Called()
	return args.Bool(0), args.Error(1)
}

func TestNewStorage(t *testing.T) {
	t.Run("FileBackend", func(t *testing.T) {
		assert := assert.New(t)

		cfg := StorageConfig{
			Type: "file",
			File: file.FileBackendConfig{
				Path: t.TempDir() + "/test.yaml",
			},
			Readonly: true,
		}

		s, err := NewStorage(cfg)
		assert.NoError(err)
		assert.NotNil(s)
		assert.IsType(&file.FileBackend{}, s.backend, "Should have file as backend type")
		assert.Equal(cfg.Readonly, s.readonly, "Readonly should match config")
	})

	t.Run("ValkeyBackend", func(t *testing.T) {
		assert := assert.New(t)

		cfg := StorageConfig{
			Type: "valkey",
		}

		s, err := NewStorage(cfg)
		assert.Error(err, "Should fail to connect to valkey server")
		assert.Nil(s, "Storage should be nil")

		mr := miniredis.RunT(t)
		cfg.Valkey.Addrs = []string{mr.Addr()}

		s, err = NewStorage(cfg)
		assert.NoError(err, "Should not return an error")
		assert.NotNil(s, "Storage should not be nil")
		assert.IsType(&valkey.ValkeyBackend{}, s.backend, "Should have valkey as backend type")
	})

	t.Run("UnknownBackend", func(t *testing.T) {
		assert := assert.New(t)

		s, err := NewStorage(StorageConfig{Type: "unknown"})
		assert.Nil(s, "Storage should be nil")
		assert.Error(err)
		assert.Contains(err.Error(), "unknown storage backend type")
	})

	t.Run("WritableBackend", func(t *testing.T) {
		assert := assert.New(t)

		cfg := StorageConfig{
			Type: "file",
			File: file.FileBackendConfig{
				Path: t.TempDir() + "/test.yaml",
			},
			Readonly: false,
		}

		s, err := NewStorage(cfg)
		assert.NoError(err, "Should not return an error")
		assert.NotNil(s, "Storage should not be nil")

		assert.False(s.Readonly())
	})

	t.Run("OverwriteConfigWhenReadonlyBackend", func(t *testing.T) {
		if testsuite.IsRoot() {
			t.Skip("Running as root")
		}
		assert := assert.New(t)

		path := t.TempDir() + "/test.yaml"

		f, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0444)
		require.NoError(t, err, "Should create file")
		f.Close()

		cfg := StorageConfig{
			Type: "file",
			File: file.FileBackendConfig{
				Path: path,
			},
			Readonly: false,
		}

		s, err := NewStorage(cfg)
		assert.NoError(err, "Should not return an error")
		assert.NotNil(s, "Storage should not be nil")

		assert.True(s.Readonly(), "Should ignore config if backend is readonly")
	})

	t.Run("SeededStorageReadonly", func(t *testing.T) {
		assert := assert.New(t)

		cfg := StorageConfig{
			Type: "file",
			File: file.FileBackendConfig{
				Path: t.TempDir() + "/test.yaml",
			},
			Readonly:    true,
			SeededHosts: "file/testdata/basic.yaml",
		}

		s, err := NewStorage(cfg)
		assert.Nil(s, "Storage should be nil")
		assert.Error(err, "Should return an error when trying to seed hosts in readonly mode")
	})

	t.Run("SeededMissingFile", func(t *testing.T) {
		assert := assert.New(t)

		cfg := StorageConfig{
			Type: "file",
			File: file.FileBackendConfig{
				Path: t.TempDir() + "/test.yaml",
			},
			SeededHosts: "not-a-file.yaml",
		}

		s, err := NewStorage(cfg)
		assert.Nil(s, "Storage should be nil")
		assert.Error(err, "Should return an error when seed file is missing")
	})

	t.Run("SeededFileWrongFormat", func(t *testing.T) {
		assert := assert.New(t)

		cfg := StorageConfig{
			Type: "file",
			File: file.FileBackendConfig{
				Path: t.TempDir() + "/test.yaml",
			},
			SeededHosts: "file/testdata/not-yaml.txt",
		}

		s, err := NewStorage(cfg)
		assert.Nil(s, "Storage should be nil")
		assert.Error(err, "Should return an error when seed file is not yaml")
	})

	t.Run("Seeded", func(t *testing.T) {
		assert := assert.New(t)

		cfg := StorageConfig{
			Type: "file",
			File: file.FileBackendConfig{
				Path: t.TempDir() + "/test.yaml",
			},
			SeededHosts: "file/testdata/basic.yaml",
		}

		s, err := NewStorage(cfg)
		assert.NotNil(s, "Storage should not be nil")
		assert.NoError(err, "Should create storage")

		hosts, err := s.backend.GetHosts()
		assert.NoError(err, "Should get hosts without error")
		assert.Len(hosts, 2, "Should have 2 hosts from seeding")
	})
}

func TestStorageReadonly(t *testing.T) {
	assert := assert.New(t)

	s := &Storage{
		readonly: true,
	}

	assert.True(s.Readonly())

	s.readonly = false
	assert.False(s.Readonly())
}

func TestStorageAddHost(t *testing.T) {
	mockBackend := new(MockBackend)
	s := &Storage{
		backend:  mockBackend,
		readonly: false,
	}

	t.Run("Readonly", func(t *testing.T) {
		assert := assert.New(t)

		s.readonly = true
		err := s.AddHost(types.Host{MAC: "00:11:22:33:44:55", Name: "test"})
		assert.Error(err)
		assert.Contains(err.Error(), "storage is readonly")
	})

	t.Run("Success", func(t *testing.T) {
		assert := assert.New(t)

		s.readonly = false
		mockBackend.On("AddHost", types.Host{MAC: "00:11:22:33:44:55", Name: "test"}).Return(nil)

		err := s.AddHost(types.Host{MAC: "00:11:22:33:44:55", Name: "test"})
		assert.NoError(err, "Should add host without error")
		mockBackend.AssertExpectations(t)
	})
}

func TestStorageRemoveHost(t *testing.T) {
	mockBackend := new(MockBackend)
	s := &Storage{
		backend:  mockBackend,
		readonly: false,
	}

	t.Run("ReadonlyStorage", func(t *testing.T) {
		assert := assert.New(t)

		s.readonly = true
		err := s.RemoveHost("00:11:22:33:44:55")
		assert.Error(err)
		assert.Contains(err.Error(), "storage is readonly")
	})

	t.Run("SuccessfulRemove", func(t *testing.T) {
		assert := assert.New(t)

		s.readonly = false
		mockBackend.On("RemoveHost", "00:11:22:33:44:55").Return(nil)

		err := s.RemoveHost("00:11:22:33:44:55")
		assert.NoError(err, "Should remove host without error")
		mockBackend.AssertExpectations(t)
	})
}
