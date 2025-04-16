terraform plan -out terraform.tfplan
terraform show -json terraform.tfplan > tfplan.json
terraform apply -auto-approve
terraform show -json terraform.tfstate > tfstate.json
jq -r '.. | objects |  .address' tfplan.json | grep -v "^null$" > addresses
jq -r '.. | objects |  .address' tfstate.json | grep -v "^null$" >> addresses
