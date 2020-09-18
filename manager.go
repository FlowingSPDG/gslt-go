package gslt

import "fmt"

// Manager GSLT Manager.
type Manager struct {
	APIToken string
	Servers  []*GSLT
}

// NewManager Get new manager
func NewManager(token string) Manager {
	return Manager{
		APIToken: token,
		Servers:  make([]*GSLT, 0),
	}
}

// GetList Gets a list of game server accounts with their logon tokens
func (m *Manager) GetList() error {
	if m == nil {
		return fmt.Errorf("Nil manager")
	}
	ListPointer, err := ListGSLT(m.APIToken)
	if err != nil {
		return  err
	}
	m.Servers = ListPointer
	return nil
}

// Generate Creates a persistent game server account
func (m *Manager) Generate(memo string, appid uint32) (*GSLT, error) {
	if m == nil {
		return nil, fmt.Errorf("Nil manager")
	}
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
