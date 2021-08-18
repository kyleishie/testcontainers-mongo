package testingdb

import (
	mc "github.com/kyleishie/testcontainers-mongo/pkg/mongodbcontainer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetup(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		testingDB, err := Setup(mc.ContainerRequest{})
		assert.Nil(t, err)
		assert.NotNil(t, testingDB)

		client, err := testingDB.container.NewClient()
		assert.Nil(t, err)
		assert.NotNil(t, client)

		err = testingDB.TearDown()
		assert.Nil(t, err)
	})
}
