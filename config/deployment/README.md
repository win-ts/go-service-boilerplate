# deployment-cd-config-template

Configuration repository for `<service-name>`<br>
ArgoCD will sync configmap, secrets, and resources config from this repository.

Available Environments
- alpha
- production (prod)

## How to use
### Configmaps
To set configmaps, add environment variables into `config.env` file of the desired environment.

### Secrets
To set secrets, add environment variables into `secrets.env` file of the desired environment.

### Resources
To set resources (replicas, limits/requests), change the values in `set_resources.yaml` file of the desired environment.
