# Apollo Studio Plugin for Mach Composer

This plugin adds an integration for Apollo Studio for use with MACH composer.

This allows you to streamline your configuration and share it as an integration with your MACH components.

## Requirements

- [MACH Composer >=2.5](https://github.com/labd/mach-composer)
- [terraform-provider-apollo-studio](https://github.com/labd/terraform-provider-apollostudio)

## Usage

```yaml
mach_composer:
  plugins:
    apollostudio:
      source: mach-composer/apollostudio
      version: 0.1.0

global:
  # ...
  apollostudio:
    api_key: your-api-key
sites:
  - identifier: my-site
    # ...
    apollostudio:
      graph_ref: your-graph-ref
```
