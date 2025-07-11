# Create the ALB
resource "aws_security_group" "alb_sg" {
  name        = "alb-sg"
  description = "Allow HTTP from internet"
  vpc_id      = aws_vpc.main.id

  ingress {
    description = "HTTP from anywhere"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "alb-sg"
  }
}

resource "aws_lb_target_group" "http" {
  name     = "alb-target-group"
  port     = 80
  protocol = "HTTP"
  vpc_id   = aws_vpc.main.id

  health_check {
    path                = "/"
    protocol            = "HTTP"
    matcher             = "200"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 2
  }

  tags = {
    Name = "alb-target-group"
  }
}

resource "aws_lb" "http" {
  name               = "beard-alb"
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb_sg.id]
  subnets = [
    aws_subnet.public[0].id,
    aws_subnet.public[1].id
  ]
  enable_deletion_protection = false
  internal                   = false

  tags = {
    Name = "beard-alb"
  }
}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.http.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.http.arn
  }
}

resource "aws_lb_target_group_attachment" "ec2" {
  count            = var.instance_count
  target_group_arn = aws_lb_target_group.http.arn
  target_id        = module.ec2_beard.instance_id[count.index]
  port             = 80
}
