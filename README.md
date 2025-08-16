# Aurora Multi-Site Manager API 
version 0.3.0

## Overview
The Aurora Multi-Site Manager API is designed to facilitate the management of multiple sites within the Aurora ecosystem. It provides endpoints for creating, updating, and deleting sites, as well as managing site-specific configurations.

## Features
- Create, update, and delete sites
- Manage site configurations
- Retrieve site information
- Support for multiple site management
- Integration with Aurora's existing services

## Installation
- **Install atlas**: Follow the [Atlas installation guide](https://atlasgo.io/docs/installation) to install Atlas.
- **Install dependencies**: Run `go mod tidy` to install the required dependencies for the Aurora Multi-Site Manager API.
- **Run migrations**: Use Atlas to run the migrations for the Aurora Multi-Site Manager API. This will set up the necessary database schema.

```bash
atlas migrate apply --env "gorm" --url "mysql://<USERNAME>:<PASSWORD>@<DB_HOST>:<BD_PASSWORD>/<DB_NAME>?parseTime=true"
```

## Development
- **Clone the repository**: Clone the Aurora Multi-Site Manager API repository to your local machine.
  - The development branch is `main`.
  - The production branch is `release`.
- **Configure the API**:
  - **Set up environment variables**: Create a `.env` file in the root directory of the project and set the necessary environment variables for database connection, API keys, and other configurations.
  - **Set up config file**: Create a `<AURORA_APP_ENV>.json` file in the `./config` directory of the project. This file should contain the configuration settings for the Aurora Multi-Site Manager API, including database connection details, API keys, and other necessary configurations.
- **Start the API server**: Use the following command to start the Aurora Multi-Site Manager API server:
    ```bash
    go run ./cmd/aurora-api/main.go
    ```
- **Unit tests**: Run the unit tests to ensure that the API is functioning correctly. Use the following command:
    ```bash
    go test ./...
    ```

## Data structure
The data structure for the Aurora Multi-Site Manager API is designed to be flexible and extensible. It includes the following key components:
- **Site**: Represents a single site within the Aurora ecosystem. Each site has a unique identifier, name, and configuration settings.
- **Tenant**: Represents a tenant associated with a site. Each tenant has its own configuration and can be managed independently.
- **Template**: Represents a template that can be applied to a site or tenant. Templates define the structure and configuration of sites and tenants.
- **TemplateSetting**: Represents a setting within a template. Each setting can be customized for a specific site or tenant.
- **TemplateSettingOverride**: Represents an override for a template setting. Overrides allow for customization of settings at the site or tenant level.
- **Page**: Represents a page within a site tree. Pages can be organized hierarchically and can have their own configurations.
- **PageContent**: Represents the versioned content of a page. Each page can have multiple versions, allowing for content management and versioning. Only one version can be published at a time.
