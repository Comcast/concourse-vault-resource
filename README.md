# Vault Resource
Reads secrets from [Vault](https://www.vaultproject.io/). This resource supports [KV1](https://www.vaultproject.io/docs/secrets/kv/index.html#kv-version-1) and [KV2](https://www.vaultproject.io/docs/secrets/kv/index.html#kv-version-2) and can check for new versions or specific versions if using KV2.

## Source Configuration
* `vault_addr`: *Required.* The location of the Vault server. `https://vault.example.com:8200`.

* `vault_token`: *Required. if secret_id and role_id are not set* The token to use for authentication. `abc123f4k3T0k3n!&`.

* `vault_paths`: *Required.* A list of paths:version to secrets in vault. You can place this in the source configuration or you may pass it a parameter when fetching the resource. 

```yaml
vault_paths:
  path/to/secret: -1 # -1 means latest
  path/to/secret/w/version: 1 # grab version 1
```

*AppRole Authentication*
* `role_name`: *Optional.* If set, `vault_token` is required. Resource will use the `vault_token` and `role_name` to obtain a `role_id` and `secret_id` and use that to authenticate the approle.

* `role_id`: *Optional.* The role_id to authenticate with. Must be used with `secret_id`.

* `secret_id`: *Optional.* The secret_id to authenticate with. must be used with `role_id`


*General Parameters*
* `debug`: *Optional.* Print debug information. Will not expose secrets

* `format`: *Optional.* Choose output format of either `json` or `yaml`. Default: `json`

* `prefix`: *Optional.* Prepends a prefix to the secret key

* `retries`: *Optional.* The amount of retries. Default: 3
    
* `upcase`: *Optional.* Converts all secret keys to UPPERCASE

* `sanitize`: *Optional.* Converts dots and dashes in a secret key to underscores

* `vault_insecure`: *Optional.* Skips Vault SSL verification 

### Example
Resource configuration 

``` yaml
resource_types:
- name: vault
  type: docker-image
  source:
    repository: hub.example.com/foo/concourse-vault-resource
    tag: latest

resources:
- name: vault
  type: vault
  source:
    vault_addr: https://vault.example.com:8200
    vault_token: {{token}}
```

Resource configuration with AppRole 

``` yaml
resource_types:
- name: vault
  type: docker-image
  source:
    repository: hub.example.com/foo/concourse-vault-resource
    tag: latest

resources:
- name: vault
  type: vault
  source:
    vault_addr: https://vault.example.com:8200
    vault_token: {{token}}
    role_name: atu_vault-admins_approle
```

Resource configuration with AppRole using `role_id` and `secret_id` 

``` yaml
resource_types:
- name: vault
  type: docker-image
  source:
    repository: hub.example.com/foo/concourse-vault-resource
    tag: latest

resources:
- name: vault
  type: vault
  source:
    vault_addr: https://vault.example.com:8200
    role_id: 123456zzxROLE_IDjhdjkfafpfwefwa
    secret_id: faffdsfafdSECRET_IDdsfsdfadfd
```

Fetching secrets:

``` yaml
- get: vault
  params:
    vault_paths: 
      # KV1 Engine Test
      secret/foo: -1
      # KV2 Engine Test
      kv2/data/foo/bar: 2
```

## Behavior

### `check`: Check for new versions.

### `in`: Read secrets from Vault
Reads secrets from Vault and stores them in /opt/resource/secrets as JSON or YAML.
