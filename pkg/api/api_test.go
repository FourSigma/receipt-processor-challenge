package api

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
)

const EXAMPLE1 = `
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },{
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },{
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },{
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },{
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}
`

const EXAMPLE2 = `
{
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-03-20",
  "purchaseTime": "14:33",
  "items": [
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    },{
      "shortDescription": "Gatorade",
      "price": "2.25"
    }
  ],
  "total": "9.00"
}
`

var pointsMap = map[string]int64{
	EXAMPLE1: 28,
	EXAMPLE2: 109,
}

func TestAPIHappyPath(t *testing.T) {
	api := New()

	for _, v := range []string{EXAMPLE1, EXAMPLE2} {

		example := strings.NewReader(v)
		req := httptest.NewRequest("POST", "/receipts/process", example)

		rec := httptest.NewRecorder()
		api.ProcessReceipt(rec, req)

		if rec.Code != 200 {
			t.Error("got", rec.Code, "want 200")
		}
		var response map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Errorf("got %v, want nil", err)
		}

		// Should have an id propterty
		val, ok := response["id"]
		if !ok {
			t.Errorf("got %v, want id", response)
		}

		// Should be a string
		_, ok = val.(string)
		if !ok {
			t.Errorf("got %v, want string", val)
		}

		// Should be a valid UUID
		_, err := uuid.Parse(fmt.Sprintf("%v", val))
		if err != nil {
			t.Errorf("got %v, want nil", err)
		}

		// Request the points
		recPoints := httptest.NewRecorder()
		reqPoints := httptest.NewRequest("GET", fmt.Sprintf("/receipts/%s/points", val), nil)
		reqPoints.SetPathValue("id", fmt.Sprintf("%v", val))

		api.GetReceipt(recPoints, reqPoints)
		if recPoints.Code != 200 {
			t.Error("got", recPoints.Code, "want 200")
		}

		responsePoints := map[string]int64{}
		if err := json.Unmarshal(recPoints.Body.Bytes(), &responsePoints); err != nil {
			t.Errorf("got %v, want nil", err)
		}

		gotPoints, ok := responsePoints["points"]
		if !ok {
			t.Errorf("got %v, want points key", responsePoints)
		}

		expectedPoints, ok := pointsMap[v]
		if !ok {
			t.Errorf("got %v, want example key", v)
		}

		if gotPoints != expectedPoints {
			t.Errorf("got %v, want %v", gotPoints, expectedPoints)
		}
	}
}

func TestAPISadPath(t *testing.T) {
	api := New()

	t.Run("Invalid JSON for POST /receipts/process", func(t *testing.T) {
		example := strings.NewReader(`{}`)
		req := httptest.NewRequest("POST", "/receipts/process", example)

		rec := httptest.NewRecorder()
		api.ProcessReceipt(rec, req)

		if rec.Code != 400 {
			t.Error("got", rec.Code, "want 400")
		}
	})

	t.Run("Invalid ID for GET /receipts/{id}/points", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/receipts/id/points", nil)
		req.SetPathValue("id", "232ajdf- 2323 322-322") // Invalid ID

		api.GetReceipt(rec, req)
		if rec.Code != 400 {
			t.Error("got", rec.Code, "want 400")
		}
	})
	t.Run("Not found ID for GET /receipts/{id}/points", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/receipts/id/points", nil)
		req.SetPathValue("id", "7fb1377b-b223-49d9-a31a-5a02701dd310") // Valid ID

		api.GetReceipt(rec, req)
		if rec.Code != 404 {
			t.Error("got", rec.Code, "want 404")
		}
	})
}
