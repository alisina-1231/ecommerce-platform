aws_region   = "us-east-1"
project_name = "ecommerce-platform"
environment  = "dev"

vpc_cidr = "10.0.0.0/16"

availability_zones = [
  "us-east-1a",
  "us-east-1b"
]

# Enable only when testing private subnet internet access.
enable_nat_gateway = true