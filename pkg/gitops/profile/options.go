package profile

import (
	"errors"
)

// Options holds options to configure the source of a profile.
type Options struct {
	Name         string
	Overlay      string
	Revision     string
	ManifestOnly bool
}

// Validate validates this Options object.
func (o Options) Validate() error {
	if o.Name == "" {
		return errors.New("please supply a valid profile name or URL")
	}
	return nil
}
