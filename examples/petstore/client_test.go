package petstore

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

type DryRunTransport struct {
	Handler func(*http.Request) (*http.Response, error)
}

func (dr *DryRunTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return dr.Handler(r)
}

type PetStoreSuite struct {
	suite.Suite
	transport *DryRunTransport
	client    ClientInterface
}

func (suite *PetStoreSuite) SetupTest() {
	var err error
	suite.transport = &DryRunTransport{}
	suite.client, err = NewClient("http://invalid.localdomain/", WithRoundTripper(suite.transport))
	suite.Require().Nil(err)
}

func TestPetStoreSuite(t *testing.T) {
	suite.Run(t, new(PetStoreSuite))
}

func (suite *PetStoreSuite) TestAdd() {
	suite.transport.Handler = func(r *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(200, map[string]any{
			"id":        123,
			"name":      "Midnight",
			"photoUrls": nil,
			"status":    "sold",
		})
	}

	resp, err := suite.client.Pet().Add(context.Background())
	suite.Nil(err)
	suite.NotNil(resp)
	suite.EqualValues(*resp.Id, 123)
	suite.Equal(resp.Name, "Midnight")
}

func (suite *PetStoreSuite) TestDelete() {
	suite.transport.Handler = func(r *http.Request) (*http.Response, error) {
		return httpmock.NewBytesResponse(200, []byte{}), nil
	}
	err := suite.client.Pet().Delete(context.Background(), 0, nil)
	suite.Nil(err)
}
