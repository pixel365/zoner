package server

import (
	"context"
	"fmt"

	"github.com/pixel365/goepp/command"
	"github.com/pixel365/goepp/command/create"

	"github.com/pixel365/zoner/epp/server/response"
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

	if data.Contact != nil || data.Host != nil {
		errResponse := response.AnyError(2101, response.UnimplementedCommand)
		if err := connection.Write(ctx, errResponse, e.Metrics.IncBytes); err != nil {
			return fmt.Errorf("write error response for unimplemented command: %w", err)
		}
		return nil
	}

	passwordHash, _ := password.Hash(data.Domain.AuthInfo.Password, password.DefaultParams)

	domain := model.DomainCreateInput{
		Name:         data.Domain.Name,
		Punycode:     data.Domain.Punycode,
		RegistrarID:  connection.UserId(),
		AuthInfoHash: passwordHash,
		Period:       data.Domain.Period.Value,
		PeriodUnit:   string(data.Domain.Period.Unit),
	}

	return e.DomainService.Create(ctx, domain)
}
