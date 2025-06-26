# Aurora MySQL Serverless v2 RDS モジュール

このTerraformモジュールは、既存のVPC・サブネット・セキュリティグループを利用してAurora MySQL Serverless v2クラスターを作成します。

## 前提条件
- VPC、サブネット、セキュリティグループは既存リソースを利用します
- Terraform >= 1.0

## 使い方

```hcl
module "aurora" {
  source = "./module/rds"

  db_subnet_group_name    = "example-aurora-subnet-group"
  subnet_ids              = ["subnet-xxxxxxxx", "subnet-yyyyyyyy"]
  vpc_security_group_ids  = ["sg-xxxxxxxx"]
  cluster_identifier      = "example-aurora-cluster"
  engine_version          = "8.0.mysql_aurora.3.05.1"
  database_name           = "myapp"
  master_username         = "admin"
  master_password         = "your-secure-password-here"
  min_capacity            = 0.5
  max_capacity            = 16
  instance_count          = 2
  skip_final_snapshot     = true
  deletion_protection     = false
  tags = {
    Environment = "dev"
    Project     = "example"
  }
}
```

## 主要変数
- `subnet_ids` : プライベートサブネットIDリスト
- `vpc_security_group_ids` : AuroraにアタッチするセキュリティグループIDリスト
- `master_password` : 強力なパスワードを指定してください

## 出力値
- `cluster_endpoint` : クラスターエンドポイント
- `cluster_reader_endpoint` : リーダーエンドポイント
- `cluster_identifier` : クラスター識別子
- `database_name` : データベース名
- `master_username` : マスターユーザー名

## 注意事項
- `master_password`は機密情報として管理してください
- 本番環境では`deletion_protection = true`を推奨します

## サンプル
`terraform.tfvars.example` を参考にしてください。 