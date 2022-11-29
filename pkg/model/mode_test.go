package model

import (
	"encoding/json"
	"testing"
)

func TestMilliSatoshiInt(t *testing.T) {
	millisat := ListFundsChannel{}
	if millisat.InternalTotAmountMsat != nil {
		t.Errorf("msat is not nil")
	}

	millisat.InternalTotAmountMsat = 1
	if millisat.InternalTotAmountMsat == nil {
		t.Errorf("msat is nil")
	}

	jsonObj, _ := json.Marshal(millisat)
	var newMillisat ListFundsChannel
	if err := json.Unmarshal(jsonObj, &newMillisat); err != nil {
		t.Errorf("%s", err)
	}

	if newMillisat.TotAmountMsat() != 1 {
		t.Errorf("expected `1` but received %d", newMillisat.TotAmountMsat())
	}
}

func TestMilliSatoshiStr(t *testing.T) {
	millisat := ListFundsChannel{}
	if millisat.InternalTotAmountMsat != nil {
		t.Errorf("msat is not nil")
	}

	millisat.InternalTotAmountMsat = "1msat"
	if millisat.InternalTotAmountMsat == nil {
		t.Errorf("msat is nil")
	}

	jsonObj, _ := json.Marshal(millisat)
	var newMillisat ListFundsChannel
	if err := json.Unmarshal(jsonObj, &newMillisat); err != nil {
		t.Errorf("%s", err)
	}

	if newMillisat.InternalTotAmountMsat != "1msat" {
		t.Errorf("expected `1msat` but received %s", newMillisat.InternalTotAmountMsat)
	}

	if newMillisat.TotAmountMsat() != 1 {
		t.Errorf("expected `1` but received %d", newMillisat.TotAmountMsat())
	}
}
