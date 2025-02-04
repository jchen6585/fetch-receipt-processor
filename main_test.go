package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jchen6585/fetch-receipt-processor/utils"
)

var (
	testDB = make(map[string]int)
	testId = "-1"
)

func TestRequests(t *testing.T) {
	testCases := []struct {
		name       string
		receipt    utils.Receipt
		httpMethod string
		want       int
	}{
		{
			name: "POST method",
			receipt: utils.Receipt{
				Retailer:     utils.ConvertStringToPointer("Target"),
				PurchaseDate: utils.ConvertStringToPointer("2022-01-01"),
				PurchaseTime: utils.ConvertStringToPointer("13:01"),
				Items: []utils.Item{
					{
						ShortDescription: utils.ConvertStringToPointer("Mountain Dew 12PK"),
						Price:            utils.ConvertStringToPointer("6.49"),
					}, {
						ShortDescription: utils.ConvertStringToPointer("Emils Cheese Pizza"),
						Price:            utils.ConvertStringToPointer("12.25"),
					}, {
						ShortDescription: utils.ConvertStringToPointer("Knorr Creamy Chicken"),
						Price:            utils.ConvertStringToPointer("1.26"),
					}, {
						ShortDescription: utils.ConvertStringToPointer("Doritos Nacho Cheese"),
						Price:            utils.ConvertStringToPointer("3.35"),
					}, {
						ShortDescription: utils.ConvertStringToPointer("   Klarbrunn 12-PK 12 FL OZ  "),
						Price:            utils.ConvertStringToPointer("12.00"),
					},
				},
				Total: utils.ConvertStringToPointer("35.35"),
			},
			httpMethod: http.MethodPost,
			want:       1,
		},
		{
			name: "POST method fail",
			receipt: utils.Receipt{
				Retailer: utils.ConvertStringToPointer("Target"),
				Items: []utils.Item{
					{
						ShortDescription: utils.ConvertStringToPointer("Mountain Dew 12PK"),
						Price:            utils.ConvertStringToPointer("6.49"),
					},
				},
				Total: utils.ConvertStringToPointer("35.35"),
			},
			httpMethod: http.MethodPost,
			want:       0,
		},
		{
			name:       "GET method",
			httpMethod: http.MethodGet,
			want:       28,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.httpMethod == http.MethodPost {
				body, _ := json.Marshal(testCase.receipt)
				request := httptest.NewRequest(testCase.httpMethod, "/receipts/process", bytes.NewReader(body))
				response := httptest.NewRecorder()
				postHandler(response, request)
				if testCase.want == 1 && len(response.Body.String()) <= testCase.want {
					t.Error("Want a valid ID, but instead got an empty response")
				} else if testCase.want == 0 && len(response.Body.String()) > testCase.want {
					t.Error("Expected no ID, but instead got an ID")
				}
				if len(response.Body.String()) != 0 {
					str := response.Body.String()
					str = str[1 : len(str)-1]
					testDB[str] = utils.CalculatePoints(testCase.receipt)
					testId = str
				}
			} else if testCase.httpMethod == http.MethodGet {
				url := "/receipts/" + testId + "/points"
				request := httptest.NewRequest(testCase.httpMethod, url, nil)
				response := httptest.NewRecorder()
				getHandler(response, request)
				got := testDB[testId]
				if got != testCase.want {
					t.Errorf("Wanted %d Got %d", testCase.want, got)
				}
			}
		})
	}
}
