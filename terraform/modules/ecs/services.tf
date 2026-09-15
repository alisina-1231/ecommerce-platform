resource "aws_ecs_service" "this" {
  for_each = toset(var.services)

  name = "${var.project_name}-${var.environment}-${each.value}"

  cluster = aws_ecs_cluster.this.id

  task_definition = aws_ecs_task_definition.this[each.value].arn

  desired_count = var.desired_count

  launch_type = "FARGATE"

  platform_version = "LATEST"

  deployment_minimum_healthy_percent = 0
  deployment_maximum_percent         = 100

  health_check_grace_period_seconds = 60

  network_configuration {
    subnets = var.subnet_ids

    security_groups = [
      var.security_group_id
    ]

    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = var.target_group_arns[each.value]

    container_name = each.value

    container_port = var.container_port
  }

  depends_on = [
    aws_iam_role_policy_attachment.ecs_task_execution
  ]

  tags = {
    Name    = "${var.project_name}-${var.environment}-${each.value}"
    Service = each.value
  }
}