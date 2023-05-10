package config

type ModelConfig struct {
	ProgramPath string `mapstructure:"program_path" json:"program_path" yaml:"program_path"`
	ProgramName string `mapstructure:"program_name" json:"program_name" yaml:"program_name"`
	App         string `mapstructure:"app" json:"app" yaml:"app"`
	DataType    string `mapstructure:"data_type" json:"data_type" yaml:"data_type"`
}
type Detection struct {
	ModelConfig []ModelConfig `mapstructure:"model_config" json:"model_config" yaml:"model_config"`
}
