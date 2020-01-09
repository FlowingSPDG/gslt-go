package gslt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Manager GSLT Manager.
type Manager struct {
	APIToken string
	Servers  []*GSLT
}

// GSLT GSLT Struct.
type GSLT struct {
	Steamid     string `json:"steamid"`
	Appid       int    `json:"appid"`
	LoginToken  string `json:"login_token"`
	Memo        string `json:"memo"`
	IsDeleted   bool   `json:"is_deleted"`
	IsExpired   bool   `json:"is_expired"`
	RtLastLogon int    `json:"rt_last_logon"`

	APIToken string `json:"-"`
}

// sGenerateGSLTJSON Response Struct for generating GSLT. private
type sGenerateGSLTJSON struct {
	sGenerateGSLT `json:"response"`
}

// sGenerateGSLT Response Struct for generating GSLT
type sGenerateGSLT struct {
	Steamid    string `json:"steamid"`
	LoginToken string `json:"login_token"`
}

// GSLTList Response Struct for getting GSLTs list
type sGSLTListJSON struct {
	sGSLTList `json:"response"`
}

type sGSLTList struct {
	Servers []*GSLT `json:"servers"`

	IsBanned       bool   `json:"is_banned"`
	Expires        int    `json:"expires"`
	Actor          string `json:"actor"`
	LastActionTime int    `json:"last_action_time"`
}

// GetList Get and refresh server list
func (m *Manager) GetList() ([]*GSLT, error) {
	ListPointer, err := ListGSLT(m.APIToken)
	if err != nil {
		return nil, err
	}
	m.Servers = ListPointer
	return m.Servers, nil
}

// Generate Generates GSLT Token.
func (m *Manager) Generate(memo string, appid int) (*GSLT, error) {
	if m.APIToken == "" {
		return nil, fmt.Errorf("API Token Empty")
	}
	gslt, err := GetGSLT(m.APIToken, memo, appid)
	if err != nil {
		return nil, err
	}
	gslt.APIToken = m.APIToken
	gslt.Appid = appid
	gslt.Memo = memo
	m.Servers = append(m.Servers, &gslt)
	return &gslt, nil
}

// Delete Deletes token
func (g *GSLT) Delete() error {
	err := DeleteGSLT(g.APIToken, g.LoginToken)
	if err != nil {
		return err
	}
	g = nil
	return nil
}

// GetGSLT Generates GSLT. You can use Manager or this function
func GetGSLT(token string, memo string, appid int) (GSLT, error) {
	GetURL := fmt.Sprintf("https://api.steampowered.com/IGameServersService/CreateAccount/v1/?key=%s&appid=%d&memo=%s", token, appid, memo)

	resp, err := http.Post(GetURL, "json", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)

	jsonBytes := ([]byte)(string(byteArray))
	data := sGenerateGSLTJSON{}
	gslt := GSLT{}

	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return gslt, err
	}
	gslt.Steamid = data.sGenerateGSLT.Steamid
	gslt.LoginToken = data.sGenerateGSLT.LoginToken
	gslt.APIToken = token
	gslt.Appid = appid
	gslt.Memo = memo
	return gslt, nil
}

// DeleteGSLT Deletes specified GSLT. Returns error if GSLT was invalid
func DeleteGSLT(token string, gslt string) error {
	steamid := ""
	List, err := ListGSLT(token)
	if err != nil {
		return err
	}
	for i := 0; i < len(List); i++ {
		if List[i].LoginToken == gslt {
			steamid = List[i].Steamid
		}
	}
	if steamid == "" {
		return fmt.Errorf("SteamID not found")
	}
	DeleteURL := fmt.Sprintf("https://api.steampowered.com/IGameServersService/DeleteAccount/v1/?key=%s&steamid=%s", token, steamid)
	postresp, err := http.Post(DeleteURL, "json", nil)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer postresp.Body.Close()
	return nil
}

// ListGSLT Get GSLT Account list
func ListGSLT(token string) ([]*GSLT, error) {
	ListURL := "https://api.steampowered.com/IGameServersService/GetAccountList/v1/?key=" + token
	listresp, err := http.Get(ListURL)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer listresp.Body.Close()

	GetByteArray, _ := ioutil.ReadAll(listresp.Body)
	GetjsonBytes := ([]byte)(string(GetByteArray))
	GetData := sGSLTListJSON{}

	if err := json.Unmarshal(GetjsonBytes, &GetData); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return nil, err
	}
	for i := 0; i < len(GetData.Servers); i++ {
		GetData.Servers[i].APIToken = token
	}
	return GetData.Servers, nil
}
