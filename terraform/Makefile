.POSIX:

.PHONY: ensure-env-variable-is-set
ensure-env-variable-is-set:
ifeq ($(env),)
	echo "Parameter <env> is missing" >&2 && exit 1
endif

tee_log=tee terraform-log-`date +%Y%m%d-%H%M%S`.log

encrypt: ensure-env-variable-is-set
encrypt: ## Envrypt sensitive variables
ifeq (,$(wildcard ./.vault_pass))
	echo "[ERROR] Vault file Not exists. Download it from bitwarden 'Ansible Vault Latency Monitor'";
	exit 1
else
	ansible-vault encrypt --vault-password-file .vault_pass ./terraform-secret-$(env).tfvars
endif

decrypt: ensure-env-variable-is-set
decrypt: ## Envrypt sensitive variables
ifeq (,$(wildcard ./.vault_pass))
	echo "[ERROR] Vault file Not exists. Download it from bitwarden 'Ansible Vault Latency Monitor'";
	exit 1
else
	ansible-vault decrypt --vault-password-file .vault_pass ./terraform-secret-$(env).tfvars
endif


init: ensure-env-variable-is-set
init: ## Initialize backend
	terraform init -backend-config="bucket=terraform-latency-monitor-$(env)" -reconfigure | $(tee_log)

is_encrypted := $(shell grep -q ANSIBLE_VAULT terraform-secret-$(env).tfvars && echo "YES" || echo "NO")
plan: ensure-env-variable-is-set
plan: ## Plan changes
ifeq ($(is_encrypted),NO)
	terraform plan -out=tfplan-$(env) -var-file=terraform-$(env).tfvars -var-file=terraform-secret-$(env).tfvars | $(tee_log)
else
	echo "[ERROR] Secret values file './terraform-secret-$(env).tfvars' are not decrypted. Execute first 'make decrypt env=$(env)'";
	exit 1
endif

apply: ensure-env-variable-is-set
apply: ## Apply changes
	terraform apply tfplan-$(env) | $(tee_log)

destroy: ensure-env-variable-is-set
destroy: ## Destroy infrastructure
ifeq ($(is_encrypted),NO)
	terraform destroy -var-file=terraform-$(env).tfvars -var-file=terraform-secret-$(env).tfvars | $(tee_log)
else
	echo "[ERROR] Secret values file './terraform-secret-$(env).tfvars' are not decrypted. Execute first 'make decrypt env=$(env)'";
	exit 1
endif

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
