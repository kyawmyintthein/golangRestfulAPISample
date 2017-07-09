package bootstrap

import "strings"

var (
	APP_NAME = `App Name`

	APP_VERSION = "1.0.0"

	DEV_ENV = "dev"

	STAGING_ENV = "staging"

	PROD_ENV = "prod"

	REPLACER *strings.Replacer = strings.NewReplacer(".", "_")

	APP_CONFIG_PREFIX = `APP`

	APP_CONFIG_NAME = `application`

	DB_CONFIG_PREFIX = `DB`

	DB_CONFIG_NAME = `database`

	CONFIG_PATH = `config`

	CONFIG_FILE_TYPE = `yaml`
)
