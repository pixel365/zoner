package renew

import (
	"errors"
	"time"

	normalizer "github.com/pixel365/domain-normalizer"

	"github.com/pixel365/zoner/epp/server/command/command"
	"github.com/pixel365/zoner/epp/server/command/internal"
)

type Renew struct {
	Domain *Domain
}

func (r *Renew) Name() command.CommandName {
	return command.Renew
}

func (r *Renew) NeedAuth() bool {
	return true
}

func (r *Renew) Validate() error {
	if r.Domain == nil {
		return errors.New("exactly one renew object must be present")
	}

	if r.Domain.Name == "" {
		return errors.New("domain:name is required")
	}

	name, err := normalizer.Parse(r.Domain.Name)
	if err != nil {
		return errors.New("domain:name is invalid")
	}

	r.Domain.Name = name.Normalized

	if r.Domain.Period != nil {
		if r.Domain.Period.Value <= 0 {
			return errors.New("domain:period must be > 0")
		}

		switch r.Domain.Period.Unit {
		case internal.PeriodUnitYear, internal.PeriodUnitMonth:
		default:
			return errors.New("domain:period unit must be y or m")
		}
	}

	return r.validateDate()
}

func (r *Renew) validateDate() error {
	if r.Domain.CurrentExpirationDate == "" {
		return errors.New("domain:curExpDate is required")
	}

	t, err := time.Parse(time.DateOnly, r.Domain.CurrentExpirationDate)
	if err != nil {
		return errors.New("domain:curExpDate must be YYYY-MM-DD")
	}

	if t.Format(time.DateOnly) != r.Domain.CurrentExpirationDate {
		return errors.New("domain:curExpDate must be YYYY-MM-DD")
	}

	return nil
}
