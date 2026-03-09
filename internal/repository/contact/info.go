package contact

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/pixel365/zoner/internal/model"
)

func (r *Repository) Info(
	ctx context.Context,
	data model.ContactInfoInput,
) (model.ContactInfo, error) {
	var result model.ContactInfo
	var discloseRaw []byte
	var statusesRaw []byte
	var postalInfoRaw []byte

	sql := `
SELECT
      c.id,
      c.contact_id,
      c.roid,
      c.registrar_id,
      c.name,
      c.email,
      c.organization,
      c.voice,
      c.fax,
      c.auth_info_hash,
      c.disclose,
      c.created_at,
      c.updated_at,
      c.created_by_client_id,
      c.updated_by_client_id,

      COALESCE((
          SELECT jsonb_agg(
              jsonb_build_object(
                  'status', cs.status,
                  'source', cs.source,
                  'reason', cs.reason,
                  'created_at', cs.created_at,
                  'created_by_client', cs.created_by_client
              )
              ORDER BY cs.id
          )
          FROM contact_statuses cs
          WHERE cs.contact_id = c.id
      ), '[]'::jsonb) AS statuses,

      COALESCE((
          SELECT jsonb_agg(
              jsonb_build_object(
                  'type', pi.type,
                  'name', pi.name,
                  'postal_name', pi.postal_name,
                  'postal_org', pi.postal_org,
                  'postal_code', pi.postal_code,
                  'city', pi.city,
                  'country_code', pi.country_code,
                  'streets', pi.streets,
                  'state_province', pi.state_province
              )
              ORDER BY pi.id
          )
          FROM contacts_postal_info pi
          WHERE pi.contact_id = c.id
      ), '[]'::jsonb) AS postal_info

  FROM contacts c
  WHERE c.contact_id = $1
    AND c.registrar_id = $2
    AND c.deleted_at IS NULL
  LIMIT 1
`

	err := r.db.QueryRow(ctx, sql, data.ContactID, data.RegistrarID).Scan(
		&result.ID,
		&result.ContactID,
		&result.Roid,
		&result.RegistrarID,
		&result.Name,
		&result.Email,
		&result.Organization,
		&result.Voice,
		&result.Fax,
		&result.AuthInfoHash,
		&discloseRaw,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.CreatedByClientID,
		&result.UpdatedByClientID,
		&statusesRaw,
		&postalInfoRaw,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return model.ContactInfo{}, fmt.Errorf(
				"%w, contact with id %s not found",
				ErrNotFound,
				data.ContactID,
			)
		default:
			return model.ContactInfo{}, fmt.Errorf("%w, %w", ErrInternal, err)
		}
	}

	if len(discloseRaw) > 0 && string(discloseRaw) != "{}" {
		result.Disclose = &model.Disclose{}
		if err = json.Unmarshal(discloseRaw, result.Disclose); err != nil {
			return model.ContactInfo{}, fmt.Errorf("%w, unmarshal disclose: %w", ErrInternal, err)
		}
	}

	if len(statusesRaw) > 0 {
		if err = json.Unmarshal(statusesRaw, &result.Statuses); err != nil {
			return model.ContactInfo{}, fmt.Errorf("%w, unmarshal statuses: %w", ErrInternal, err)
		}
	}

	if len(postalInfoRaw) > 0 {
		if err = json.Unmarshal(postalInfoRaw, &result.PostalInfo); err != nil {
			return model.ContactInfo{}, fmt.Errorf(
				"%w, unmarshal postal info: %w",
				ErrInternal,
				err,
			)
		}
	}

	return result, nil
}
