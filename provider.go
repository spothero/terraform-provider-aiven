package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jelmersnoeck/aiven"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"email": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Aiven email address",
			},
			"otp": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Aiven One-Time password",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Aiven password",
			},
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Aiven Authentication Token",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"aiven_project":      resourceProject(),
			"aiven_service":      resourceService(),
			"aiven_database":     resourceDatabase(),
			"aiven_service_user": resourceServiceUser(),
		},

		ConfigureFunc: func(d *schema.ResourceData) (interface{}, error) {
			if d.Get("api_token") {
				return aiven.NewTokenClient(
					d.Get("api_token"),
				)
			}
			return aiven.NewMFAUserClient(
				d.Get("email").(string),
				d.Get("otp").(string),
				d.Get("password").(string),
			)
		},
	}
}

func optionalString(d *schema.ResourceData, key string) string {
	str, ok := d.Get(key).(string)
	if !ok {
		return ""
	}
	return str
}
