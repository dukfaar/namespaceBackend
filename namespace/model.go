package namespace

import (
	"github.com/dukfaar/goUtils/relay"
	"github.com/globalsign/mgo/bson"
)

type Model struct {
	ID   bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string        `json:"name,omitempty" bson:"name,omitempty"`
}

type MutationInput struct {
	Name *string
}

var GraphQLType = `
type Namespace {
	_id: ID
	name: String
}
` +
	relay.GenerateConnectionTypes("Namespace")
