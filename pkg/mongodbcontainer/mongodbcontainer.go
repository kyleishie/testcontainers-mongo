package mongodbcontainer

import (
	"context"
	"fmt"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultUser       = "user"
	defaultPassword   = "password"
	defaultImage      = "mongo:latest"
	defaultMappedPort = "27017/tcp"
	defaultReadyLog   = "Waiting for connections"
)

const (
	env_MONGO_INITDB_ROOT_USERNAME = "MONGO_INITDB_ROOT_USERNAME"
	env_MONGO_INITDB_ROOT_PASSWORD = "MONGO_INITDB_ROOT_PASSWORD"
)

type ContainerRequest struct {
	tc.GenericContainerRequest
	User     string
	Password string
}

type Container struct {
	container tc.Container
	req       ContainerRequest
	ctx       context.Context
}

// NewContainer creates an instance testcontainers.Container configured to run mongo:tag.
// By default, i.e, passing an empty req, the container is configured with the mongo:latest image, port 27017,
// a default user "user" with password "password". To override any of these options see ContainerRequest.
func NewContainer(ctx context.Context, req ContainerRequest) (mc *Container, err error) {
	provider, err := req.ProviderType.GetProvider()
	if err != nil {
		return nil, err
	}

	req.Started = true
	req.ExposedPorts = []string{defaultMappedPort}

	if req.Image == "" {
		req.Image = defaultImage
	}

	if req.User == "" {
		req.User = defaultUser
	}

	if req.Password == "" {
		req.Password = defaultPassword
	}

	req.Env = map[string]string{
		env_MONGO_INITDB_ROOT_USERNAME: req.User,
		env_MONGO_INITDB_ROOT_PASSWORD: req.Password,
	}

	req.WaitingFor = wait.ForLog(defaultReadyLog)

	c, err := provider.CreateContainer(ctx, req.ContainerRequest)
	if err != nil {
		return nil, err
	}

	mc = &Container{
		container: c,
		req:       req,
		ctx:       ctx,
	}

	if req.Started {
		if err = c.Start(ctx); err != nil {
			return
		}
	}

	return
}

func (c *Container) Start() error {
	return c.container.Start(c.ctx)
}

func (c *Container) Stop() error {
	return c.container.Terminate(c.ctx)
}

// NewClient creates a new mongo client using the connection string of the Container.
// The connection is tested once before returning the new client.
func (c *Container) NewClient() (*mongo.Client, error) {

	host, err := c.container.Host(c.ctx)
	if err != nil {
		return nil, err
	}

	port, err := c.container.MappedPort(c.ctx, defaultMappedPort)
	if err != nil {
		return nil, err
	}

	clientOptions := options.Client().ApplyURI(fmt.Sprintf(
		"mongodb://%s:%s@%s:%d",
		c.req.User,
		c.req.Password,
		host,
		port.Int(),
	))

	client, err := mongo.Connect(c.ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// NewDatabase creates a new mongo client then a new database with then given name and options.
func (c *Container) NewDatabase(name string, opts ...*options.DatabaseOptions) (*mongo.Database, error) {
	client, err := c.NewClient()
	if err != nil {
		return nil, err
	}
	return client.Database(name, opts...), nil
}
