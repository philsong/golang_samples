package stock

import "testing"

func TestDeal(t *testing.T) {
	var Response string
	Response = Deal("000001")
	if Response != "" {
		t.Error("Expected 1.5, got ", Response)
	}
}
