locals {
  container_images = {
    for service in var.services :
    service => "${var.repository_urls[service]}:${var.container_image_tag}"
  }
}

locals {
  service_environment = {
    for service in var.services :
    service => concat(
      [
        {
          name  = "PORT"
          value = tostring(var.container_port)
        },
        {
          name  = "APP_ENV"
          value = var.environment
        },
        {
          name  = "SERVICE_NAME"
          value = service
        }
      ],
      lookup(var.extra_environment, service, [])
    )
  }
}

resource "aws_ecs_task_definition" "this" {
  for_each = toset(var.services)

  family                   = "${var.project_name}-${var.environment}-${each.value}"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"

  cpu    = var.service_cpu
  memory = var.service_memory

  execution_role_arn = aws_iam_role.ecs_task_execution.arn
  task_role_arn      = aws_iam_role.ecs_task.arn

  container_definitions = jsonencode([
    {
      name      = each.value
      image     = local.container_images[each.value]
      essential = true

      portMappings = [
        {
          name          = each.value
          containerPort = var.container_port
          hostPort      = var.container_port
          protocol      = "tcp"
        }
      ]

      environment = local.service_environment[each.value]

      logConfiguration = {
        logDriver = "awslogs"

        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.ecs.name
          "awslogs-region"        = "us-east-1"
          "awslogs-stream-prefix" = "ecs"
        }
      }

      healthCheck = {
        command = [
          "CMD-SHELL",
          "wget --no-verbose --tries=1 --spider http://localhost:${var.container_port}/health || exit 1"
        ]

        interval    = 30
        timeout     = 5
        retries     = 3
        startPeriod = 30
      }
    }
  ])

  lifecycle {
    ignore_changes = [
      container_definitions
    ]
  }

  tags = {
    Name    = "${var.project_name}-${var.environment}-${each.value}"
    Service = each.value
  }
}