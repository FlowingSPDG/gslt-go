package gslt

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	gameServerServiceEndpoint = "https://api.steampowered.com/IGameServersService"

	getAccountListPath    = "/GetAccountList/v1"
	createAccountPath     = "/CreateAccount/v1"
	setMemoPath           = "/SetMemo/v1"
	resetLoginToken       = "/ResetLoginToken/v1"
	deleteAccount         = "/DeleteAccount/v1"
	getAccountPublicInfo  = "/GetAccountPublicInfo/v1"
	queryLoginToken       = "/QueryLoginToken/v1"
	getServerSteamIDsByIP = "/GetServerSteamIDsByIP/v1"
	getServerIPsBySteamID = "/GetServerIPsBySteamID/v1"
)

func buildURL(key string, path string, queries map[string]string) string {
	// get endpoint
	basePath := gameServerServiceEndpoint
	if queries == nil {
		queries = map[string]string{}
	}
	queries["key"] = key

	// build query
	u, _ := url.ParseRequestURI(basePath)
	q := u.Query()
	for k, v := range queries {
		q.Add(k, v)
	}
	// Encode query
	u.RawQuery = q.Encode()
	u = u.JoinPath(path)

	return u.String()
}

func get(endpoint string, p any) error {
	resp, err := http.Get(endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if p == nil {
		return nil
	}

	return json.Unmarshal(b, p)
}

func post(endpoint string, p any) error {
	resp, err := http.Post(endpoint, "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if p == nil {
		return nil
	}

	return json.Unmarshal(b, p)
}

func GetAccontList(key string) (*GetAccountListResponse, error) {
	u := buildURL(key, getAccountListPath, nil)
	resp := &GetAccountListResponse{}
	if err := get(u, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func CreateAccount(key string, appid uint32, memo string) (*CreateAccountResponse, error) {
	u := buildURL(key, createAccountPath, map[string]string{
		"appid": fmt.Sprintf("%d", appid),
		"memo":  memo,
	})
	resp := &CreateAccountResponse{}
	if err := post(u, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func SetMemo(key string, steamid string, memo string) error {
	u := buildURL(key, setMemoPath, map[string]string{
		"steamid": steamid,
		"memo":    memo,
	})
	if err := post(u, nil); err != nil {
		return err
	}

	return nil
}

func ResetLoginToken(key string, steamid string) (*ResetLoginTokenResponse, error) {
	u := buildURL(key, resetLoginToken, map[string]string{
		"steamid": steamid,
	})
	resp := &ResetLoginTokenResponse{}
	if err := post(u, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func DeleteAccount(key string, steamid string) error {
	u := buildURL(key, deleteAccount, map[string]string{
		"steamid": steamid,
	})
	if err := post(u, nil); err != nil {
		return err
	}

	return nil
}

func GetAccountPublicInfo(key string, steamid string) (*GetAccountPublicInfoResponse, error) {
	u := buildURL(key, getAccountPublicInfo, map[string]string{
		"steamid": steamid,
	})
	resp := &GetAccountPublicInfoResponse{}
	if err := get(u, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func QueryLoginToken(key string, loginToken string) (*QueryLoginTokenResponse, error) {
	u := buildURL(key, queryLoginToken, map[string]string{
		"login_token": loginToken,
	})
	resp := &QueryLoginTokenResponse{}
	if err := get(u, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func GetServerSteamIDsByIP(key string, serverIP string) (*GetServerSteamIDsByIPResponse, error) {
	u := buildURL(key, getServerSteamIDsByIP, map[string]string{
		"server_ips": serverIP,
	})
	resp := &GetServerSteamIDsByIPResponse{}
	if err := get(u, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func GetServerIPsBySteamID(key string, steamid string) (*GetServerIPsBySteamIDResponse, error) {
	u := buildURL(key, getServerIPsBySteamID, map[string]string{
		"server_steamids": steamid,
	})
	resp := &GetServerIPsBySteamIDResponse{}
	if err := get(u, resp); err != nil {
		return nil, err
	}

	return resp, nil
}
