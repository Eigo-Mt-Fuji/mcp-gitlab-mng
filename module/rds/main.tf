resource "aws_rds_cluster" "this" {
  cluster_identifier     = var.cluster_identifier
  engine                = "aurora-mysql"
  engine_mode           = "provisioned"
  engine_version        = var.engine_version
  database_name         = var.database_name
  master_username       = var.master_username
  master_password       = var.master_password
  skip_final_snapshot   = var.skip_final_snapshot
  deletion_protection   = var.deletion_protection

  db_subnet_group_name   = aws_db_subnet_group.this.name
  vpc_security_group_ids = var.vpc_security_group_ids

  serverlessv2_scaling_configuration {
    min_capacity = var.min_capacity
    max_capacity = var.max_capacity
  }

  tags = var.tags
}

resource "aws_db_subnet_group" "this" {
  name       = var.db_subnet_group_name
  subnet_ids = var.subnet_ids

  tags = var.tags
}

resource "aws_rds_cluster_instance" "this" {
  count                = var.instance_count
  identifier           = "${var.cluster_identifier}-instance-${count.index + 1}"
  cluster_identifier   = aws_rds_cluster.this.id
  instance_class       = "db.serverless"
  engine               = aws_rds_cluster.this.engine
  engine_version       = aws_rds_cluster.this.engine_version

  tags = var.tags
} 