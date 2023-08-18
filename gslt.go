package gslt

import (
	"encoding/json"

	steamworks "github.com/FlowingSPDG/steamworks-web-api-gen-go/generated"
)

type gameServerService struct {
	key string
	srv steamworks.IGameServersService
}

type GameServerService interface {
	GetAccountList() (*GetAccountListResponse, error)
	CreateAccount(appid uint32, memo string) (*CreateAccountResponse, error)
	SetMemo(steamid uint64, memo string) error
	ResetLoginToken(steamid uint64) (*ResetLoginTokenResponse, error)
	DeleteAccount(steamid uint64) error
	GetAccountPublicInfo(steamid uint64) (*GetAccountPublicInfoResponse, error)
	QueryLoginToken(loginToken string) (*QueryLoginTokenResponse, error)
	GetServerSteamIDsByIP(serverIP string) (*GetServerSteamIDsByIPResponse, error)
	GetServerIPsBySteamID(steamid uint64) (*GetServerIPsBySteamIDResponse, error)
}

func NewGameServerService(key string) GameServerService {
	return &gameServerService{
		key: key,
		srv: steamworks.NewIGameServersService(),
	}
}

func reMarshal(m map[string]any, dest any) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dest)
}

func (gs *gameServerService) GetAccountList() (*GetAccountListResponse, error) {
	resp, err := gs.srv.GetAccountListV1(steamworks.IGameServersServiceGetAccountListV1Input{
		Key: gs.key,
	})
	if err != nil {
		return nil, err
	}
	ret := &GetAccountListResponse{}
	if err := reMarshal(resp, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (gs *gameServerService) CreateAccount(appid uint32, memo string) (*CreateAccountResponse, error) {
	resp, err := gs.srv.CreateAccountV1(steamworks.IGameServersServiceCreateAccountV1Input{
		Key:   gs.key,
		Appid: appid,
		Memo:  memo,
	})
	if err != nil {
		return nil, err
	}
	ret := &CreateAccountResponse{}
	if err := reMarshal(resp, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (gs *gameServerService) SetMemo(steamid uint64, memo string) error {
	if _, err := gs.srv.SetMemoV1(steamworks.IGameServersServiceSetMemoV1Input{
		Key:     gs.key,
		Steamid: steamid,
		Memo:    memo,
	}); err != nil {
		return err
	}
	return nil
}

func (gs *gameServerService) ResetLoginToken(steamid uint64) (*ResetLoginTokenResponse, error) {
	resp, err := gs.srv.ResetLoginTokenV1(steamworks.IGameServersServiceResetLoginTokenV1Input{
		Key:     gs.key,
		Steamid: steamid,
	})
	if err != nil {
		return nil, err
	}
	ret := &ResetLoginTokenResponse{}
	if err := reMarshal(resp, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (gs *gameServerService) DeleteAccount(steamid uint64) error {
	if _, err := gs.srv.DeleteAccountV1(steamworks.IGameServersServiceDeleteAccountV1Input{
		Key:     gs.key,
		Steamid: steamid,
	}); err != nil {
		return err
	}

	return nil
}

func (gs *gameServerService) GetAccountPublicInfo(steamid uint64) (*GetAccountPublicInfoResponse, error) {
	resp, err := gs.srv.GetAccountPublicInfoV1(steamworks.IGameServersServiceGetAccountPublicInfoV1Input{
		Key:     gs.key,
		Steamid: steamid,
	})
	if err != nil {
		return nil, err
	}
	ret := &GetAccountPublicInfoResponse{}
	if err := reMarshal(resp, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (gs *gameServerService) QueryLoginToken(loginToken string) (*QueryLoginTokenResponse, error) {
	resp, err := gs.srv.QueryLoginTokenV1(steamworks.IGameServersServiceQueryLoginTokenV1Input{
		Key:        gs.key,
		LoginToken: loginToken,
	})
	if err != nil {
		return nil, err
	}
	ret := &QueryLoginTokenResponse{}
	if err := reMarshal(resp, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (gs *gameServerService) GetServerSteamIDsByIP(serverIPs string) (*GetServerSteamIDsByIPResponse, error) {
	resp, err := gs.srv.GetServerSteamIDsByIPV1(steamworks.IGameServersServiceGetServerSteamIDsByIPV1Input{
		Key:       gs.key,
		ServerIps: serverIPs,
	})
	if err != nil {
		return nil, err
	}
	ret := &GetServerSteamIDsByIPResponse{}
	if err := reMarshal(resp, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (gs *gameServerService) GetServerIPsBySteamID(steamids uint64) (*GetServerIPsBySteamIDResponse, error) {
	resp, err := gs.srv.GetServerIPsBySteamIDV1(steamworks.IGameServersServiceGetServerIPsBySteamIDV1Input{
		Key:            gs.key,
		ServerSteamids: steamids,
	})
	if err != nil {
		return nil, err
	}
	ret := &GetServerIPsBySteamIDResponse{}
	if err := reMarshal(resp, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
