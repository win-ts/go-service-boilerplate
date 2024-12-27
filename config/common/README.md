# domain-name-common-cd-config-template

Common secrets repository for `domain-name`<br>
ArgoCD will sync secrets from this repository.

Available Environments
- alpha
- uat
- production (prod)

## How to use
### Secrets
To add or set secrets, add environment variables into `secrets.env` file of the desired environment.
