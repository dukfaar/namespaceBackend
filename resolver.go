package main

import (
	"context"

	"github.com/dukfaar/goUtils/relay"
	"github.com/dukfaar/namespaceBackend/namespace"
	graphql "github.com/graph-gophers/graphql-go"
)

type Resolver struct {
}

func (r *Resolver) Namespaces(ctx context.Context, args struct {
	First  *int32
	Last   *int32
	Before *string
	After  *string
}) (*namespace.ConnectionResolver, error) {
	namespaceService := ctx.Value("namespaceService").(namespace.Service)

	var totalChannel = make(chan int)
	go func() {
		var total, _ = namespaceService.Count()
		totalChannel <- total
	}()

	var namespacesChannel = make(chan []namespace.Model)
	go func() {
		result, _ := namespaceService.List(args.First, args.Last, args.Before, args.After)
		namespacesChannel <- result
	}()

	var (
		start string
		end   string
	)

	var namespaces = <-namespacesChannel

	if len(namespaces) == 0 {
		start, end = "", ""
	} else {
		start, end = namespaces[0].ID.Hex(), namespaces[len(namespaces)-1].ID.Hex()
	}

	hasPreviousPageChannel, hasNextPageChannel := relay.GetHasPreviousAndNextPage(len(namespaces), start, end, namespaceService)

	return &namespace.ConnectionResolver{
		Models: namespaces,
		ConnectionResolver: relay.ConnectionResolver{
			relay.Connection{
				Total:           int32(<-totalChannel),
				From:            start,
				To:              end,
				HasNextPage:     <-hasNextPageChannel,
				HasPreviousPage: <-hasPreviousPageChannel,
			},
		},
	}, nil
}

func setDataOnModel(model *namespace.Model, input *namespace.MutationInput) {
	model.Name = *input.Name
}

func (r *Resolver) CreateNamespace(ctx context.Context, args struct {
	Input *namespace.MutationInput
}) (*namespace.Resolver, error) {
	namespaceService := ctx.Value("namespaceService").(namespace.Service)

	inputModel := namespace.Model{}
	setDataOnModel(&inputModel, args.Input)

	newModel, err := namespaceService.Create(&inputModel)

	if err == nil {
		return &namespace.Resolver{
			Model: newModel,
		}, nil
	}

	return nil, err
}

func (r *Resolver) UpdateNamespace(ctx context.Context, args struct {
	Id    string
	Input *namespace.MutationInput
}) (*namespace.Resolver, error) {
	namespaceService := ctx.Value("namespaceService").(namespace.Service)

	inputModel := namespace.Model{}
	setDataOnModel(&inputModel, args.Input)

	newModel, err := namespaceService.Update(args.Id, &inputModel)

	if err == nil {
		return &namespace.Resolver{
			Model: newModel,
		}, nil
	}

	return nil, err
}

func (r *Resolver) DeleteNamespace(ctx context.Context, args struct {
	Id string
}) (*graphql.ID, error) {
	namespaceService := ctx.Value("namespaceService").(namespace.Service)

	deletedID, err := namespaceService.DeleteByID(args.Id)
	result := graphql.ID(deletedID)

	if err == nil {
		return &result, nil
	}

	return nil, err
}

func (r *Resolver) Namespace(ctx context.Context, args struct {
	Id string
}) (*namespace.Resolver, error) {
	namespaceService := ctx.Value("namespaceService").(namespace.Service)

	queryNamespace, err := namespaceService.FindByID(args.Id)

	if err == nil {
		return &namespace.Resolver{
			Model: queryNamespace,
		}, nil
	}

	return nil, err
}

func (r *Resolver) NamespaceByName(ctx context.Context, args struct {
	Name string
}) (*namespace.Resolver, error) {
	namespaceService := ctx.Value("namespaceService").(namespace.Service)

	queryNamespace, err := namespaceService.FindByName(args.Name)

	if err == nil {
		return &namespace.Resolver{
			Model: queryNamespace,
		}, nil
	}

	return nil, err
}
