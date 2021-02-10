package env

type Env interface {
	GetEnvs(appMode string) map[string]string
}

type EnvVariables struct{}

func NewEnvService() Env {

	return &EnvVariables{}
}

func (e *EnvVariables) GetEnvs(appMode string) map[string]string {

	if appMode == "prod" {
		return prodMap
	}

	return devMap
}
