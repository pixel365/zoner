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
	var resp response.Marshaller

	data, _ := cmd.(*create.Create)

	switch {
	case data.Domain != nil:
		return createDomain(ctx, connection, e, *data.Domain)
	case data.Contact != nil:
		resp = e.ContactService.Create(ctx, *data.Contact, connection.UserId())
	default:
		resp = response.AnyError(2101, response.UnimplementedCommand)
	}

	if err := connection.Write(ctx, resp, e.Metrics.IncBytes); err != nil {
		return fmt.Errorf("write response error: %w", err)
	}

	return nil
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
