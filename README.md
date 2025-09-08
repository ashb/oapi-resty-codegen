# oapi-resty-codegen

An OpenAPI client generator (built on top of [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)) that provides "literate" programming style clients

The main goals are:

- Namespace/group operations rather than have them all be in a single flat
  namespace.

  I use the first tag for operations as the name fo the group to put it in

- Automatic HTTP error code handling. net/http doesn't treat a 4xx or 5xx
  response as an error, which makes some sense for raw HTTP client, but as a
  REST API client, an http error should be reported as an error to the Go
  Code.

- Provide "just enough" support of OpenAPI 3.1.

  There are lots of changes in 3.1 that this won't support, but it does enough
  that most simple uses of 3.1 will be understood. This was driven by wanting to
  generate a client from Apache Airflow's 3.1/FastAPI-based oapi schema.  

## Example

The "canonical example" [Pet Store][petstore] creates the following interfaces:

```go
type ClientInterface interface {
	// Pet deals with all the pet endpoints
	Pet() PetClient
	// Store deals with all the store endpoints
	Store() StoreClient
	// User deals with all the user endpoints
	User() UserClient
}

type PetClient interface {
	// Add a new pet to the store.
	Add(ctx context.Context, body *Pet) (*Pet, error)
	// AddResponse is a lower level version of [Add] and provides access to the raw [resty.Response]
	AddResponse(ctx context.Context, body *Pet) (*resty.Response, error)

	// Update an existing pet by Id.
	Update(ctx context.Context, body *Pet) (*Pet, error)
	// UpdateResponse is a lower level version of [Update] and provides access to the raw [resty.Response]
	UpdateResponse(ctx context.Context, body *Pet) (*resty.Response, error)

	// Multiple status values can be provided with comma separated strings.
	FindByStatus(ctx context.Context, params *FindPetsByStatusParams) (*[]Pet, error)
	// FindByStatusResponse is a lower level version of [FindByStatus] and provides access to the raw [resty.Response]
	FindByStatusResponse(ctx context.Context, params *FindPetsByStatusParams) (*resty.Response, error)

	// Multiple tags can be provided with comma separated strings. Use tag1, tag2, tag3 for testing.
	FindByTags(ctx context.Context, params *FindPetsByTagsParams) (*[]Pet, error)
	// FindByTagsResponse is a lower level version of [FindByTags] and provides access to the raw [resty.Response]
	FindByTagsResponse(ctx context.Context, params *FindPetsByTagsParams) (*resty.Response, error)

	// Delete a pet.
	Delete(ctx context.Context, petId int64, params *DeletePetParams) error
	// DeleteResponse is a lower level version of [Delete] and provides access to the raw [resty.Response]
	DeleteResponse(ctx context.Context, petId int64, params *DeletePetParams) (*resty.Response, error)

	// Returns a single pet.
	GetById(ctx context.Context, petId int64) (*Pet, error)
	// GetByIdResponse is a lower level version of [GetById] and provides access to the raw [resty.Response]
	GetByIdResponse(ctx context.Context, petId int64) (*resty.Response, error)

    // ...
```

And an example of how this is used in a client:

```go
	var id int64 = 123
	status := PetStatus("sold")
	resp, err := suite.client.Pet().Add(context.Background(), &Pet{
		Category: &Category{},
		Id:       &id,
		Name:     "Midnight",
		Status:   &status,
	})
```

`err` will be non-nil if there is a network issue _or_ if the HTTP status code is `> 399`. You can see more examples of this in [`examples/petstore/client_test.go`](examples/petstore/client_test.go)

[petstore]: https://raw.githubusercontent.com/swagger-api/swagger-petstore/refs/tags/swagger-petstore-v3-1.0.26/src/main/resources/openapi.yaml
