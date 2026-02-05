package service

import "anyrouter-checkin/internal/repository"

func GetConfigs(category string) (map[string]string, error) {
	configs, err := repository.ListConfigs(category)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, cfg := range configs {
		result[cfg.Key] = cfg.Value
	}
	return result, nil
}

func GetConfig(key string) string {
	value, err := repository.GetConfigValue(key)
	if err != nil {
		return ""
	}
	return value
}

func SetConfig(key, value, category string) error {
	return repository.SetConfigValue(key, value, category)
}
