package server

import (
	"context"
	"fmt"

	"github.com/pixel365/goepp/command"
	"github.com/pixel365/goepp/command/create"
	"github.com/pixel365/goepp/response"

	"github.com/pixel365/zoner/internal/model"
	"github.com/pixel365/zoner/internal/stringutils/password"

	"github.com/pixel365/zoner/epp/server/conn"
)

func handleCreate(
	ctx context.Context,
	connection *conn.Connection,
	cmd command.Commander,
	e *Epp,
) error {
	data, _ := cmd.(*create.Create)

	if data.Host != nil {
		errResponse := response.AnyError(2101, response.UnimplementedCommand)
		if err := connection.Write(ctx, errResponse, e.Metrics.IncBytes); err != nil {
			return fmt.Errorf("write error response for unimplemented command: %w", err)
		}
		return nil
	}

	switch {
	case data.Domain != nil:
		return createDomain(ctx, connection, e, *data.Domain)
	case data.Contact != nil:
		return createContact(ctx, connection, e, *data.Contact)
	default:
		return nil
	}
}

func createContact(
	ctx context.Context,
	connection *conn.Connection,
	e *Epp,
	data create.Contact,
) error {
	passwordHash, _ := password.Hash(data.AuthInfo.Password, password.DefaultParams)

	contact := model.ContactCreateInput{
		ContactID:    data.ID,
		Name:         "",
		Organization: "",
		Email:        data.Email,
		Voice:        data.Voice.Num,
		Fax:          data.Fax.Num,
		AuthInfoHash: passwordHash,
		Disclose:     nil,
		RegistrarID:  connection.UserId(),
	}

	if len(data.PostalInfo) > 0 {
		contact.PostalInfo = make([]model.ContactPostalFields, 0, len(data.PostalInfo))

		for i := range data.PostalInfo {
			info := model.ContactPostalFields{
				Typ:           string(data.PostalInfo[i].Type),
				PostalName:    data.PostalInfo[i].Name,
				PostalOrg:     data.PostalInfo[i].Org,
				PostalCode:    data.PostalInfo[i].Addr.Pc,
				City:          data.PostalInfo[i].Addr.City,
				Country:       data.PostalInfo[i].Addr.Cc,
				Streets:       data.PostalInfo[i].Addr.Street,
				StateProvince: data.PostalInfo[i].Addr.Sp,
			}

			contact.PostalInfo = append(contact.PostalInfo, info)
		}
	}

	if data.Disclose != nil && len(data.Disclose.Items) > 0 {
		disclose := &model.Disclose{
			Flag:   uint8(data.Disclose.Flag),
			Fields: make([]string, 0, len(data.Disclose.Items)),
		}
		di := map[string]struct{}{}
		for i := range data.Disclose.Items {
			di[data.Disclose.Items[i].Name] = struct{}{}
		}
		for k := range di {
			disclose.Fields = append(disclose.Fields, k)
		}

		contact.Disclose = disclose
	}

	return e.ContactService.Create(ctx, contact)
}

func createDomain(
	ctx context.Context,
	connection *conn.Connection,
	e *Epp,
	data create.Domain,
) error {
	passwordHash, _ := password.Hash(data.AuthInfo.Password, password.DefaultParams)

	domain := model.DomainCreateInput{
		Name:         data.Name,
		Punycode:     data.Punycode,
		RegistrarID:  connection.UserId(),
		AuthInfoHash: passwordHash,
		Period:       data.Period.Value,
		PeriodUnit:   string(data.Period.Unit),
	}

	return e.DomainService.Create(ctx, domain)
}
