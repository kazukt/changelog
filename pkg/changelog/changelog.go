package changelog

import (
	"errors"
	"fmt"
	"io"
)

// Changelog represents a changelog in its entirety, containing all the
// versions that are tracked in the changelog.
type Changelog struct {
	Title    string
	Preamble string
	Versions []*Version
}

// NewChangelog creates a new Changelog.
func NewChangelog() *Changelog {
	return &Changelog{Versions: []*Version{}}
}

// Version gets the Version struct which matches the name.
// Returns nil if no version was found matching the given name.
func (c *Changelog) Version(name string) *Version {
	for _, v := range c.Versions {
		if v.Name == name {
			return v
		}
	}
	return nil
}

// AddItem add a message under the type of changes of Unreleased version.
func (c *Changelog) AddItem(ctype ChangeType, message string) {
	v := c.Version("Unreleased")
	if v == nil {
		v := &Version{Name: "Unreleased"}
		c.Versions = append(c.Versions, v)
	}
	cc := v.ChangeCollection(ctype)
	if cc == nil {
		cc = &ChangeCollection{Type: ctype, Items: []string{}}
		v.Changes = append(v.Changes, cc)
	}
	cc.Items = append([]string{message}, cc.Items...)
}

// Release transforms Unreleased into the fixed version.
func (c *Changelog) Release(name, date string) {
	v := c.Version("Unreleased")
	v.Name = name
	v.Date = date
	unrelease := &Version{Name: "Unreleased"}
	c.Versions = append([]*Version{unrelease}, c.Versions...)
}

// Version contains the data for the changes for a given version.
type Version struct {
	// Name is a release version in a semantic versioning format:
	// https://semver.org/
	Name string

	// Date is a release date and optional.
	// Acceptable formats:
	//     YYYY-MM-DD
	Date string

	Changes []*ChangeCollection
}

// ChangeCollection gets the change collection witch matches the type of changes.
// Returns nil if no type was found matching the given type.
func (v *Version) ChangeCollection(ctype ChangeType) *ChangeCollection {
	for _, c := range v.Changes {
		if c.Type == ctype {
			return c
		}
	}
	return nil
}

// NewVersion creates a new Version.
func NewVersion() *Version {
	return &Version{}
}

// ChangeCollection is a collection of change types.
type ChangeCollection struct {
	Type  ChangeType
	Items []string
}

// TODO: refactoring
func (c *Changelog) Write(w io.Writer) error {
	_, err := fmt.Fprintf(w, "# %s\n\n", c.Title)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "%s\n\n", c.Preamble)
	if err != nil {
		return err
	}
	for _, version := range c.Versions {
		var nameHeader string
		if version.Name == "" {
			return errors.New("version is empty")
		} else {
			nameHeader += "[" + version.Name + "]"
		}
		if version.Date != "" {
			nameHeader += fmt.Sprintf(" - %s", version.Date)
		}
		_, err = fmt.Fprintf(w, "## %s\n", nameHeader)
		if err != nil {
			return err
		}
		if len(version.Changes) == 0 {
			_, err = fmt.Fprintln(w)
			if err != nil {
				return err
			}
		}
		for _, cc := range version.Changes {
			_, err = fmt.Fprintf(w, "### %s\n", cc.Type)
			if err != nil {
				return err
			}
			for _, item := range cc.Items {
				_, err = fmt.Fprintf(w, "- %s\n", item)
				if err != nil {
					return err
				}
			}
			_, err = fmt.Fprintln(w)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
