package mongodbcontainer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestNewMongoDBContainer(t *testing.T) {
	t.Run("should succeed", func(t *testing.T) {
		t.Run("default options", func(t *testing.T) {
			req := ContainerRequest{}
			c, err := NewContainer(context.Background(), req)
			assert.Nil(t, err)
			assert.NotNil(t, c)
			assert.NotNil(t, c.Con)
			assert.Equal(t, defaultImage, c.req.Image)
			assert.Len(t, c.req.ExposedPorts, 1)
			assert.Equal(t, defaultMappedPort, c.req.ExposedPorts[0])
			assert.Equal(t, defaultUser, c.req.User)
			assert.Equal(t, defaultPassword, c.req.Password)
			assert.Equal(t, defaultUser, c.req.Env[env_MONGO_INITDB_ROOT_USERNAME])
			assert.Equal(t, defaultPassword, c.req.Env[env_MONGO_INITDB_ROOT_PASSWORD])
			assert.Equal(t, wait.ForLog(defaultReadyLog), c.req.WaitingFor)
		})

	})
}
