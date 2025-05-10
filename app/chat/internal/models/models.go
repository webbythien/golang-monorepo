package models

var registeredModels = make(map[string]interface{})

//nolint:unused
type gormTable interface {
	TableName() string
}

//nolint:unparam,unused
func registerAutoMigrate(mdls ...gormTable) error {
	for _, m := range mdls {
		registeredModels[m.TableName()] = m
	}
	return nil
}

//nolint:unused
func LoadMigrationModels() []interface{} {
	models := make([]interface{}, 0, len(registeredModels))
	for _, m := range registeredModels {
		models = append(models, m)
	}
	return models
}
