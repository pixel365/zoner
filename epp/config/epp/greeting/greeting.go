package greeting

import (
	"errors"
	"fmt"
	"strings"
)

type ObjURI string
type Extension string
type Language string

func (o ObjURI) String() string {
	return string(o)
}

func (o ObjURI) Validate() error {
	return nil
}

func (e Extension) String() string {
	return string(e)
}

func (e Extension) Validate() error {
	return nil
}

func (l *Language) String() string {
	return string(*l)
}

func (l *Language) Validate() error {
	if l == nil {
		return errors.New("language is nil")
	}

	lang := strings.ToLower(string(*l))
	if len(lang) != 2 {
		return errors.New("language must be 2 chars")
	}

	*l = Language(lang)

	return nil
}

// Greeting https://datatracker.ietf.org/doc/html/rfc5730#section-2.4
type Greeting struct {
	ServerID   string      `yaml:"sv_id"`
	Versions   []string    `yaml:"versions"`
	Extensions []Extension `yaml:"extensions"`
	Languages  []Language  `yaml:"languages"`
	Dcp        *Dcp        `yaml:"dcp,omitempty"`
	ObjURI     []ObjURI    `yaml:"objURI,omitempty"`
}

func (g Greeting) Validate() error {
	if g.ServerID == "" {
		return errors.New("server id is empty")
	}

	if len(g.Versions) == 0 {
		return errors.New("versions is empty")
	}

	for i, l := range g.Languages {
		if err := l.Validate(); err != nil {
			return fmt.Errorf("languages[%d] validation error: %w", i, err)
		}
	}

	if g.Dcp != nil {
		if err := g.Dcp.Validate(); err != nil {
			return err
		}
	}

	for i, uri := range g.ObjURI {
		if err := uri.Validate(); err != nil {
			return fmt.Errorf("objURI[%d] validation error: %w", i, err)
		}
	}

	return nil
}
