package whoishiring

import "types/location"

type WhoIsHiringType struct {
	Apply       string
	Address     string
	Company     string
	Description string
	Kind        string
	Location    location.Location
	Url         string
	Title       string
	Source      string
	Source_name string
	Remote      bool
	Time        uint64
	Tags        []string
}
