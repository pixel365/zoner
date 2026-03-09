package repository

import (
	"context"

	"github.com/pixel365/zoner/internal/model"
)

type ContactRepository interface {
	Create(context.Context, model.ContactCreateInput) (int64, error)
	Check(context.Context, model.ContactsIdentifiersInput) ([]model.CheckedContact, error)
}
