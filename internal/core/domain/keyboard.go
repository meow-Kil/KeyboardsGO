package domain

type Keyboard struct{
	Id uint
	KeycapType string
	BaseType string
	SwitchType string
	Color string
}

type User struct {
    ID       uint
    Login    string
    Password string
    IsAdmin  bool
}
