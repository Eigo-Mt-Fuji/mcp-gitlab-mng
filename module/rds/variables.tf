variable "db_subnet_group_name" {
  description = "DBサブネットグループ名"
  type        = string
}

variable "subnet_ids" {
  description = "Aurora用のサブネットIDリスト"
  type        = list(string)
}

variable "vpc_security_group_ids" {
  description = "AuroraにアタッチするセキュリティグループIDリスト"
  type        = list(string)
}

variable "cluster_identifier" {
  description = "Auroraクラスターの識別子"
  type        = string
}

variable "engine_version" {
  description = "Aurora MySQLエンジンバージョン"
  type        = string
  default     = "8.0.mysql_aurora.3.05.1"
}

variable "database_name" {
  description = "データベース名"
  type        = string
}

variable "master_username" {
  description = "マスターユーザー名"
  type        = string
}

variable "master_password" {
  description = "マスターパスワード"
  type        = string
  sensitive   = true
}

variable "min_capacity" {
  description = "Serverless v2の最小ACU"
  type        = number
  default     = 0.5
}

variable "max_capacity" {
  description = "Serverless v2の最大ACU"
  type        = number
  default     = 16
}

variable "instance_count" {
  description = "Auroraインスタンス数"
  type        = number
  default     = 2
}

variable "skip_final_snapshot" {
  description = "削除時にスナップショットをスキップするか"
  type        = bool
  default     = true
}

variable "deletion_protection" {
  description = "削除保護の有効/無効"
  type        = bool
  default     = false
}

variable "tags" {
  description = "リソース共通タグ"
  type        = map(string)
  default     = {}
} 