package setting

import "time"

type AppSettingS struct {
	Server struct {
		RunMode      string
		HttpPort     string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}

	Log struct {
		SavePath string
		FileName string
		FileExt  string
	}

	Auth struct {
		ServerKey  string
		TokenKey   string
		TokenDated int
	}

	Email struct {
		Sender   string
		Password string
		Host     string
		Port     string
	}

	Admin struct {
		Username string
		Password string
	}
}

type DbSettingS struct {
	Database struct {
		DBType    string
		UserName  string
		Password  string
		Host      string
		DBName    string
		Charset   string
		ParseTime bool

		MaxIdleConns int
		MaxOpenConns int
	}

	Redis struct {
		AddressPort string
		Password    string
		DefaultDB   int
		DialTimeout time.Duration
	}
}
