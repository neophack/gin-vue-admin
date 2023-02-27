package config
type ModelConfig struct {
        ModelPath string `mapstructure:"model_path" json:"model_path" yaml:"model_path"`
        App string `mapstructure:"app" json:"app" yaml:"app"`
        Algorithm string `mapstructure:"algorithm" json:"algorithm" yaml:"algorithm"`
}
type Detection struct {
        ModelConfig []ModelConfig `mapstructure:"model_config" json:"model_config" yaml:"model_config"`
}
