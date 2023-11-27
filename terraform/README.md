# Usage

- Download bitwarden secret 'Ansible Vault Latency Monitor' and stored as `./.vault_pass` file
- Download service account file from Bitwarden `RPCH Terraform Service Account`
- `export $(grep -v '^#' .envrc | xargs)`: Setup default environment variables
- `make env=staging encrypt`: Encrypt secret values
- `make env=staging decrypt`: Decrypt secret values
- `make env=staging init`: Initialized terraform state
- `make env=staging plan`: Create terraform plan of changes to be applied
- `make env=staging apply`: Apply planned changes into GCP infrastructure
- `make env=staging destroy`: Apply planned changes into GCP infrastructure

