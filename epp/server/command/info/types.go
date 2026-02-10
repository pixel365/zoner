package info

import (
	"errors"

	"github.com/pixel365/zoner/epp/server/command/internal"
)

type DomainInfo struct {
	AuthInfo *internal.AuthInfo `xml:"authInfo,omitempty"`
	internal.DomainRef
}

type ContactInfo struct {
	AuthInfo *internal.AuthInfo `xml:"authInfo,omitempty"`
	internal.ContactRef
}

type HostInfo struct {
	AuthInfo *internal.AuthInfo `xml:"authInfo,omitempty"`
	internal.HostRef
}

func (i *DomainInfo) Validate() error {
	if i.AuthInfo != nil {
		if i.AuthInfo.Password == "" {
			return errors.New("domain:authInfo/domain:pw is required if authInfo is present")
		}
	}

	return i.DomainRef.Validate()
}

func (i *ContactInfo) Validate() error {
	if i.AuthInfo != nil {
		if i.AuthInfo.Password == "" {
			return errors.New("contact:authInfo/contact:pw is required if authInfo is present")
		}
	}

	return i.ContactRef.Validate()
}

func (i *HostInfo) Validate() error {
	if i.AuthInfo != nil {
		if i.AuthInfo.Password == "" {
			return errors.New("host:authInfo/host:pw is required if authInfo is present")
		}
	}

	return i.HostRef.Validate()
}
