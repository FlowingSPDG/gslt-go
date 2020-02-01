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
	Appid       uint32 `json:"appid"`
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

// sResetGSLTJSON Response Struct for generating GSLT. private
type sResetGSLTJSON struct {
	sResetGSLT `json:"response"`
}

// sResetGSLT Response Struct for generating GSLT
type sResetGSLT struct {
	LoginToken string `json:"login_token"`
}

// PublicInfoJSON Response Struct for generating GSLT. private
type PublicInfoJSON struct {
	sPublicInfo `json:"response"`
}

// sResetGSLT Response Struct for generating GSLT
type sPublicInfo struct {
	Steamid string `json:"steamid"`
	Appid   uint32 `json:"appid"`
}

// QueryInfoJSON Response Struct for generating GSLT. private
type QueryInfoJSON struct {
	sQueryInfo `json:"response"`
}

// sQueryInfo Response Struct for generating GSLT
type sQueryInfo struct {
	IsBanned bool   `json:"is_banned"`
	Expires  uint32 `json:"expires"`
	Steamid  string `json:"steamid"`
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
func (m *Manager) Generate(memo string, appid uint32) (*GSLT, error) {
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
	err := DeleteGSLTByToken(g.APIToken, g.LoginToken)
	if err != nil {
		return err
	}
	g = nil
	return nil
}

// SetMemo SetMemo Sets/updates memo
func (g *GSLT) SetMemo(memo string) error {
	err := SetMemo(g.APIToken, g.Steamid, memo)
	if err != nil {
		return err
	}
	g.Memo = memo
	return nil
}

// ResetLoginToken Resets login_token. login_token will not changed if its never used
func (g *GSLT) ResetLoginToken() error {
	token, err := ResetLoginTokenBySteamID(g.APIToken, g.Steamid)
	if err != nil {
		return err
	}
	if token != "" {
		g.LoginToken = token
	}
	return nil
}

func (g *GSLT) GetAccountPublicInfo() (*PublicInfoJSON, error) {
	info, err := GetAccountPublicInfo(g.APIToken, g.Steamid)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (g *GSLT) QueryLoginToken() (*QueryInfoJSON, error) {
	info, err := QueryLoginToken(g.APIToken, g.Steamid)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// GetGSLT Generates GSLT. You can use Manager or this function
func GetGSLT(token string, memo string, appid uint32) (GSLT, error) {
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

// DeleteGSLT for back compatibility...
func DeleteGSLT(token string, gslt string) error {
	return DeleteGSLTByToken(token, gslt)
}

// DeleteGSLT Deletes specified GSLT. Returns error if GSLT was invalid
func DeleteGSLTByToken(token string, gslt string) error {
	var steamid string = ""
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

// DeleteGSLTBySteamID Deletes specified GSLT. Returns error if GSLT was invalid
func DeleteGSLTBySteamID(token string, steamid string) error {
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

// SetMemo Sets/Updates specified GSLT's memo. Returns error if GSLT was invalid
func SetMemo(token string, steamid string, memo string) error {
	SetMemoURL := fmt.Sprintf("https://api.steampowered.com/IGameServersService/SetMemo/v1/?key=%s&steamid=%s&memo=%s", token, steamid, memo)
	postresp, err := http.Post(SetMemoURL, "json", nil)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer postresp.Body.Close()
	return nil
}

// ResetLoginToken Resets specified GSLT's token. Returns error if GSLT was invalid. login_token will not be changed if its never used
func ResetLoginTokenBySteamID(token string, steamid string) (string, error) {
	SetLoginTokenURL := fmt.Sprintf("https://api.steampowered.com/IGameServersService/ResetLoginToken/v1/?key=%s&steamid=%s", token, steamid)
	postresp, err := http.Post(SetLoginTokenURL, "json", nil)
	if err != nil {
		fmt.Printf("ERR : %v\n", err)
		return "", err
	}
	defer postresp.Body.Close()

	ResetByteArray, _ := ioutil.ReadAll(postresp.Body)
	ResetjsonBytes := ([]byte)(string(ResetByteArray))
	ResetData := sResetGSLTJSON{}

	if err := json.Unmarshal(ResetjsonBytes, &ResetData); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return "", err
	}
	if ResetData.LoginToken == "" {
		fmt.Printf("login_token not changed...\n")
		return "", nil
	}

	return ResetData.LoginToken, nil
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

// GetAccountPublicInfo Gets GSLT Account's public information
func GetAccountPublicInfo(token string, steamid string) (*PublicInfoJSON, error) {
	ListURL := fmt.Sprintf("https://api.steampowered.com/IGameServersService/GetAccountPublicInfo/v1/?key=%s&steamid=%s", token, steamid)
	listresp, err := http.Get(ListURL)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer listresp.Body.Close()

	PublicByteArray, _ := ioutil.ReadAll(listresp.Body)
	PublicjsonBytes := ([]byte)(string(PublicByteArray))
	PublicData := PublicInfoJSON{}

	if err := json.Unmarshal(PublicjsonBytes, &PublicData); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return nil, err
	}
	return &PublicData, nil
}

// QueryLoginToken Queries login token info
func QueryLoginToken(token string, logintoken string) (*QueryInfoJSON, error) {
	ListURL := fmt.Sprintf("https://api.steampowered.com/IGameServersService/QueryLoginToken/v1/?key=%s&login_token=%s", token, logintoken)
	listresp, err := http.Get(ListURL)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer listresp.Body.Close()

	QueryByteArray, _ := ioutil.ReadAll(listresp.Body)
	QueryjsonBytes := ([]byte)(string(QueryByteArray))
	QueryData := QueryInfoJSON{}

	if err := json.Unmarshal(QueryjsonBytes, &QueryData); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return nil, err
	}
	return &QueryData, nil
}

// SetBanStatus // only for partners...?
// GET https://api.steampowered.com/IGameServersService/GetServerSteamIDsByIP/v1/

// GetServerSteamIDsByIP // Gets a list of server SteamIDs given a list of IPs
// GetServerIPsBySteamID // Gets a list of server IP addresses given a list of SteamIDs
