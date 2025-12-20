package dto

type Keyboard struct {
	Id         uint   `json:"id"`
	KeycapType string `json:"keycap_type"`
	BaseType   string `json:"base_type"`
	SwitchType string `json:"switch_type"`
	Color      string `json:"color"`
}