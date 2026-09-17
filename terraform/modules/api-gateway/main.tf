resource "aws_apigatewayv2_api" "this" {
  name          = var.name
  protocol_type = "HTTP"

  cors_configuration {
    allow_headers = [
      "Authorization",
      "Content-Type"
    ]

    allow_methods = [
      "GET",
      "POST",
      "PUT",
      "DELETE",
      "OPTIONS"
    ]

    allow_origins = ["*"]

    max_age = 300
  }

  tags = var.tags
}


# ---------------------------------------------------------
# VPC LINK
# ---------------------------------------------------------

resource "aws_apigatewayv2_vpc_link" "this" {
  name               = "${var.name}-vpc-link"
  security_group_ids = var.security_group_ids
  subnet_ids         = var.subnet_ids

  tags = var.tags
}


# ---------------------------------------------------------
# COGNITO JWT AUTHORIZER
# ---------------------------------------------------------

resource "aws_apigatewayv2_authorizer" "cognito" {
  api_id = aws_apigatewayv2_api.this.id

  authorizer_type  = "JWT"
  identity_sources = ["$request.header.Authorization"]
  name             = "${var.name}-cognito"

  jwt_configuration {
    audience = var.jwt_audience
    issuer   = var.jwt_issuer
  }
}


# ---------------------------------------------------------
# ALB INTEGRATION
# ---------------------------------------------------------

resource "aws_apigatewayv2_integration" "alb" {
  api_id = aws_apigatewayv2_api.this.id

  integration_type = "HTTP_PROXY"

  integration_uri = var.alb_listener_arn

  integration_method = "ANY"

  connection_type = "VPC_LINK"
  connection_id   = aws_apigatewayv2_vpc_link.this.id

  payload_format_version = "1.0"

  timeout_milliseconds = 30000
}


# ---------------------------------------------------------
# PRODUCT
# ---------------------------------------------------------

resource "aws_apigatewayv2_route" "product" {
  api_id = aws_apigatewayv2_api.this.id

  route_key = "ANY /product/{proxy+}"

  target = "integrations/${aws_apigatewayv2_integration.alb.id}"

  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
  authorization_type = "JWT"
}


# ---------------------------------------------------------
# CART
# ---------------------------------------------------------

resource "aws_apigatewayv2_route" "cart" {
  api_id = aws_apigatewayv2_api.this.id

  route_key = "ANY /cart/{proxy+}"

  target = "integrations/${aws_apigatewayv2_integration.alb.id}"

  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
  authorization_type = "JWT"
}


# ---------------------------------------------------------
# CHECKOUT
# ---------------------------------------------------------

resource "aws_apigatewayv2_route" "checkout" {
  api_id = aws_apigatewayv2_api.this.id

  route_key = "ANY /checkout/{proxy+}"

  target = "integrations/${aws_apigatewayv2_integration.alb.id}"

  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
  authorization_type = "JWT"
}


# ---------------------------------------------------------
# ORDERS
# ---------------------------------------------------------

resource "aws_apigatewayv2_route" "orders" {
  api_id = aws_apigatewayv2_api.this.id

  route_key = "ANY /order/{proxy+}"

  target = "integrations/${aws_apigatewayv2_integration.alb.id}"

  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
  authorization_type = "JWT"
}


# ---------------------------------------------------------
# PAYMENTS
# ---------------------------------------------------------

resource "aws_apigatewayv2_route" "payments" {
  api_id = aws_apigatewayv2_api.this.id

  route_key = "ANY /payment/{proxy+}"

  target = "integrations/${aws_apigatewayv2_integration.alb.id}"

  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
  authorization_type = "JWT"
}


# ---------------------------------------------------------
# STAGE + THROTTLING
# ---------------------------------------------------------

resource "aws_apigatewayv2_stage" "default" {
  api_id = aws_apigatewayv2_api.this.id

  name = "$default"

  auto_deploy = true

  default_route_settings {
    throttling_burst_limit = var.throttling_burst_limit
    throttling_rate_limit  = var.throttling_rate_limit
  }

  tags = var.tags
}