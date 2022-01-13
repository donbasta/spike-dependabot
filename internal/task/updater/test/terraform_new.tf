module "instance_01" {
  source                          = "git::ssh://git@source.golabs.io/gopay_infra/terraform/aws-basic-instance.git?ref=v5.0.20"
  ami_name                        = local.ami_name
  description                     = "Basic instance"
  app_name                        = local.labels.app_name
  country                         = local.labels.country
  environment                     = local.labels.environment
  focus_area                      = local.labels.focus_area
}

module "instance_02" {
  ami_name                        = local.ami_name
  description                     = "PostgreSQL"
  app_name                        = local.labels.app_name
  country                         = local.labels.country
  environment                     = local.labels.environment
  focus_area                      = local.labels.focus_area
  source                          = "git::ssh://git@source.golabs.io/gopay_infra/terraform/gcloud-postgresql.git?ref=v1.2.4"
}

module "instance_03" {
  ami_name                        = local.ami_name
  description                     = "Redis"
  app_name                        = local.labels.app_name
  country                         = local.labels.country
  source                          = "git::ssh://git@source.golabs.io/gopay_infra/terraform/aws-redis.git?ref=v4.1.2"
  environment                     = local.labels.environment
  focus_area                      = local.labels.focus_area
}
