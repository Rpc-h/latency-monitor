# Usage


export $(grep -v '^#' .envrc | xargs)
terraform init -backend-config="bucket=latency-monitor-terraform"
terraform plan -out=tfplan
terraform apply tfplan

## Resources

- Cloud Run

## Inputs

- Input 1
- Input 2

## Outputs

- Output 1