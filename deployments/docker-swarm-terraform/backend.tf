terraform {
  backend "s3" {
    bucket                      = "terraform"
    key                         = "monitor.tfstate"
    endpoint                    = "s3.storage.selcloud.ru"
    region                      = "ru-1"
    skip_credentials_validation = true
    skip_region_validation      = true
    force_path_style            = true
  }
}
