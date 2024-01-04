# Run Terraform

## Setting up Azure access
First you will need to login using AzureCli command: 

``` bash
az login
#set your subscription 
az account list --all --output table
az account set --subscription {id}
```

## Terraform remote state
In terraform.tf change Storage account and container which must exist before running init command 
terraform init -upgrade -migrate-state

## Terraform commands
``` bash
terraform init -upgrade -migrate-state #Initialization of providers
### update terraform.tf-> backend "azurerm" block
terraform plan -var-file="environments/coffeemachine-{environment}.tfvars" -out coffeemachine-{environment}-{clientName}.tfplan # create a plan for terraform deployment
terraform apply coffeemachine-{environment}-{clientName}.tfplan #run terraform deployment
### enter variables needed
terraform destroy #remove terraform deployment
```

# Modify Terraform

## Setting up Azure access
First you will need to login using AzureCli command: 

``` bash
az login
#set your subscription 
az account list --all --output table
az account set --subscription {id}
```

## State
### Local state
In terraform.tf comment backend section so that state will be saved locally in terraform folder and run:
terraform init -upgrade
### Remote state
In terraform.tf change backend section azurerm->Storage account and container which must exist before running init command 
terraform init -upgrade -migrate-state

## Terraform Code modifications
To modify terraform code, modify file's you need and separate if possible by sections.
After modifications run this commands to validate and run terraform with local state. 
```bash
terraform fmt # format terraform files
terraform validate # validate that there are no errors
terraform init -upgrade #Initialization of providers
terraform plan -var-file="environments/coffeemachine.tfvars" -out coffeemachine-{clientName}.tfplan # create a plan for terraform deployment
terraform apply coffeemachine-{clientName}.tfplan #run terraform deployment
terraform destroy #remove terraform deployment
```

If there is some resource that is stuck at upgrading, you will have to replace VM, use with care:
```bash
terraform plan -replace azurerm_linux_virtual_machine.vms[0] -out terraform.tfplan
```

## Terraform folder structure

│   .gitignore
│   .terraform.lock.hcl # lock file for terraform modules and providers
│   general.tf # general resources like resource_group
│   locals.tf # local variables file
│   network.tf # all network related resources
│   outputs.tf # output variables file
│   providers.tf # terraform providers initialiaztion
│   README.md # we are here :)
│   terraform.tf # terraform providers and modules list
│   variables.tf # external input variables file 
│   virtual_machines.tf # all VM related resources
│
└───environments # possible coffeemachine environments
        coffeemachine-{qa|prod}.tfvars 

# Variables section

## .tfvars variables

.tfvars file are storing some default variables for certain use cases/scenarios. Those variables can also be modified but do not put sensitive variables in those files:

| Name                | Description                                     |
| ------------------- | ----------------------------------------------- |
| `environment`       | `environment sufix/prefix`                      |
| `instance_type_VM`  | `VM type from Azure - remove Standard_ prefix`  |
| `location`          | `location of azure deployment`                  |

## Prompt variables

Variables in prompt at runtime are either sensitive (should be stored at env variables) or should be dynamic

| Name             | Description                                                                          | Default Value          |
| ---------------- | ------------------------------------------------------------------------------------ | ---------------------- |
| `project`        | `project sufix/prefix`                                                               | ""                     |
| `client`         | `client name for which is this being used`                                           | ""                     |
| `sufix`          | `OPTIONAL sufix for resources if you want to add`                                    | ""                     |
| `resource_group` | `name of resource_group - must be unique in subscription`                            | ""                     |
| `username`       | `administrator username VMs ->it's ssh pub key will be added from ~/.ssh/id_rsa.pub` | ""                     |