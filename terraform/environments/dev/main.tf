terraform {
  required_version = ">= 1.6.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }

  backend "s3" {
    bucket       = "terraform-test-bucket-1231"
    key          = "ecommerce/dev/terraform.tfstate"
    region       = "us-east-1"
    encrypt      = true
    use_lockfile = true
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = var.project_name
      Environment = var.environment
      ManagedBy   = "Terraform"
    }
  }
}

module "vpc" {
  source = "../../modules/vpc"

  project_name       = var.project_name
  environment        = var.environment
  vpc_cidr           = var.vpc_cidr
  availability_zones = var.availability_zones

  enable_nat_gateway = var.enable_nat_gateway
}

module "security_groups" {
  source = "../../modules/security-groups"

  project_name = var.project_name
  environment  = var.environment
  vpc_id       = module.vpc.vpc_id
}


module "ecr" {
  source = "../../modules/ecr"

  project_name = var.project_name
  environment  = var.environment

  repositories = [
    "product",
    "cart",
    "checkout",
    "order",
    "payment"
  ]
}

module "ecs" {
  source = "../../modules/ecs"

  project_name = var.project_name
  environment  = var.environment

  services = [
    "product",
    "cart",
    "checkout",
    "order",
    "payment"
  ]

  repository_urls = module.ecr.repository_urls

  service_cpu         = 256
  service_memory      = 512
  container_port      = 8080
  container_image_tag = "v1"

  log_retention_days = 14

  subnet_ids = module.vpc.private_subnet_ids

  security_group_id = module.security_groups.ecs_security_group_id

  target_group_arns = module.alb.target_group_arns

  desired_count = 1
}

module "vpc_endpoints" {
  source = "../../modules/vpc-endpoints"

  project_name = var.project_name
  environment  = var.environment
  vpc_id       = module.vpc.vpc_id
  aws_region   = var.aws_region

  private_subnet_ids = module.vpc.private_subnet_ids

  private_route_table_ids = module.vpc.private_route_table_ids

  ecs_security_group_id = module.security_groups.ecs_security_group_id
}

module "alb" {
  source = "../../modules/alb"

  project_name = var.project_name
  environment  = var.environment

  vpc_id = module.vpc.vpc_id

  public_subnet_ids = module.vpc.public_subnet_ids

  security_group_id = module.security_groups.alb_security_group_id

  services = [
    "product",
    "cart",
    "checkout",
    "order",
    "payment"
  ]

  container_port = 8080
}

module "s3_static_website" {
  source = "../../modules/s3"

}