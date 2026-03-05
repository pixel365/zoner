package contact

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/pixel365/zoner/internal/db/postgres"
	"github.com/pixel365/zoner/internal/model"
)

func (r *Repository) Create(
	ctx context.Context,
	data model.ContactCreateInput,
) (int64, error) {
	var contactId int64
	err := postgres.Tx(ctx, r.db, pgx.Serializable,
		insertContact(ctx, data, &contactId),
		insertContactPostalInfo(ctx, data.PostalInfo, &contactId),
	)

	return contactId, err
}

func insertContact(
	ctx context.Context,
	data model.ContactCreateInput,
	contactId *int64,
) func(tx pgx.Tx) error {
	sql := `
INSERT INTO contacts (
                      contact_id,
                      registrar_id,
                      name,
                      email,
                      organization,
                      voice,
                      fax,
                      auth_info_hash,
                      disclose
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9::jsonb)
RETURNING id
`
	var disclose any = struct{}{}
	if data.Disclose != nil {
		disclose = data.Disclose
	}

	b, _ := json.Marshal(disclose)

	return func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, sql,
			data.ContactID, data.RegistrarID, data.Name, data.Email, data.Organization,
			data.Voice, data.Fax, data.AuthInfoHash, b).Scan(contactId)

		if err == nil {
			return nil
		}

		switch {
		case strings.Contains(err.Error(), "duplicate key value violates unique constraint"):
			return fmt.Errorf(
				"%w, contact with id %s already exists",
				ErrAlreadyExists,
				data.ContactID,
			)
		case strings.Contains(err.Error(), "violates foreign key constraint"):
			return fmt.Errorf(
				"%w, registrar with id %d does not exist",
				ErrValidation,
				data.RegistrarID,
			)
		case strings.Contains(err.Error(), "violates check constraint \"contacts_disclose_check\""):
			return fmt.Errorf("%w, invalid contact disclose data", ErrValidation)
		default:
			return fmt.Errorf("%w, %w", ErrInternal, err)
		}
	}
}

func insertContactPostalInfo(
	ctx context.Context,
	data []model.ContactPostalFields,
	contactId *int64,
) func(tx pgx.Tx) error {
	sql := `
INSERT INTO contacts_postal_info (
                                 contact_id,
                                 name,
                                 type,
                                 postal_name,
                                 postal_org,
                                 postal_code,
                                 city,
                                 country_code,
                                 streets,
                                 state_province
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`

	return func(tx pgx.Tx) error {
		if *contactId == 0 {
			return nil
		}

		for i := range data {
			_, err := tx.Exec(
				ctx,
				sql,
				*contactId,
				nil, //TODO: save name
				data[i].Typ,
				data[i].PostalName,
				data[i].PostalOrg,
				data[i].PostalCode,
				data[i].City,
				data[i].Country,
				data[i].Streets,
				data[i].StateProvince,
			)
			if err != nil {
				switch {
				case strings.Contains(err.Error(), "violates foreign key constraint"):
					return fmt.Errorf(
						"%w, contact with id %d does not exist",
						ErrValidation,
						*contactId,
					)
				default:
					return fmt.Errorf("%w, %w", ErrInternal, err)
				}
			}
		}
		return nil
	}
}
