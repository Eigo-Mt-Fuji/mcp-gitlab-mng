output "cluster_endpoint" {
  description = "Auroraクラスターのエンドポイント"
  value       = aws_rds_cluster.this.endpoint
}

output "cluster_reader_endpoint" {
  description = "Auroraクラスターのリーダーエンドポイント"
  value       = aws_rds_cluster.this.reader_endpoint
}

output "cluster_identifier" {
  description = "Auroraクラスターの識別子"
  value       = aws_rds_cluster.this.cluster_identifier
}

output "database_name" {
  description = "データベース名"
  value       = aws_rds_cluster.this.database_name
}

output "master_username" {
  description = "マスターユーザー名"
  value       = aws_rds_cluster.this.master_username
} 