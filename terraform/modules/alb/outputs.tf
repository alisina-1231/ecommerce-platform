output "alb_id" {
  description = "Application Load Balancer ID"
  value       = aws_lb.this.id
}

output "alb_arn" {
  description = "Application Load Balancer ARN"
  value       = aws_lb.this.arn
}

output "alb_dns_name" {
  description = "Application Load Balancer DNS name"
  value       = aws_lb.this.dns_name
}

output "alb_zone_id" {
  description = "Application Load Balancer hosted zone ID"
  value       = aws_lb.this.zone_id
}

output "listener_arn" {
  description = "HTTP listener ARN"
  value       = aws_lb_listener.http.arn
}

output "target_group_arns" {
  description = "Target group ARNs by service"
  value = {
    for service, target_group in aws_lb_target_group.this :
    service => target_group.arn
  }
}

output "target_group_names" {
  description = "Target group names by service"
  value = {
    for service, target_group in aws_lb_target_group.this :
    service => target_group.name
  }
}