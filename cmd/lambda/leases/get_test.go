package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/Optum/Redbox/pkg/api/response"
	"github.com/Optum/Redbox/pkg/db"
	"github.com/Optum/Redbox/pkg/db/mocks"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
)

func TestGetLeaseByID(t *testing.T) {

	t.Run("When the invoking Call and there are no errors", func(t *testing.T) {
		expectdLease := &db.RedboxLease{
			ID:             "unique-id",
			AccountID:      "123456789",
			PrincipalID:    "test",
			LeaseStatus:    db.Active,
			LastModifiedOn: 1561149393,
		}
		expectedLeaseResponse := &response.LeaseResponse{
			ID:             "unique-id",
			AccountID:      "123456789",
			PrincipalID:    "test",
			LeaseStatus:    db.Active,
			LastModifiedOn: 1561149393,
		}
		mockDb := mocks.DBer{}
		mockDb.On("GetLeaseByID", "unique-id").Return(expectdLease, nil)
		mockRequest := events.APIGatewayProxyRequest{HTTPMethod: http.MethodGet, Path: "/leases/unique-id"}

		controller := getController{
			Dao: &mockDb,
		}

		actualResponse, err := controller.Call(context.TODO(), &mockRequest)
		require.Nil(t, err)

		parsedResponse := &response.LeaseResponse{}
		err = json.Unmarshal([]byte(actualResponse.Body), parsedResponse)
		require.Nil(t, err)

		require.Equal(t, expectedLeaseResponse, parsedResponse, "Returns a single lease.")
		require.Equal(t, actualResponse.StatusCode, 200, "Returns a 200.")
	})

	t.Run("When the query fails", func(t *testing.T) {
		expectedError := errors.New("Error")
		mockDb := mocks.DBer{}
		mockDb.On("GetLeaseByID", "unique-id").Return(nil, expectedError)
		mockRequest := events.APIGatewayProxyRequest{HTTPMethod: http.MethodGet, Path: "/accounts/unique-id"}

		controller := getController{
			Dao: &mockDb,
		}

		actualResponse, err := controller.Call(context.TODO(), &mockRequest)
		require.Nil(t, err)

		require.Equal(t, actualResponse.StatusCode, 500, "Returns a 500.")
		require.Equal(t, actualResponse.Body, "{\"error\":{\"code\":\"ServerError\",\"message\":\"Failed Get on Lease unique-id\"}}", "Returns an error response.")
	})

}
