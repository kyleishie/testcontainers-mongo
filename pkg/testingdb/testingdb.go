package testingdb

import (
	"context"
	mc "github.com/kyleishie/testcontainers-mongo/pkg/mongodbcontainer"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

// RunWithDatabase creates a test using the given database's name. The new test simply runs fn injecting db.
// Use this function if you need to customize the database with custom options.
func RunWithDatabase(t *testing.T, db *mongo.Database, fn func(*testing.T, *mongo.Database)) {
	t.Run(db.Name(), func(t *testing.T) {
		fn(t, db)
	})
}

// _testingDB - An object to store the container reference.
type _testingDB struct {
	container *mc.Container
}

// Setup provides an entry point for you to config the mongo container.
func Setup(req mc.ContainerRequest) (*_testingDB, error) {
	c, err := mc.NewContainer(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return &_testingDB{container: c}, nil
}

// TearDown performs clean up operations such as stopping the container.
func (tdb *_testingDB) TearDown() error {
	return tdb.container.Stop()
}

// Run creates a test using the given name.
// A new mongo.Database will be created against the _testingDB.container using the given name.
// The new test simply runs fn injecting the new mongo.Database.
//
// Returns an error if a database cannot be created.
func (tdb *_testingDB) Run(t *testing.T, name string, fn func(*testing.T, *mongo.Database)) error {
	db, err := tdb.container.NewDatabase(name)
	if err != nil {
		return err
	}
	RunWithDatabase(t, db, fn)
	if err := db.Drop(context.Background()); err != nil {
		return err
	}
	return nil
}

// MustRun creates a test using the given name.
// A new mongo.Database will be created against the _testingDB.container using the given name.
// The new test simply runs fn injecting the new mongo.Database.
//
// If a database cannot be created, t.Fatal is called.
func (tdb *_testingDB) MustRun(t *testing.T, name string, fn func(*testing.T, *mongo.Database)) {
	if err := tdb.Run(t, name, fn); err != nil {
		t.Fatal(err)
	}
}
