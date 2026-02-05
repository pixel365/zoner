package login

import "errors"

type Options struct {
	Version string `xml:"version"`
	Lang    string `xml:"lang"`
}

func (o Options) Validate() error {
	//TODO: compare version and lang with greeting

	if o.Version == "" {
		return errors.New("version is empty")
	}

	if o.Lang == "" {
		return errors.New("lang is empty")
	}

	return nil
}
