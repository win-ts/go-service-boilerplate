# deployment-cd-config-template

Configuration repository for `service-name`<br>
ArgoCD will sync configmap, secrets, and resources config from this repository.

Available Environments
- alpha
- uat
- production (prod)

Code Repo URL: 

## How to use
### Image Tag
To set image tag, change the newTag value in `kustomization.yaml` of the desired environment.

### Configmaps
To set configmaps, add environment variables into `config.env` file of the desired environment.

### Secrets
To set secrets, add environment variables into `secrets.env` file of the desired environment.

### Resources
To set resources (replicas, limits/requests), change the values in `set_resources.yaml` file of the desired environment.
