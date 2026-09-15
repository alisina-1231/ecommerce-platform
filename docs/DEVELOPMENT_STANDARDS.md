# E-commerce Platform — Development Standards

## 1. Project Overview

This project is a cloud-native e-commerce platform built using AWS, Terraform, Go, Python, and modern DevOps practices.

### Main Technologies

* **Infrastructure:** Terraform
* **Cloud:** AWS
* **Backend:** Go microservices
* **Frontend:** To be defined
* **Database:** PostgreSQL
* **Messaging:** Amazon SQS
* **Streaming:** Amazon Kinesis / Apache Flink
* **Authentication:** Amazon Cognito
* **Monitoring:** Prometheus, Grafana, OpenTelemetry
* **CI/CD:** GitHub Actions
* **Machine Learning:** Python, Amazon SageMaker

---

## 2. Git Standards

### Branch Naming

Use the following branch naming conventions:

* `main` — Production-ready code
* `develop` — Development integration
* `feature/<name>` — New features
* `fix/<name>` — Bug fixes
* `hotfix/<name>` — Critical production fixes
* `infra/<name>` — Infrastructure changes

### Commit Messages

Use clear and descriptive commit messages.

Format:

`type: description`

Examples:

* `feat: add product service`
* `fix: resolve order validation`
* `infra: create vpc module`
* `docs: update readme`
* `test: add payment tests`
* `refactor: improve checkout logic`

### Pull Requests

* Create a pull request for changes to `main`.
* Review code before merging.
* Ensure tests pass.
* Keep pull requests focused on one task.

---

## 3. Terraform Standards

### Directory Structure

Reusable infrastructure must be placed inside:

`terraform/modules/`

Environment-specific configuration must be placed inside:

`terraform/environments/`

### Naming

Use descriptive resource names.

Example:

```hcl
resource "aws_vpc" "this" {
  # Configuration
}
```

### Terraform Rules

* Run `terraform fmt` before committing.
* Run `terraform validate`.
* Run `terraform plan` before applying changes.
* Never commit `terraform.tfstate`.
* Never commit AWS credentials.
* Use variables instead of hardcoded values.
* Use outputs to expose important resource IDs.
* Avoid unnecessary changes to existing resources.

### State Management

Terraform state must be stored securely.

The production environment should use a remote backend with locking.

---

## 4. Go Standards

### Code Style

* Follow standard Go formatting.
* Run `gofmt`.
* Use meaningful variable and function names.
* Keep functions small and focused.
* Handle errors explicitly.
* Write unit tests for business logic.

### Project Layout

Each Go service should follow a consistent structure:

```text
service/
├── cmd/
├── internal/
├── migrations/
├── Dockerfile
├── go.mod
└── README.md
```

### API Standards

* Use RESTful HTTP APIs where appropriate.
* Return consistent JSON responses.
* Use proper HTTP status codes.
* Validate incoming requests.
* Add health and readiness endpoints.

---

## 5. Security Standards

* Never commit secrets or credentials.
* Use AWS IAM roles instead of long-lived access keys.
* Follow least-privilege access.
* Store secrets in AWS Secrets Manager or Parameter Store.
* Validate and sanitize user input.
* Use HTTPS for public APIs.
* Do not expose database ports publicly.
* Do not expose private service endpoints unnecessarily.

---

## 6. Docker Standards

* Use official base images.
* Prefer multi-stage builds.
* Keep images small.
* Do not run applications as root unless necessary.
* Do not include secrets in images.
* Add health checks where appropriate.
* Use versioned image tags.

---

## 7. Testing Standards

Each service should include:

* Unit tests
* API tests where applicable
* Integration tests where applicable

Before merging:

```bash
go test ./...
```

For Terraform:

```bash
terraform fmt -check -recursive
terraform validate
```

---

## 8. Observability Standards

Each service should provide:

* Structured logging
* Health endpoint
* Readiness endpoint
* Metrics where appropriate
* Distributed tracing where appropriate

Use OpenTelemetry for distributed tracing.

---

## 9. Environment Standards

Supported environments:

* `dev`
* `prod`

Development resources must be isolated from production resources.

Use environment-specific variables and Terraform configurations.

---

## 10. Documentation Standards

Every service should include a README covering:

* Service purpose
* Local setup
* Environment variables
* API endpoints
* How to run tests
* How to run the service

Keep documentation updated with major changes.
