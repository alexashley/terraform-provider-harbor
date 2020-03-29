package provider

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/liatrio/terraform-provider-harbor/harbor"
)

func resourceRegistry() *schema.Resource {
	return &schema.Resource{
		Create: resourceRegistryCreate,
		Read:   resourceRegistryRead,
		Update: resourceRegistryUpdate,
		Delete: resourceRegistryDelete,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"google-gcr", "azure-acr", "jfrog-artifactory", "docker-hub", "huawei-SWR", "aws-ecr", "ali-acr", "quay-io", "helm-hub", "gitlab", "docker-registry", "harbor"},
					false,
				),
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1024), // Self Imposed Limit, Harbor doesn't seem to have a specific limit
			},
			"endpoint_url": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^https?://.*$`), "validation error: endpoint_url must begin with http:// or https://'"),
			},
			"access_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_secret": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"verify_remote_cert": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func mapRegistryToData(d *schema.ResourceData, registry *harbor.Registry) error {
	err := d.Set("type", registry.Type)
	if err != nil {
		return err
	}
	err = d.Set("name", registry.Name)
	if err != nil {
		return err
	}
	err = d.Set("description", registry.Description)
	if err != nil {
		return err
	}
	err = d.Set("endpoint_url", registry.URL)
	if err != nil {
		return err
	}
	err = d.Set("access_id", registry.Credential.AccessKey)
	if err != nil {
		return err
	}
	err = d.Set("verify_remote_cert", !registry.Insecure)
	if err != nil {
		return err
	}

	return nil
}

func mapDataToRegistry(d *schema.ResourceData, registry *harbor.Registry) error {
	registry.Type = d.Get("type").(string)
	registry.Name = d.Get("name").(string)
	registry.Description = d.Get("description").(string)
	registry.URL = d.Get("endpoint_url").(string)
	registry.Insecure = !d.Get("verify_remote_cert").(bool)
	if d.Get("access_id") != nil && d.Get("access_secret") != nil {
		registry.Credential = harbor.RegistryCredential{
			Type:         "basic",
			AccessKey:    d.Get("access_id").(string),
			AccessSecret: d.Get("access_secret").(string),
		}
	}
	return nil
}

func resourceRegistryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)
	registry, err := client.GetRegistry(d.Id())
	if err != nil {
		return handleNotFoundError(err, d)
	}

	return mapRegistryToData(d, registry)
}

func resourceRegistryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)
	registry := &harbor.Registry{}

	mapDataToRegistry(d, registry)

	location, err := client.NewRegistry(registry)
	if err != nil {
		return err
	}

	d.SetId(location)

	return resourceRegistryRead(d, meta)
}

func resourceRegistryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)
	registry := &harbor.Registry{}

	mapDataToRegistry(d, registry)

	err := client.UpdateRegistry(d.Id(), registry)
	if err != nil {
		return err
	}

	return resourceRegistryRead(d, meta)
}

func resourceRegistryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*harbor.Client)

	err := client.DeleteRegistry(d.Id())
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
