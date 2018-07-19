package main

import (
	"github.com/dukfaar/goUtils/relay"
	"github.com/dukfaar/namespaceBackend/namespace"
)

var Schema string = `
		schema {
			query: Query
			mutation: Mutation
		}

		type Query {
			namespaces(first: Int, last: Int, before: String, after: String): NamespaceConnection!
			namespace(id: ID!): Namespace!
		}

		input NamespaceMutationInput {
			name: String
		}

		type Mutation {
			createNamespace(input: NamespaceMutationInput): Namespace!
			updateNamespace(id: ID!, input: NamespaceMutationInput): Namespace!
			deleteNamespace(id: ID!): ID
		}` +
	relay.PageInfoGraphQLString +
	namespace.GraphQLType
