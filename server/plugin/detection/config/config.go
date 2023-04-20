package config
type ModelConfig struct {
        ModelPath string `mapstructure:"model_path" json:"model_path" yaml:"model_path"`
        App string `mapstructure:"app" json:"app" yaml:"app"`
        Algorithm string `mapstructure:"algorithm" json:"algorithm" yaml:"algorithm"`
        Width int `mapstructure:"width" json:"width" yaml:"width"`
        Height int `mapstructure:"height" json:"height" yaml:"height"`
}
type Detection struct {
        ModelConfig []ModelConfig `mapstructure:"model_config" json:"model_config" yaml:"model_config"`
}
