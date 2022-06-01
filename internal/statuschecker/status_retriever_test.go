package statuschecker

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	statuschecker "github.com/sortednet/statuschecker/internal/statuschecker/statuschecker_test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestStatusReceiver_checkService(t *testing.T) {
	ctx := context.TODO()

	testCases := map[string]struct {
		serviceName      string
		serviceUrl       string
		mockExpectations func(m *statuschecker.MockHttpClient)
		expectedStatus   Status
	}{
		"up": {
			serviceName: "testSvcUp",
			serviceUrl:  "http://testservice.com",
			mockExpectations: func(m *statuschecker.MockHttpClient) {
				m.EXPECT().Get("http://testservice.com").Return(&http.Response{StatusCode: 200}, nil)
			},
			expectedStatus: Up,
		},
		"down": {
			serviceName: "testSvcDown",
			serviceUrl:  "http://badservice.com",
			mockExpectations: func(m *statuschecker.MockHttpClient) {
				m.EXPECT().Get("http://badservice.com").Return(nil, fmt.Errorf("Timeout"))
			},
			expectedStatus: Down,
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			httpClient := statuschecker.NewMockHttpClient(mockCtrl)
			tc.mockExpectations(httpClient)

			retriever := newStatusRetriever(httpClient, tc.serviceName, tc.serviceUrl)
			retriever.checkService(ctx)
			assert.Equal(t, tc.expectedStatus, retriever.getStatus())
		})
	}
}
