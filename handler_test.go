package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateOrder(t *testing.T) {
	catalog := []Product{
		{ID: 1, Name: "Raqueta", Price: 9000},
		{ID: 2, Name: "Pelotas", Price: 500},
	}

	tests := []struct {
		name      string
		requested []string
		wantTotal int
		wantCount int
	}{
		{"empty order", nil, 0, 0},
		{"single product", []string{"Pelotas"}, 500, 1},
		{"two products", []string{"Raqueta", "Pelotas"}, 9500, 2},
		{"duplicate counts twice", []string{"Pelotas", "Pelotas"}, 1000, 2},
		{"unknown products ignored", []string{"Inexistente"}, 0, 0},
		{"mixed known and unknown", []string{"Raqueta", "Inexistente"}, 9000, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			names, total := calculateOrder(tt.requested, catalog)
			if total != tt.wantTotal {
				t.Errorf("total = %d, want %d", total, tt.wantTotal)
			}
			if len(names) != tt.wantCount {
				t.Errorf("matched %d products, want %d", len(names), tt.wantCount)
			}
		})
	}
}

func doRequest(t *testing.T, router http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	t.Helper()
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, bytes.NewReader(body))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestListProductsEndpoint(t *testing.T) {
	router := newRouter()
	rec := doRequest(t, router, http.MethodGet, "/products", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var got []Product
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("response is not a product list: %v", err)
	}
	if len(got) == 0 {
		t.Error("product list is empty")
	}
}

func TestCreateOrderEndpoint(t *testing.T) {
	router := newRouter()

	body := map[string]any{"products": []string{"Pelotas de tenis Penn"}}
	raw, _ := json.Marshal(body)

	rec := doRequest(t, router, http.MethodPost, "/orders", raw)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var resp OrderResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("bad response: %v", err)
	}
	if resp.Total != 499 {
		t.Errorf("total = %d, want 499", resp.Total)
	}
	if len(resp.Products) != 1 {
		t.Errorf("products = %v, want 1 entry", resp.Products)
	}
}

func TestCreateOrderInvalidJSON(t *testing.T) {
	router := newRouter()
	rec := doRequest(t, router, http.MethodPost, "/orders", []byte("{not json"))

	// The original replied 200 with an empty order here — a bug.
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestCreateOrderEmptyList(t *testing.T) {
	router := newRouter()
	body, _ := json.Marshal(Order{Products: []string{}})
	rec := doRequest(t, router, http.MethodPost, "/orders", body)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

func TestHealthEndpoint(t *testing.T) {
	router := newRouter()
	rec := doRequest(t, router, http.MethodGet, "/health", nil)
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
}
