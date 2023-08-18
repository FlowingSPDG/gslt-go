package gslt

import (
	"encoding/json"
	"fmt"
)

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
