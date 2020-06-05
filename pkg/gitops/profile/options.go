package profile

import (
	"errors"
)

// Options holds options to configure the source of a Quick Start profile.
type Options struct {
	Name         string
	Overlay      string
	Revision     string
	ManifestOnly bool
}

// Validate validates this Options object.
func (o Options) Validate() error {
	if o.Name == "" {
		return errors.New("please supply a valid Quick Start profile name or URL")
	}
	return nil
}
