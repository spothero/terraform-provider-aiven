package main

import (
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spothero/aiven"
)

func resourceProjectVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectVpcCreate,
		Read:   resourceProjectVpcRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Target project",
				ForceNew:    true,
			},
			"cloud": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Target cloud.",
				ForceNew:    true,
			},
			"network_cidr": {
				Type:        schema.TypeString,
				Optional:    false,
				Description: "The network CIDR to associate with the VPC. One of 10.0.0.0/24, 10.30.0.0/24, 172.16.40.0/24, 192.168.0.0/24, 192.168.80.0/24",
				ForceNew:    true,
			},
		},
	}
}

func resourceProjectVpcCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*aiven.Client)
	vpc, err := client.Vpcs.Create(
		d.Get("project").(string),
		aiven.CreateVpcRequest{
			d.Get("cloud").(string),
			d.Get("network_cidr").(string),
			nil,
		})

	if err != nil {
		return err
	}

	d.SetId(vpc.ProjectVpcID)
	return nil
}

func resourceProjectVpcRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*aiven.Client)

	vpc, err := client.Vpcs.Get(d.Get("project").(string), d.Id())
	if err != nil {
		return err
	} else if vpc == nil {
		return errors.New("VPC not found")
	}

	d.Set("cloud", vpc.CloudName)
	d.Set("network_cidr", vpc.NetworkCidr)

	return nil
}
