package core

import (
	"time"

	"github.com/MustWin/baremetal-sdk-go"
	"github.com/MustWin/terraform-Oracle-BareMetal-Provider/client"
	"github.com/hashicorp/terraform/helper/schema"
)

type ShapeDatasourceCrud struct {
	D      *schema.ResourceData
	Client client.BareMetalClient
	Res    *baremetal.ListShapes
}

func (r *ShapeDatasourceCrud) Get() (e error) {
	compartmentID := r.D.Get("compartment_id").(string)
	opts := getCoreOptionsFromResourceData(r.D, "availability_domain", "image_id", "page", "limit")

	r.Res = &baremetal.ListShapes{
		Shapes: []baremetal.Shape{},
	}

	for {
		var list *baremetal.ListShapes
		if list, e = r.Client.ListShapes(compartmentID, opts...); e != nil {
			break
		}

		r.Res.Shapes = append(r.Res.Shapes, list.Shapes...)

		var hasNextPage bool
		if opts, hasNextPage = getOptionsWithNextPageID(list.NextPage, opts); !hasNextPage {
			break
		}
	}

	return
}

func (r *ShapeDatasourceCrud) SetData() {
	if r.Res != nil {
		// Important, if you don't have an ID, make one up for your datasource
		// or things will end in tears
		r.D.SetId(time.Now().UTC().String())
		shapes := []map[string]string{}
		for _, v := range r.Res.Shapes {
			shape := map[string]string{
				"name": v.Name,
			}
			shapes = append(shapes, shape)
		}
		r.D.Set("shapes", shapes)
	}
	return
}
