package contact

import (
	"context"
	"errors"
	"time"

	"github.com/pixel365/goepp/command/create"
	"github.com/pixel365/goepp/response"

	contactrepo "github.com/pixel365/zoner/internal/repository/contact"

	"github.com/pixel365/zoner/internal/model"
	"github.com/pixel365/zoner/internal/stringutils/password"
)

func (s *Service) Create(
	ctx context.Context,
	payload create.Contact,
	registrarId int64,
) response.Marshaller {
	passwordHash, _ := password.Hash(payload.AuthInfo.Password, password.DefaultParams)

	contact := model.ContactCreateInput{
		ContactID:    payload.ID,
		Name:         "",
		Organization: "",
		Email:        payload.Email,
		Voice:        payload.Voice.Num,
		Fax:          payload.Fax.Num,
		AuthInfoHash: passwordHash,
		Disclose:     nil,
		RegistrarID:  registrarId,
	}

	if len(payload.PostalInfo) > 0 {
		contact.PostalInfo = make([]model.ContactPostalFields, 0, len(payload.PostalInfo))

		for i := range payload.PostalInfo {
			info := model.ContactPostalFields{
				Typ:           string(payload.PostalInfo[i].Type),
				PostalName:    payload.PostalInfo[i].Name,
				PostalOrg:     payload.PostalInfo[i].Org,
				PostalCode:    payload.PostalInfo[i].Addr.Pc,
				City:          payload.PostalInfo[i].Addr.City,
				Country:       payload.PostalInfo[i].Addr.Cc,
				Streets:       payload.PostalInfo[i].Addr.Street,
				StateProvince: payload.PostalInfo[i].Addr.Sp,
			}

			contact.PostalInfo = append(contact.PostalInfo, info)
		}
	}

	if payload.Disclose != nil && len(payload.Disclose.Items) > 0 {
		disclose := &model.Disclose{
			Flag:   uint8(payload.Disclose.Flag),
			Fields: make([]string, 0, len(payload.Disclose.Items)),
		}
		di := map[string]struct{}{}
		for i := range payload.Disclose.Items {
			di[payload.Disclose.Items[i].Name] = struct{}{}
		}
		for k := range di {
			disclose.Fields = append(disclose.Fields, k)
		}

		contact.Disclose = disclose
	}

	var resp response.Marshaller

	contactId, err := s.repo.Create(ctx, contact)

	switch {
	case errors.Is(err, contactrepo.ErrAlreadyExists):
		s.log.WithUserId(registrarId).Error("contact already exists error", err)
		return response.AnyError(
			response.CodeObjectExists,
			response.MessageForCode(response.CodeObjectExists),
		)
	case errors.Is(err, contactrepo.ErrValidation):
		s.log.WithUserId(registrarId).Error("contact validation error", err)
		return response.AnyError(
			response.CodeParameterValuePolicyError,
			response.MessageForCode(response.CodeParameterValuePolicyError),
		)
	case err != nil:
		s.log.WithUserId(registrarId).Error("contact create error", err)
		return response.AnyError(
			response.CodeCommandFailed,
			response.MessageForCode(response.CodeCommandFailed),
		)
	}

	data := ContactCreateResData{
		ID:     contact.ContactID,
		CRDate: time.Now().UTC().Format(time.RFC3339),
	}
	resp = response.NewResponse[ContactCreateResData, struct{}](
		response.CodeCommandCompletedSuccessfully,
		response.MessageForCode(response.CodeCommandCompletedSuccessfully),
	).
		WithResData(data)

	s.log.WithUserId(registrarId).Info("contact %s created", contactId)

	return resp
}
