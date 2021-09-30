package artifactory

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
)

var goVirtualSchema = mergeSchema(baseVirtualRepoSchema, map[string]*schema.Schema{

	"external_dependencies_enabled": {
		Type:        schema.TypeBool,
		Computed:    true,
		Optional:    true,
		Description: "When set (default), Artifactory will automatically follow remote VCS roots in 'go-import' meta tags to download remote modules.",
	},
	"external_dependencies_patterns": {
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		RequiredWith: []string{"external_dependencies_enabled"},
		Description: "An allow list of Ant-style path patterns that determine which remote VCS roots Artifactory will " +
			"follow to download remote modules from, when presented with 'go-import' meta tags in the remote repository response. " +
			"By default, this is set to '**', which means that remote modules may be downloaded from any external VCS source.",
	},
})

func newGoVirtStruct() interface{} {
	return &services.GoVirtualRepositoryParams{}
}

var goVirtReader = mkVirtualRepoRead(packGoVirtualRepository, newGoVirtStruct)

func resourceArtifactoryGoVirtualRepository() *schema.Resource {
	return &schema.Resource{
		Create: mkVirtualCreate(unpackGoVirtualRepository, goVirtReader),
		Read:   goVirtReader,
		Update: mkVirtualUpdate(unpackGoVirtualRepository, goVirtReader),
		Delete: resourceVirtualRepositoryDelete,
		Exists: resourceVirtualRepositoryExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: goVirtualSchema,
	}
}

func unpackGoVirtualRepository(s *schema.ResourceData) (interface{}, string) {
	d := &ResourceData{s}
	base, _ := unpackBaseVirtRepo(s)

	repo := services.GoVirtualRepositoryParams{
		VirtualRepositoryBaseParams:  base,
		ExternalDependenciesPatterns: d.getList("external_dependencies_patterns"),
		ExternalDependenciesEnabled:  d.getBoolRef("external_dependencies_enabled", false),
	}

	return repo, repo.Key
}

func packGoVirtualRepository(r interface{}, d *schema.ResourceData) error {
	repo := r.(*services.GoVirtualRepositoryParams)
	setValue := packBaseVirtRepo(d, repo.VirtualRepositoryBaseParams)

	setValue("external_dependencies_patterns", repo.ExternalDependenciesPatterns)
	errors := setValue("external_dependencies_enabled", repo.ExternalDependenciesEnabled)

	if errors != nil && len(errors) > 0 {
		return fmt.Errorf("failed to pack go virtual repo %q", errors)
	}

	return nil
}