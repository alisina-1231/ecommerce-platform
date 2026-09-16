resource "aws_db_subnet_group" "this" {
  name       = "${var.name}-subnet-group"
  subnet_ids = var.private_subnet_ids

  tags = merge(
    var.tags,
    {
      Name = "${var.name}-subnet-group"
    }
  )
}

resource "aws_db_instance" "this" {
  identifier = var.name

  engine         = "postgres"
  engine_version = "16"

  instance_class        = var.instance_class
  allocated_storage     = var.allocated_storage
  max_allocated_storage = var.max_allocated_storage
  storage_type          = "gp3"
  storage_encrypted     = true

  db_name  = var.database_name
  username = var.username
  password = var.password
  port     = 5432

  db_subnet_group_name   = aws_db_subnet_group.this.name
  vpc_security_group_ids = var.security_group_ids

  publicly_accessible    = false
  multi_az               = false
  deletion_protection    = false
  skip_final_snapshot    = true

  backup_retention_period = var.backup_retention_period

  auto_minor_version_upgrade = true

  tags = merge(
    var.tags,
    {
      Name = var.name
    }
  )
}