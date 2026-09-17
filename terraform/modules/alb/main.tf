resource "aws_lb" "this" {
  name               = "${var.project_name}-${var.environment}"
  internal           = false
  load_balancer_type = "application"

  security_groups = [
    var.security_group_id
  ]

  subnets = var.public_subnet_ids

  enable_deletion_protection = false

  tags = {
    Name = "${var.project_name}-${var.environment}-alb"
  }
}


# Target Groups

resource "aws_lb_target_group" "this" {
  for_each = toset(var.services)

  name = "${var.project_name}-${var.environment}-${each.value}"

  port        = var.container_port
  protocol    = "HTTP"
  target_type = "ip"

  vpc_id = var.vpc_id

  health_check {
    enabled             = true
    path                = "/health"
    protocol            = "HTTP"
    port                = "traffic-port"
    matcher             = "200"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 3
  }

  deregistration_delay = 30

  tags = {
    Name    = "${var.project_name}-${var.environment}-${each.value}"
    Service = each.value
  }
}


# HTTP Listener

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.this.arn

  port     = 80
  protocol = "HTTP"

  default_action {
    type = "fixed-response"

    fixed_response {
      content_type = "application/json"
      message_body = jsonencode({
        message = "Ecommerce Platform API"
      })
      status_code = "200"
    }
  }

  tags = {
    Name = "${var.project_name}-${var.environment}-http"
  }
}


# Product
# /product/*

resource "aws_lb_listener_rule" "product" {
  listener_arn = aws_lb_listener.http.arn

  priority = 100

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.this["product"].arn
  }

  condition {
    path_pattern {
      values = [
        "/product/*"
      ]
    }
  }
}


# Cart
# /cart/*

resource "aws_lb_listener_rule" "cart" {
  listener_arn = aws_lb_listener.http.arn

  priority = 110

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.this["cart"].arn
  }

  condition {
    path_pattern {
      values = [
        "/cart/*"
      ]
    }
  }
}


# Checkout
# /checkout/*

resource "aws_lb_listener_rule" "checkout" {
  listener_arn = aws_lb_listener.http.arn

  priority = 120

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.this["checkout"].arn
  }

  condition {
    path_pattern {
      values = [
        "/checkout/*"
      ]
    }
  }
}


# Orders
# /orders/*

resource "aws_lb_listener_rule" "order" {
  listener_arn = aws_lb_listener.http.arn

  priority = 130

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.this["order"].arn
  }

  condition {
    path_pattern {
      values = [
        "/order/*"
      ]
    }
  }
}


# Payments
# /payments/*

resource "aws_lb_listener_rule" "payment" {
  listener_arn = aws_lb_listener.http.arn

  priority = 140

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.this["payment"].arn
  }

  condition {
    path_pattern {
      values = [
        "/payment/*"
      ]
    }
  }
}