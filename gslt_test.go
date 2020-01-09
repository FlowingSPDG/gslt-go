package gslt_test

import (
	"github.com/FlowingSPDG/gslt-go"
	"testing"
)

const (
	SteamAPIToken = ""
)

func TestGSLTStruct(t *testing.T) {
	Manager := gslt.Manager{}
	Manager.APIToken = SteamAPIToken
	Manager.GetList()
	generated, err := Manager.Generate("GOTEST", 730)
	if err != nil {
		t.Errorf("Failed to generate GSLT...\nERR : %v\n", err)
		return
	} else {
		t.Logf("Generated GSLT : %v\n", generated)
	}
	err = generated.Delete()
	if err != nil {
		t.Errorf("Failed to delete GSLT...\n ERR : %v\n", err)
		return
	}
}

func TestDeleteAllGSLT(t *testing.T) {
	Manager := gslt.Manager{}
	Manager.APIToken = SteamAPIToken
	Manager.GetList()
	for i := 0; i < len(Manager.Servers); i++ {
		if Manager.Servers[i].IsExpired == true {
			Manager.Servers[i].Delete()
		}
	}
}
