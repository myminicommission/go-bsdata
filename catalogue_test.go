package bsdata_test

import (
	"testing"

	"github.com/myminicommission/go-bsdata"
)

func TestGetData(t *testing.T) {
	testRepo := "star-wars-legion"
	testTag := "1.7.0"
	catalogues, err := bsdata.GetData(testRepo, testTag)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if catalogues == nil {
		t.Error("catalogues are nil")
		t.FailNow()
	}

	if len(catalogues) == 0 {
		t.Error("catalogues collection is empty")
		t.FailNow()
	}

	t.Logf("Found %d cat files", len(catalogues))
}
