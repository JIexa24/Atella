package atella

// Logger implements logger abstraction.
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
}

// Config is configuration structure for atella.
type Config struct {
	Hostname string              `yaml:"hostname"`
	Logger   LoggerConfig        `yaml:"log"`
	Reporter ReporterConfig      `yaml:"reporter"`
	Channels []map[string]string `yaml:"channels"`
}

// ReporterConfig is configure log.
type ReporterConfig struct {
	MessagePath string `yaml:"message_path"`
	HexLen      int64  `yaml:"hex_len"`
}

// LoggerConfig is configure log.
type LoggerConfig struct {
	LogFile  string `yaml:"log_file"`
	LogLevel string `yaml:"log_level"`
}
