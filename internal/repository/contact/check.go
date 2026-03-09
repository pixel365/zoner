package contact

import (
	"context"

	"github.com/pixel365/zoner/internal/model"
)

func (r *Repository) Check(
	ctx context.Context,
	data model.ContactsIdentifiersInput,
) ([]model.CheckedContact, error) {
	var result []model.CheckedContact

	sql := `
  WITH input AS (
      SELECT
          v.contact_id AS requested_id,
          lower(v.contact_id) AS normalized_id,
          v.ord
      FROM unnest($1::text[]) WITH ORDINALITY AS v(contact_id, ord)
  )
  SELECT
      input.requested_id,
      c.id IS NULL AS available
  FROM input
  LEFT JOIN contacts c
      ON lower(c.contact_id) = input.normalized_id
     AND c.registrar_id = $2
     AND c.deleted_at IS NULL
  ORDER BY input.ord
`

	rows, err := r.db.Query(ctx, sql, data.Identifiers, data.RegistrarID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		contact := model.CheckedContact{
			ID:        "",
			Available: false,
		}

		err = rows.Scan(&contact.ID, &contact.Available)
		if err != nil {
			return nil, err
		}

		result = append(result, contact)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}
