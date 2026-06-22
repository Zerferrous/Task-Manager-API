package logger

type Config struct {
	Level  string `envconfig:"LOG_LEVEL" default:"info"`
	Folder string `envconfig:"LOG_FOLDER" default:"out/logs"`
}
