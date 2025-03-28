package marionette_client

type Cookie struct {
	Secure   bool   `json:"secure,omitempty"`
	Expiry   uint   `json:"expiry,omitempty"`
	Domain   string `json:"domain,omitempty"`
	HttpOnly bool   `json:"httpOnly,omitempty"`
	Name     string `json:"name"`
	Path     string `json:"path,omitempty"`
	Value    string `json:"value"`
}
