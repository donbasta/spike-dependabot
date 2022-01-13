module "instance_01" {
  source                          = "git::ssh://git@source.golabs.io/gopay_infra/terraform/aws-basic-instance.git?ref=v5.0.14"
  ami_name                        = local.ami_name
  description                     = "Basic instance"
  app_name                        = local.labels.app_name
  country                         = local.labels.country
  environment                     = local.labels.environment
  focus_area                      = local.labels.focus_area
  is_temporary                    = local.labels.is_temporary
  instance_type                   = local.instance_type
  product_group                   = local.labels.product_group
  security_group_names            = local.security_group_names
  route53_internal_hosted_zone_id = var.route53_internal_hosted_zone_id
  vpc_id                          = var.vpc_id
  pod                             = local.labels.pod
  stream                          = local.labels.stream
  component                       = local.labels.component
}

module "instance_02" {
  source                          = "git::ssh://git@source.golabs.io/gopay_infra/terraform/gcloud-postgresql.git?ref=v1.2.4"
  ami_name                        = local.ami_name
  description                     = "PostgreSQL"
  app_name                        = local.labels.app_name
  country                         = local.labels.country
  environment                     = local.labels.environment
  focus_area                      = local.labels.focus_area
  is_temporary                    = local.labels.is_temporary
  instance_type                   = local.instance_type
  product_group                   = local.labels.product_group
  security_group_names            = local.security_group_names
  route53_internal_hosted_zone_id = var.route53_internal_hosted_zone_id
  vpc_id                          = var.vpc_id
  pod                             = local.labels.pod
  stream                          = local.labels.stream
  component                       = local.labels.component
}
