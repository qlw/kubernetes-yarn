package flavors

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToFlavorListQuery() (string, error)
}

// ListOpts helps control the results returned by the List() function.
// For example, a flavor with a minDisk field of 10 will not be returned if you specify MinDisk set to 20.
// Typically, software will use the last ID of the previous call to List to set the Marker for the current call.
type ListOpts struct {

	// ChangesSince, if provided, instructs List to return only those things which have changed since the timestamp provided.
	ChangesSince string `q:"changes-since"`

	// MinDisk and MinRAM, if provided, elides flavors which do not meet your criteria.
	MinDisk int `q:"minDisk"`
	MinRAM  int `q:"minRam"`

	// Marker and Limit control paging.
	// Marker instructs List where to start listing from.
	Marker string `q:"marker"`

	// Limit instructs List to refrain from sending excessively large lists of flavors.
	Limit int `q:"limit"`
}

// ToFlavorListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToFlavorListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// ListDetail instructs OpenStack to provide a list of flavors.
// You may provide criteria by which List curtails its results for easier processing.
// See ListOpts for more details.
func ListDetail(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToFlavorListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	createPage := func(r pagination.PageResult) pagination.Page {
		return FlavorPage{pagination.LinkedPageBase{PageResult: r}}
	}

	return pagination.NewPager(client, url, createPage)
}

// Get instructs OpenStack to provide details on a single flavor, identified by its ID.
// Use ExtractFlavor to convert its result into a Flavor.
func Get(client *gophercloud.ServiceClient, id string) GetResult {
	var gr GetResult
	gr.Err = perigee.Get(getURL(client, id), perigee.Options{
		Results:     &gr.Body,
		MoreHeaders: client.AuthenticatedHeaders(),
	})
	return gr
}
