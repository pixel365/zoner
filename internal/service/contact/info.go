package contact

import (
	"context"
	"errors"
	"time"

	"github.com/pixel365/goepp/command/info"
	"github.com/pixel365/goepp/command/update"
	"github.com/pixel365/goepp/response"

	contactrepo "github.com/pixel365/zoner/internal/repository/contact"

	"github.com/pixel365/zoner/internal/model"
)

func (s *Service) Info(
	ctx context.Context,
	payload info.ContactInfo,
	registrarId int64,
	registrarName string,
) response.Marshaller {
	input := model.ContactInfoInput{
		ContactID:   payload.ID,
		RegistrarID: registrarId,
		Password:    "",
	}

	if payload.AuthInfo != nil {
		input.Password = payload.AuthInfo.Password
	}

	var resp response.Marshaller
	code := response.CodeCommandCompletedSuccessfully

	contactInfo, err := s.repo.Info(ctx, input)
	switch {
	case errors.Is(err, contactrepo.ErrNotFound):
		s.log.WithUserId(registrarId).Error("contact not found error", err)
		code = response.CodeObjectDoesNotExist
		return response.AnyError(code, code.String())
	case err != nil:
		s.log.WithUserId(registrarId).Error("contact info error", err)
		code = response.CodeCommandFailed
		return response.AnyError(code, code.String())
	}

	statuses := make([]update.ContactStatus, 0, len(contactInfo.Statuses))
	for i := range contactInfo.Statuses {
		statuses = append(statuses, update.ContactStatus{
			Value: contactInfo.Statuses[i].Status,
		})
	}

	postalInfo := handlePostalInfo(contactInfo)
	disclose := handleDisclose(contactInfo)

	data := ContactInfoResData{
		ID:         contactInfo.ContactID,
		Roid:       contactInfo.Roid,
		Statuses:   statuses,
		PostalInfo: postalInfo,
		Voice:      contactInfo.Voice,
		Fax:        contactInfo.Fax,
		Email:      contactInfo.Email,
		ClID:       registrarName,
		CrID:       valueOrEmpty(contactInfo.CreatedByClientID),
		CrDate:     formatTime(contactInfo.CreatedAt),
		UpID:       valueOrEmpty(contactInfo.UpdatedByClientID),
		UpDate:     formatTime(contactInfo.UpdatedAt),
		TrDate:     "",
		AuthInfo:   nil,
		Disclose:   disclose,
	}

	resp = response.NewResponse[ContactInfoResData, struct{}](code, code.String()).
		WithResData(data)

	return resp
}

func handleDisclose(contactInfo model.ContactInfo) *ContactInfoDisclose {
	var disclose *ContactInfoDisclose
	if contactInfo.Disclose != nil {
		disclose = &ContactInfoDisclose{
			Flag: contactInfo.Disclose.Flag,
		}

		for i := range contactInfo.Disclose.Fields {
			switch DiscloseField(contactInfo.Disclose.Fields[i]) {
			case DiscloseName:
				disclose.Name = &ContactInfoDisclosePostal{}
			case DiscloseOrg:
				disclose.Org = &ContactInfoDisclosePostal{}
			case DiscloseAddr:
				disclose.Addr = &ContactInfoDisclosePostal{}
			case DiscloseAddrInt:
				disclose.Addr = &ContactInfoDisclosePostal{Type: "int"}
			case DiscloseAddrLoc:
				disclose.Addr = &ContactInfoDisclosePostal{Type: "loc"}
			case DiscloseVoice:
				disclose.Voice = &struct{}{}
			case DiscloseFax:
				disclose.Fax = &struct{}{}
			case DiscloseEmail:
				disclose.Email = &struct{}{}
			}
		}
	}
	return disclose
}

func handlePostalInfo(contactInfo model.ContactInfo) []ContactInfoPostalInfo {
	postalInfo := make([]ContactInfoPostalInfo, 0, len(contactInfo.PostalInfo))
	for i := range contactInfo.PostalInfo {
		name := contactInfo.Name
		if contactInfo.PostalInfo[i].PostalName != nil &&
			*contactInfo.PostalInfo[i].PostalName != "" {
			name = *contactInfo.PostalInfo[i].PostalName
		}

		org := ""
		if contactInfo.PostalInfo[i].PostalOrg != nil {
			org = *contactInfo.PostalInfo[i].PostalOrg
		}

		addr := ContactInfoAddrData{
			Streets: append([]string(nil), contactInfo.PostalInfo[i].Streets...),
		}

		if contactInfo.PostalInfo[i].City != nil {
			addr.City = *contactInfo.PostalInfo[i].City
		}

		if contactInfo.PostalInfo[i].StateProvince != nil {
			addr.Sp = *contactInfo.PostalInfo[i].StateProvince
		}

		if contactInfo.PostalInfo[i].PostalCode != nil {
			addr.Pc = *contactInfo.PostalInfo[i].PostalCode
		}

		if contactInfo.PostalInfo[i].CountryCode != nil {
			addr.Cc = *contactInfo.PostalInfo[i].CountryCode
		}

		postalInfo = append(postalInfo, ContactInfoPostalInfo{
			Type: contactInfo.PostalInfo[i].Type,
			Name: name,
			Org:  org,
			Addr: addr,
		})
	}
	return postalInfo
}

func valueOrEmpty(v *string) string {
	if v == nil {
		return ""
	}

	return *v
}

func formatTime(v time.Time) string {
	if v.IsZero() {
		return ""
	}

	return v.UTC().Format(time.RFC3339)
}
