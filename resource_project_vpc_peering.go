package main

import (
	"errors"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/spothero/aiven"
)

func resourceProjectVpcPeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectVpcPeeringCreate,
		Read:   resourceProjectVpcPeeringRead,
		Delete: resourceProjectVpcPeeringDelete,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Target project",
				ForceNew:    true,
			},
			"aiven_vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Aiven VPC ID",
				ForceNew:    true,
			},
			"peer_account_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The account ID of the cloud to which the Aiven VPC will be peered.",
				ForceNew:    true,
			},
			"peer_vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VPC ID of the Cloud VPC to which the Aiven VPC will be peered.",
				ForceNew:    true,
			},
		},
	}
}

func resourceProjectVpcPeeringCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*aiven.Client)
	vpcPeeringConnection, err := client.VpcPeeringConnections.Create(
		d.Get("project").(string),
		d.Get("aiven_vpc_id").(string),
		aiven.CreatePeeringConnection{
			d.Get("peer_account_id").(string),
			d.Get("peer_vpc_id").(string),
		})

	if err != nil {
		return err
	}
	d.SetId(vpcPeeringConnection.PeerCloudAccount + "!" + vpcPeeringConnection.PeerVpc)
	return nil
}

func resourceProjectVpcPeeringRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*aiven.Client)

	peerSettings := strings.Split(d.Id(), ",")
	vpcPeeringConnection, err := client.VpcPeeringConnections.Get(d.Get("project").(string), d.Get("aiven_vpc_id").(string), peerSettings[0], peerSettings[1])
	if err != nil {
		return err
	} else if vpcPeeringConnection == nil {
		return errors.New("VPC Peering Connection not found")
	}

	d.Set("peer_account_id", vpcPeeringConnection.PeerCloudAccount)
	d.Set("peer_vpc_id", vpcPeeringConnection.PeerVpc)

	return nil
}

func resourceProjectVpcPeeringDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*aiven.Client)

	peerSettings := strings.Split(d.Id(), ",")
	return client.VpcPeeringConnections.Delete(
		d.Get("project").(string),
		d.Get("aiven_vpc_id").(string),
		peerSettings[0],
		peerSettings[1],
	)
}
