package config

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Backup   BackupConfig   `mapstructure:"backup"`
	Log      LogConfig      `mapstructure:"log"`
	Notify   NotifyConfig   `mapstructure:"notify"`
}

type DatabaseConfig struct {
	Type     string `mapstructure:"type"` // mysql, postgres, mongodb
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	ExtraParams string `mapstructure:"extra_params"` // e.g. sslmode=disable
	ToolPath string `mapstructure:"tool_path"` // path to mysqldump, pg_dump, etc.
}

type StorageConfig struct {
	Type            string `mapstructure:"type"` // local, s3, gcs, azure
	Path            string `mapstructure:"path"` // local path or bucket name
	Region          string `mapstructure:"region"` // for cloud
	CredentialsFile string `mapstructure:"credentials_file"` // for cloud
}

type BackupConfig struct {
	Type        string `mapstructure:"type"` // full, incremental, differential
	Compression bool   `mapstructure:"compression"`
	Schedule    string `mapstructure:"schedule"` // cron expression
}

type LogConfig struct {
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
}

type NotifyConfig struct {
	SlackWebhookURL string `mapstructure:"slack_webhook_url"`
}
