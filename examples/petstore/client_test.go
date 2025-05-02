package petstore

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

type PetStoreSuite struct {
	suite.Suite
	transport *httpmock.MockTransport
	client    ClientInterface
}

func (suite *PetStoreSuite) SetupTest() {
	var err error
	suite.transport = httpmock.NewMockTransport()
	suite.client, err = NewClient("http://invalid.localdomain/", WithRoundTripper(suite.transport))

	suite.Require().Nil(err)
}

func TestPetStoreSuite(t *testing.T) {
	suite.Run(t, new(PetStoreSuite))
}

func (suite *PetStoreSuite) TestAddOk() {
	suite.transport.RegisterResponder("POST", "/pet",
		httpmock.NewJsonResponderOrPanic(200, map[string]any{
			"id":        123,
			"name":      "Midnight",
			"photoUrls": nil,
			"status":    "sold",
		}),
	)

	var id int64 = 123
	status := PetStatus("sold")
	resp, err := suite.client.Pet().Add(context.Background(), &Pet{
		Category: &Category{},
		Id:       &id,
		Name:     "Midnight",
		Status:   &status,
	})
	suite.Nil(err)
	suite.NotNil(resp, "%s", resp)
	suite.EqualValues(*resp.Id, 123)
	suite.Equal(resp.Name, "Midnight")
}

func (suite *PetStoreSuite) TestAddResponseErrorHTML() {
	suite.transport.RegisterResponder(
		"POST",
		"/pet",
		httpmock.NewStringResponder(500, "<html><body><h1>Internal Server Error"),
	)

	resp, err := suite.client.Pet().AddResponse(context.Background(), &Pet{})
	suite.NotNil(resp)
	suite.ErrorContains(err, "server error '500 Internal Server Error'")
	target := &GeneralHTTPError{}
	suite.Require().ErrorAs(err, &target)

	suite.Equal("<html><body><h1>Internal Server Error", target.Text)
	suite.Nil(target.JSON)
}

func (suite *PetStoreSuite) TestAddResponseErrorJSON() {
	suite.transport.RegisterResponder("POST", "/pet",
		httpmock.NewJsonResponderOrPanic(404, map[string]any{
			"error": "uh-oh",
		}),
	)
	resp, err := suite.client.Pet().AddResponse(context.Background(), &Pet{})
	suite.NotNil(resp)

	suite.ErrorContains(err, "client error '404 Not Found'")
	target := &GeneralHTTPError{}
	suite.Require().ErrorAs(err, &target)
	suite.Equal("", target.Text)
	suite.Equal(map[string]any{"error": "uh-oh"}, target.JSON)
}

func (suite *PetStoreSuite) TestDelete() {
	suite.transport.RegisterResponder("DELETE", "/pet/12",
		httpmock.NewBytesResponder(200, nil),
	)
	apiKey := "abc"
	params := DeletePetParams{ApiKey: &apiKey}
	err := suite.client.Pet().Delete(context.Background(), 12, &params)
	suite.Nil(err)
}
