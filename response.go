package gslt

type Server struct {
	SteamID     string `json:"steamid"`
	AppID       uint32 `json:"appid"`
	LoginToken  string `json:"login_token"`
	Memo        string `json:"memo"`
	IsDeleted   bool   `json:"is_deleted"`
	IsExpired   bool   `json:"is_expired"`
	RtLastLogon int    `json:"rt_last_logon"`
}

type GetAccountListResponse struct {
	Response struct {
		Servers        []Server `json:"servers"`
		IsBanned       bool     `json:"is_banned"`
		Expires        int      `json:"expires"`
		Actor          string   `json:"actor"`
		LastActionTime int      `json:"last_action_time"`
	} `json:"response"`
}

type CreateAccountResponse struct {
	Response struct {
		SteamID    string `json:"steamid"`
		LoginToken string `json:"login_token"`
	} `json:"response"`
}

type SetMemoResponse struct {
}

type ResetLoginTokenResponse struct {
	Response struct {
		LoginToken string `json:"login_token"`
	} `json:"response"`
}

type DeleteAccountResponse struct {
}

type GetAccountPublicInfoResponse struct {
	Response struct {
		SteamID string `json:"steamid"`
		AppID   uint32 `json:"appid"`
	} `json:"response"`
}

type QueryLoginTokenResponse struct {
}

type GetServerSteamIDsByIPResponse struct {
	// NOT TESTED YET
}

type GetServerIPsBySteamIDResponse struct {
	// NOT TESTED YET
}
