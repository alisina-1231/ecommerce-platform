# E-commerce Platform

A cloud-native e-commerce platform built with AWS, Terraform, Go microservices, and modern DevOps practices.

## Project Goals

* Build a scalable e-commerce backend.
* Practice AWS cloud architecture.
* Implement infrastructure using Terraform.
* Build production-style microservices.
* Add CI/CD and observability.
* Explore AI-powered fraud detection and recommendations.

## Architecture

The platform will include:

* Product service
* Cart service
* Checkout service
* Order service
* Payment service
* AWS ECS
* AWS Lambda
* Amazon RDS PostgreSQL
* Amazon SQS
* Amazon Kinesis
* Amazon Cognito
* Amazon CloudFront
* AWS WAF
* Amazon SageMaker

## Project Structure

```text
ecommerce-platform/
├── terraform/
├── services/
├── stream-processing/
├── ml/
├── frontend/
├── monitoring/
├── docs/
├── .github/
├── .env.example
└── README.md
```

## Getting Started

### Prerequisites

* AWS account
* AWS CLI
* Terraform
* Docker
* Go
* Git

### Infrastructure

```bash
cd terraform

terraform init
terraform fmt -recursive
terraform validate
terraform plan
```

## Development Standards

See [Development Standards](docs/DEVELOPMENT_STANDARDS.md).

## License

