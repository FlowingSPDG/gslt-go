package gslt

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

// GSLT GSLT Struct.
type GSLT struct {
	Steamid     SteamID64 `json:"steamid"`
	Appid       uint32    `json:"appid"`
	LoginToken  string    `json:"login_token"`
	Memo        string    `json:"memo"`
	IsDeleted   bool      `json:"is_deleted"`
	IsExpired   bool      `json:"is_expired"`
	RtLastLogon int       `json:"rt_last_logon"`

	APIToken string `json:"-"`
}

// sGenerateGSLTJSON Response Struct for generating GSLT. private
type sGenerateGSLTJSON struct {
	sGenerateGSLT `json:"response"`
}

// sGenerateGSLT Response Struct for generating GSLT
type sGenerateGSLT struct {
	Steamid    SteamID64 `json:"steamid"`
	LoginToken string    `json:"login_token"`
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
	Steamid SteamID64 `json:"steamid"`
	Appid   uint32    `json:"appid"`
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

func buildURL(op string, qs map[string]string) (string, error) {
	var (
		serviceURL = "https://api.steampowered.com/IGameServersService/"
		version    = "v1"
	)
	u, err := url.Parse(serviceURL)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, op, version)
	q := u.Query()
	for k, v := range qs {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// Delete Deletes a persistent game server account
func (g *GSLT) Delete() error {
	return DeleteGSLTBySteamID(g.APIToken, g.Steamid)
}

// SetMemo This method changes the memo associated with the game server account. Memos do not affect the account in any way. The memo shows up in the GetAccountList response and serves only as a reminder of what the account is used for.
func (g *GSLT) SetMemo(memo string) error {
	err := SetMemo(g.APIToken, g.Steamid, memo)
	if err != nil {
		return err
	}
	g.Memo = memo
	return nil
}

// ResetLoginToken Generates a new login token for the specified game server
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

// GetAccountPublicInfo Gets public information about a given game server account
func (g *GSLT) GetAccountPublicInfo() (*PublicInfoJSON, error) {
	info, err := GetAccountPublicInfo(g.APIToken, g.Steamid)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// QueryLoginToken Queries the status of the specified token, which must be owned by you
func (g *GSLT) QueryLoginToken() (*QueryInfoJSON, error) {
	info, err := QueryLoginToken(g.APIToken, g.LoginToken)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// GetGSLT Generates GSLT. You can use Manager or this function
func GetGSLT(token string, memo string, appid uint32) (*GSLT, error) {
	appID := fmt.Sprintf("%d", appid) /// HACK
	u, err := buildURL("CreateAccount", map[string]string{
		"key":   token,
		"appid": appID,
		"memo":  memo,
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(u, "json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	jsonBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := sGenerateGSLTJSON{}

	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		return nil, err
	}
	gslt := &GSLT{
		Steamid:     data.sGenerateGSLT.Steamid,
		Appid:       appid,
		LoginToken:  data.sGenerateGSLT.LoginToken,
		Memo:        memo,
		IsDeleted:   false,
		IsExpired:   false,
		RtLastLogon: 0,
		APIToken:    token,
	}
	return gslt, nil
}

// DeleteGSLTByToken Lists and Deletes specified GSLT. Returns error if GSLT was invalid. High API Usage.
func DeleteGSLTByToken(token string, gslt string) error {
	var steamid SteamID64
	List, err := ListGSLT(token)
	if err != nil {
		return err
	}
	for i := 0; i < len(List); i++ {
		if List[i].LoginToken == gslt {
			steamid = List[i].Steamid
		}
	}
	if steamid == 0 {
		return fmt.Errorf("SteamID not found")
	}

	return DeleteGSLTBySteamID(token, steamid)
}

// DeleteGSLTBySteamID Deletes specified GSLT. Returns error if GSLT was invalid
func DeleteGSLTBySteamID(token string, steamid SteamID64) error {
	u, err := buildURL("DeleteAccount", map[string]string{
		"key":     token,
		"steamid": steamid.String(),
	})
	if err != nil {
		return err
	}
	postresp, err := http.Post(u, "json", nil)
	if err != nil {
		return err
	}
	defer postresp.Body.Close()
	return nil
}

// SetMemo Sets/Updates specified GSLT's memo. Returns error if GSLT was invalid
func SetMemo(token string, steamid SteamID64, memo string) error {
	u, err := buildURL("SetMemo", map[string]string{
		"key":     token,
		"steamid": steamid.String(),
		"memo":    memo,
	})
	if err != nil {
		return err
	}
	postresp, err := http.Post(u, "json", nil)
	if err != nil {
		return err
	}
	defer postresp.Body.Close()
	return nil
}

// ResetLoginTokenBySteamID Resets specified GSLT's token. Returns error if GSLT was invalid. login_token will not be changed if its never used
func ResetLoginTokenBySteamID(token string, steamid SteamID64) (string, error) {
	u, err := buildURL("ResetLoginToken", map[string]string{
		"key":     token,
		"steamid": steamid.String(),
	})
	if err != nil {
		return "", err
	}
	postresp, err := http.Post(u, "json", nil)
	if err != nil {
		return "", err
	}
	defer postresp.Body.Close()

	ResetjsonBytes, err := ioutil.ReadAll(postresp.Body)
	if err != nil {
		return "", err
	}
	ResetData := sResetGSLTJSON{}

	if err := json.Unmarshal(ResetjsonBytes, &ResetData); err != nil {
		return "", err
	}

	return ResetData.LoginToken, nil
}

// ListGSLT Get GSLT Account list
func ListGSLT(token string) ([]*GSLT, error) {
	u, err := buildURL("GetAccountList", map[string]string{
		"key": token,
	})
	if err != nil {
		return nil, err
	}
	listresp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer listresp.Body.Close()

	GetjsonBytes, err := ioutil.ReadAll(listresp.Body)
	if err != nil {
		return nil, err
	}
	GetData := sGSLTListJSON{}

	if err := json.Unmarshal(GetjsonBytes, &GetData); err != nil {
		fmt.Println("s:", string(GetjsonBytes))
		return nil, err
	}
	for _, v := range GetData.Servers {
		v.APIToken = token
	}
	return GetData.Servers, nil
}

// GetAccountPublicInfo Gets GSLT Account's public information
func GetAccountPublicInfo(token string, steamid SteamID64) (*PublicInfoJSON, error) {
	u, err := buildURL("GetAccountPublicInfo", map[string]string{
		"key":     token,
		"steamid": steamid.String(),
	})
	if err != nil {
		return nil, err
	}
	listresp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer listresp.Body.Close()

	PublicjsonBytes, err := ioutil.ReadAll(listresp.Body)
	if err != nil {
		return nil, err
	}
	PublicData := PublicInfoJSON{}

	if err := json.Unmarshal(PublicjsonBytes, &PublicData); err != nil {
		return nil, err
	}
	return &PublicData, nil
}

// QueryLoginToken Queries login token info
func QueryLoginToken(token string, logintoken string) (*QueryInfoJSON, error) {
	u, err := buildURL("QueryLoginToken", map[string]string{
		"key":         token,
		"login_token": logintoken,
	})
	if err != nil {
		return nil, err
	}
	listresp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer listresp.Body.Close()

	QueryjsonBytes, _ := ioutil.ReadAll(listresp.Body)
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

// SteamID64 is uint64 format
type SteamID64 uint64

// UnmarshalJSON Unmarshalizes SteamID64 to uint64
func (m *SteamID64) UnmarshalJSON(b []byte) error {
	var number json.Number
	if err := json.Unmarshal(b, &number); err != nil {
		return err
	}
	i, err := number.Int64()
	if err != nil {
		return err
	}
	*m = SteamID64(i)
	return nil
}

func (m SteamID64) String() string {
	return fmt.Sprintf("%d", m)
}
