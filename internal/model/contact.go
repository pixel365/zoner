package model

import "time"

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
	Disclose      *Disclose `json:",inline"`
	ContactID     string
	Name          string
	Organization  string
	Email         string
	Voice         string
	Fax           string
	AuthInfoHash  string
	RegistrarName string
	PostalInfo    []ContactPostalFields
	RegistrarID   int64
}

type ContactsIdentifiersInput struct {
	Identifiers []string
	RegistrarID int64
}

type CheckedContact struct {
	ID        string
	Available bool
}

type ContactInfoInput struct {
	ContactID   string
	Password    string
	RegistrarID int64
}

type ContactInfoStatus struct {
	CreatedAt       time.Time `json:"created_at"`
	Reason          *string   `json:"reason"`
	CreatedByClient *string   `json:"created_by_client"`
	Status          string    `json:"status"`
	Source          string    `json:"source"`
}

type ContactInfoPostalInfo struct {
	Name          *string  `json:"name"`
	PostalName    *string  `json:"postal_name"`
	PostalOrg     *string  `json:"postal_org"`
	PostalCode    *string  `json:"postal_code"`
	City          *string  `json:"city"`
	CountryCode   *string  `json:"country_code"`
	StateProvince *string  `json:"state_province"`
	Type          string   `json:"type"`
	Streets       []string `json:"streets"`
}

type ContactInfo struct {
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Fax               *string
	Disclose          *Disclose
	UpdatedByClientID *string
	CreatedByClientID *string
	Organization      *string
	Voice             *string
	AuthInfoHash      string
	Roid              string
	ContactID         string
	Email             string
	Name              string
	Statuses          []ContactInfoStatus
	PostalInfo        []ContactInfoPostalInfo
	RegistrarID       int64
	ID                int64
}
