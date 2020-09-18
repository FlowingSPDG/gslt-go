package gslt

import (
	"os"
	"testing"
)

var (
	SteamAPIToken = ""
)

func TestMain(m *testing.M) {
	SteamAPIToken = os.Getenv("STEAM_API_TOKEN")
	m.Run()
}

func TestGSLTStruct(t *testing.T) {
	Manager := Manager{}
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

func TestDeleteAllExpiredGSLT(t *testing.T) {
	Manager := Manager{}
	Manager.APIToken = SteamAPIToken
	Manager.GetList()
	for i := 0; i < len(Manager.Servers); i++ {
		if Manager.Servers[i].IsExpired == true {
			Manager.Servers[i].Delete()
		}
	}
}

func TestDeleteAllGSLT(t *testing.T) {
	Manager := Manager{}
	Manager.APIToken = SteamAPIToken
	Manager.GetList()
	for i := 0; i < len(Manager.Servers); i++ {
		Manager.Servers[i].Delete()
	}
}

func TestSetMemo(t *testing.T) {
	Manager := Manager{}
	Manager.APIToken = SteamAPIToken
	Manager.GetList()
	generated, err := Manager.Generate("GO_TEST_MEMO", 730)
	if err != nil {
		t.Errorf("Failed to generate GSLT...\nERR : %v\n", err)
		return
	}
	generated.SetMemo("GO_TEST_MEMO_SETMEMO")
}

func TestResetLoginToken(t *testing.T) {
	Manager := Manager{}
	Manager.APIToken = SteamAPIToken
	Manager.GetList()
	generated, err := Manager.Generate("GO_TEST_RESET", 730)
	if err != nil {
		t.Errorf("Failed to generate GSLT...\nERR : %v\n", err)
		return
	}
	t.Logf("GENERATED GSLT : %v\n", generated)
	err = generated.ResetLoginToken()
	if err != nil {
		t.Errorf("Failed to reset Login token...\nERR : %v\n", err)
		return
	}
	t.Logf("New login_token : %s\n", generated.LoginToken)
}

func TestGetAccountPublicInfo(t *testing.T) {
	Manager := Manager{}
	Manager.APIToken = SteamAPIToken
	Manager.GetList()
	generated, err := Manager.Generate("GO_TEST_PublicInfo", 730)
	if err != nil {
		t.Errorf("Failed to generate GSLT...\nERR : %v\n", err)
		return
	}
	t.Logf("GENERATED GSLT : %v\n", generated)
	info, err := generated.GetAccountPublicInfo()
	if err != nil {
		t.Errorf("Failed to reset Login token...\nERR : %v\n", err)
		return
	}
	t.Logf("Public Info : %v\n", info)
}

func TestQueryLoginToken(t *testing.T) {
	Manager := Manager{}
	Manager.APIToken = SteamAPIToken
	Manager.GetList()
	generated, err := Manager.Generate("GO_TEST_QueryLoginToken", 730)
	if err != nil {
		t.Errorf("Failed to generate GSLT...\nERR : %v\n", err)
		return
	}
	t.Logf("GENERATED GSLT : %v\n", generated)
	query, err := generated.QueryLoginToken()
	if err != nil {
		t.Errorf("Failed to reset Login token...\nERR : %v\n", err)
		return
	}
	t.Logf("Query Info : %v\n", query)
}
