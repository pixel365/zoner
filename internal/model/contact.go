package model

type ContactPostalFields struct {
	Typ           string   `json:"-"`
	PostalName    string   `json:"-"`
	PostalOrg     string   `json:"-"`
	PostalCode    string   `json:"-"`
	City          string   `json:"-"`
	Country       string   `json:"-"`
	StateProvince string   `json:"-"`
	Streets       []string `json:"streets"`
}

type Disclose struct {
	Fields []string `json:"fields"`
	Flag   uint8    `json:"flag"`
}

type ContactCreateInput struct {
	Disclose     *Disclose `json:",inline"`
	ContactID    string
	Name         string
	Organization string
	Email        string
	Voice        string
	Fax          string
	AuthInfoHash string
	PostalInfo   []ContactPostalFields
	RegistrarID  int64
}

type ContactsIdentifiersInput struct {
	Identifiers []string
	RegistrarID int64
}

type CheckedContact struct {
	ID        string
	Available bool
}
