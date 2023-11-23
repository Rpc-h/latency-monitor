# Usage

- Download bitwarden secret 'Ansible Vault Latency Monitor' and stored as `./.vault_pass` file
- `make encrypt env=staging`: Encrypt secret values
- `make decrypt env=staging`: Decrypt secret values
- `make init env=staging`: Initialized terraform state
- `make plan env=staging`: Create terraform plan of changes to be applied
- `make apply env=staging`: Apply planned changes into GCP infrastructure
- `make destroy env=staging`: Apply planned changes into GCP infrastructure

