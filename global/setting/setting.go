package setting

import "github.com/spf13/viper"

const ConfigDir = "configs"
const ConfigName = "conf"
const ConfigSuffix = "yaml"

type Setting struct {
	vp *viper.Viper
}

func NewSetting() (*Setting, error) {
	vp := viper.New()
	vp.SetConfigName(ConfigName)
	vp.AddConfigPath(ConfigDir)
	vp.SetConfigType(ConfigSuffix)

	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Setting{vp}, nil
}

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	return nil
}
