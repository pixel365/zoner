package greeting

import (
	"errors"
	"strings"
)

type Purpose string
type Recipient string
type Retention string
type Expiry string

func (p Purpose) String() string {
	return string(p)
}

func (r Recipient) String() string {
	return string(r)
}

func (r Retention) String() string {
	return string(r)
}

func (e Expiry) String() string {
	return string(e)
}

const (
	AdminPurpose   Purpose = "admin"
	ContactPurpose Purpose = "contact"
	ProvPurpose    Purpose = "prov"
	OtherPurpose   Purpose = "other"

	OtherRecipient     Recipient = "other"
	OursRecipient      Recipient = "ours"
	PublicRecipient    Recipient = "public"
	SameRecipient      Recipient = "same"
	UnrelatedRecipient Recipient = "unrelated"

	BusinessRetention   Retention = "business"
	IndefiniteRetention Retention = "indefinite"
	NoneRetention       Retention = "none"
	LegalRetention      Retention = "legal"
	StatedRetention     Retention = "stated"

	AbsoluteExpiry Expiry = "absolute"
	RelativeExpiry Expiry = "relative"
)

type Statement struct {
	Expiry    *Expiry     `yaml:"expiry,omitempty"`
	Retention string      `yaml:"retention"`
	Purpose   []Purpose   `yaml:"purpose"`
	Recipient []Recipient `yaml:"recipient"`
}

func (s *Statement) Validate() error {
	if s.Retention == "" {
		return errors.New("retention is empty")
	}

	if s.Expiry != nil && *s.Expiry == "" {
		return errors.New("expiry is empty")
	}

	if len(s.Purpose) == 0 {
		return errors.New("purpose is empty")
	}

	if len(s.Recipient) == 0 {
		return errors.New("recipient is empty")
	}

	purposeSet := make(map[Purpose]struct{})
	for _, p := range s.Purpose {
		purposeSet[p] = struct{}{}
	}

	recipientSet := make(map[Recipient]struct{})
	for _, r := range s.Recipient {
		recipientSet[r] = struct{}{}
	}

	s.Purpose = make([]Purpose, 0, len(purposeSet))
	for p := range purposeSet {
		s.Purpose = append(s.Purpose, p)
	}

	s.Recipient = make([]Recipient, 0, len(recipientSet))
	for r := range recipientSet {
		s.Recipient = append(s.Recipient, r)
	}

	return nil
}

func (s *Statement) WriteXML(b *strings.Builder) {
	b.WriteString(`<purpose>`)
	for _, p := range s.Purpose {
		b.WriteString(`<` + p.String() + `/>`)
	}
	b.WriteString(`</purpose>`)

	b.WriteString(`<recipient>`)
	for _, r := range s.Recipient {
		b.WriteString(`<` + r.String() + `/>`)
	}
	b.WriteString(`</recipient>`)

	b.WriteString(`<retention><` + s.Retention + `/></retention>`)

	if s.Expiry != nil {
		b.WriteString(`<expiry>` + s.Expiry.String() + `</expiry>`)
	}
}
