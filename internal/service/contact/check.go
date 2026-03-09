package contact

import (
	"context"

	"github.com/pixel365/goepp/command/check"
	"github.com/pixel365/goepp/response"

	"github.com/pixel365/zoner/internal/model"
)

func (s *Service) Check(
	ctx context.Context,
	payload check.ContactCheck,
	registrarId int64,
) response.Marshaller {
	//TODO: check current limits on the number of contacts requested

	identifiers := model.ContactsIdentifiersInput{
		Identifiers: payload.IDs,
		RegistrarID: registrarId,
	}

	var resp response.Marshaller
	code := response.CodeCommandCompletedSuccessfully

	checkedContacts, err := s.repo.Check(ctx, identifiers)
	if err != nil {
		s.log.WithUserId(registrarId).Error("contacts check error", err)
		code = response.CodeCommandFailed
		return response.AnyError(code, code.String())
	}

	inUse := "In Use"
	data := ContactsCheckResData{}
	for _, contact := range checkedContacts {
		var reason *string
		var available uint8
		if contact.Available {
			available = 1
		}

		if available == 0 {
			reason = &inUse
		}

		data.Contacts = append(data.Contacts, SingleCheckContact{
			ID: CheckContactID{
				Available: available,
				Value:     contact.ID,
			},
			Reason: reason,
		})
	}

	resp = response.NewResponse[ContactsCheckResData, struct{}](code, code.String()).
		WithResData(data)

	return resp
}
