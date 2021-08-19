package testingdb

import (
	"context"
	"github.com/google/uuid"
	mc "github.com/kyleishie/testcontainers-mongo/pkg/mongodbcontainer"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"testing"
)

// RunWithDatabase creates a test using the given database's name. The new test simply runs fn injecting db.
// Use this function if you need to customize the database with custom options.
func RunWithDatabase(t *testing.T, name string, db *mongo.Database, fn func(*testing.T, *mongo.Database)) {
	t.Run(name, func(t *testing.T) {
		fn(t, db)
	})
}

// TestingDB - An object to store the container reference.
type TestingDB struct {
	container *mc.Container
}

// NewTestDB provides an entry point for you to config the mongo container.
func NewTestDB(req mc.ContainerRequest) (*TestingDB, error) {
	c, err := mc.NewContainer(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return &TestingDB{container: c}, nil
}

// TearDown performs clean up operations such as stopping the container.
func (tdb *TestingDB) TearDown() error {
	return tdb.container.Stop()
}

// Run creates a test using the given name.
// A new mongo.Database will be created against the TestingDB.container using the given name.
// The new test simply runs fn injecting the new mongo.Database.
//
// Returns an error if a database cannot be created.
func (tdb *TestingDB) Run(t *testing.T, name string, fn func(*testing.T, *mongo.Database)) error {
	dbId := uuid.New().String()
	dbId = strings.ReplaceAll(dbId, "-", "")
	db, err := tdb.container.NewDatabase(dbId)
	if err != nil {
		return err
	}
	RunWithDatabase(t, name, db, fn)
	if err := db.Drop(context.Background()); err != nil {
		return err
	}
	return nil
}

// MustRun creates a test using the given name.
// A new mongo.Database will be created against the TestingDB.container using the given name.
// The new test simply runs fn injecting the new mongo.Database.
//
// If a database cannot be created, t.Fatal is called.
func (tdb *TestingDB) MustRun(t *testing.T, name string, fn func(*testing.T, *mongo.Database)) {
	if err := tdb.Run(t, name, fn); err != nil {
		t.Fatal(err)
	}
}

////// SINGLETON /////////

func Setup(req mc.ContainerRequest) error {
	c, err := mc.NewContainer(context.Background(), req)
	if err != nil {
		return err
	}
	_testDB = &TestingDB{container: c}
	return nil
}

var _testDB *TestingDB

func get_testDB() *TestingDB {
	if _testDB == nil {
		_ = Setup(mc.ContainerRequest{})
	}
	return _testDB
}

// TearDown performs clean up operations such as stopping the container.
func TearDown() error {
	return get_testDB().container.Stop()
}

// Run creates a test using the given name.
// A new mongo.Database will be created against the TestingDB.container using the given name.
// The new test simply runs fn injecting the new mongo.Database.
//
// Returns an error if a database cannot be created.
func Run(t *testing.T, name string, fn func(*testing.T, *mongo.Database)) error {
	dbId := uuid.New().String()
	dbId = strings.ReplaceAll(dbId, "-", "")
	db, err := get_testDB().container.NewDatabase(dbId)
	if err != nil {
		return err
	}
	RunWithDatabase(t, name, db, fn)
	if err := db.Drop(context.Background()); err != nil {
		return err
	}
	return nil
}

// MustRun creates a test using the given name.
// A new mongo.Database will be created against the TestingDB.container using the given name.
// The new test simply runs fn injecting the new mongo.Database.
//
// If a database cannot be created, t.Fatal is called.
func MustRun(t *testing.T, name string, fn func(*testing.T, *mongo.Database)) {
	if err := get_testDB().Run(t, name, fn); err != nil {
		t.Fatal(err)
	}
}
