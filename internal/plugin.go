package internal

import (
	"fmt"

	"github.com/mach-composer/mach-composer-plugin-helpers/helpers"
	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"
	"github.com/mach-composer/mach-composer-plugin-sdk/schema"
	"github.com/mitchellh/mapstructure"
)

type ApollostudioPlugin struct {
	provider     string
	environment  string
	globalConfig *ApollostudioConfig
	siteConfigs  map[string]*ApollostudioConfig
	enabled      bool
}

func NewApollostudioPlugin() schema.MachComposerPlugin {
	state := &ApollostudioPlugin{
		provider:    "0.0.4", // Provider version of `labd/apollostudio`
		siteConfigs: map[string]*ApollostudioConfig{},
	}
	return plugin.NewPlugin(&schema.PluginSchema{
		Identifier:          "apollostudio",
		Configure:           state.Configure,
		IsEnabled:           func() bool { return state.enabled },
		GetValidationSchema: state.GetValidationSchema,

		SetGlobalConfig: state.SetGlobalConfig,
		SetSiteConfig:   state.SetSiteConfig,

		// Renders
		RenderTerraformProviders: state.RenderTerraformProviders,
		RenderTerraformResources: state.RenderTerraformResources,
		RenderTerraformComponent: state.RenderTerraformComponent,
	})
}

func (p *ApollostudioPlugin) Configure(environment string, provider string) error {
	p.environment = environment
	if provider != "" {
		p.provider = provider
	}
	return nil
}

func (p *ApollostudioPlugin) GetValidationSchema() (*schema.ValidationSchema, error) {
	result := getSchema()
	return result, nil
}

func (p *ApollostudioPlugin) SetGlobalConfig(data map[string]any) error {
	cfg := ApollostudioConfig{}

	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}
	p.globalConfig = &cfg
	p.enabled = true

	return nil
}

func (p *ApollostudioPlugin) SetSiteConfig(site string, data map[string]any) error {
	cfg := ApollostudioConfig{}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}
	p.siteConfigs[site] = &cfg
	p.enabled = true
	return nil
}

func (p *ApollostudioPlugin) getSiteConfig(site string) *ApollostudioConfig {
	result := &ApollostudioConfig{}
	if p.globalConfig != nil {
		result.ApiKey = p.globalConfig.ApiKey
	}

	cfg, ok := p.siteConfigs[site]
	if ok {
		if cfg.GraphRef != "" {
			result.GraphRef = cfg.GraphRef
		}
		if cfg.ApiKey != "" {
			result.ApiKey = cfg.ApiKey
		}
	}

	return result
}

func (p *ApollostudioPlugin) RenderTerraformStateBackend(site string) (string, error) {
	return "", nil
}

func (p *ApollostudioPlugin) RenderTerraformProviders(site string) (string, error) {
	cfg := p.getSiteConfig(site)

	if cfg == nil {
		return "", nil
	}

	result := fmt.Sprintf(`
		apollostudio = {
			source = "labd/apollostudio"
			version = "%s"
		}
	`, helpers.VersionConstraint(p.provider))

	return result, nil
}

func (p *ApollostudioPlugin) RenderTerraformResources(site string) (string, error) {
	cfg := p.getSiteConfig(site)

	if cfg == nil {
		return "", nil
	}

	template := `
		provider "apollostudio" {
			{{ renderProperty "api_key" .ApiKey }}
			{{ renderProperty "graph_ref" .GraphRef }}
		}
	`
	return helpers.RenderGoTemplate(template, cfg)
}

func (p *ApollostudioPlugin) RenderTerraformComponent(site string, component string) (*schema.ComponentSchema, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return nil, nil
	}

	template := ``

	vars, err := helpers.RenderGoTemplate(template, cfg)
	if err != nil {
		return nil, err
	}

	result := &schema.ComponentSchema{
		Variables: vars,
	}

	return result, nil
}
