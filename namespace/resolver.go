package namespace

import graphql "github.com/graph-gophers/graphql-go"

type Resolver struct {
	Model *Model
}

func (r *Resolver) ID() *graphql.ID {
	id := graphql.ID(r.Model.ID.Hex())
	return &id
}

func (r *Resolver) Name() *string {
	return &r.Model.Name
}
