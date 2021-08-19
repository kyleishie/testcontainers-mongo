package simplerepository

import (
	"context"
	"github.com/kyleishie/testcontainers-mongo/pkg/testingdb"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestRepository_CreateUser(t *testing.T) {
	testingdb.MustRun(t, "should succeed", func(t *testing.T, database *mongo.Database) {
		/// NewTestDB
		ctx := context.Background()
		repo := NewRepository(database)
		testUser := User{
			FirstName: "Hola",
			LastName:  "Mundo",
		}

		/// Run test
		id, err := repo.CreateUser(ctx, testUser)

		/// Make assertions
		assert.Nil(t, err)
		assert.NotEqual(t, primitive.ObjectID{}, id)

		/// Since we have a live, unique database we can query it state and make more in depth assertions.
		/// In this case we are simply checking that one document was inserted that meets the requirements we
		/// started with.  One could use findOne here with id, but I find this simpler and just as effective.
		count, countErr := database.Collection(usersCollection).CountDocuments(ctx, bson.D{
			{"_id", id},
			{"firstName", testUser.FirstName},
			{"lastName", testUser.LastName},
		})
		if countErr != nil {
			t.Error(countErr)
		}
		assert.True(t, count == 1)

		/// After this point database will automatically be destroyed
	})

	testingdb.MustRun(t, "should error about empty firstName", func(t *testing.T, database *mongo.Database) {
		ctx := context.Background()
		repo := NewRepository(database)
		testUser := User{
			FirstName: "",
			LastName:  "Mundo",
		}

		/// Run test
		id, err := repo.CreateUser(ctx, testUser)

		/// Make assertions
		assert.Equal(t, primitive.ObjectID{}, id)
		assert.NotNil(t, err)
		assert.Equal(t, errorFirstNameRequired, err)

		/// Count documents
		count, countErr := database.Collection(usersCollection).CountDocuments(ctx, bson.D{
			{"_id", id},
			{"firstName", testUser.FirstName},
			{"lastName", testUser.LastName},
		})
		if countErr != nil {
			t.Error(countErr)
		}
		assert.True(t, count == 0)
	})

}
