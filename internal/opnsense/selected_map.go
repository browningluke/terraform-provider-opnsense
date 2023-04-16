package opnsense

import "encoding/json"

/*
	OPNsense responses to some queries with json data that looks like:
	"some_key" : {
		"K1": {
			"selected": 0,
			"value": "...",
		},
		"K2": {
			"selected": 1,
			"value": "...",
		},
	}

	This type allows the JSON library to unmarshal that map into a string containing only
	the key that is selected (i.e. "K2", in the example above).
*/

type SelectedMap string

func (s *SelectedMap) UnmarshalJSON(data []byte) error {
	var dataMap map[string]struct {
		Value    string `json:"value"`
		Selected int    `json:"selected"`
	}

	err := json.Unmarshal(data, &dataMap)
	if err != nil {
		return err
	}

	// Find selected element
	for k, v := range dataMap {
		if v.Selected == 1 {
			*s = SelectedMap(k)
		}
	}

	return nil
}

func (s *SelectedMap) String() string {
	return string(*s)
}
