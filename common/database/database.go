package database

// ConnectionInfo store information to connect to a database
type ConnectionInfo struct {
	DBType          string `yaml:"db_type" json:"db_type"`
	DBUsername      string `yaml:"db_username" json:"db_username"`
	DBPassword      string `yaml:"db_password" json:"db_password"`
	DBName          string `yaml:"db_name" json:"db_name"`
	DBSSLModeOption string `yaml:"db_ssl_mode_option" json:"db_ssl_mode_option"`
	DBHostname      string `yaml:"db_hostname" json:"db_hostname"`
	DBPort          string `yaml:"db_port" json:"db_port"`
	DBEnvironment   string `yaml:"db_environment" json:"db_environment"`
	DBSchemaName    string `yaml:"db_schema_name" json:"db_schema_name"`
}

// ExecutiveAccount store information of user manage tables define in model_list
type ExecutiveAccount struct {
	UserName      string `yaml:"username" json:"username"`
	Password      string `yaml:"password" json:"password"`
	ForceRecreate bool   `yaml:"force_recreate" json:"force_recreate"`
}

// Model store information of model can manage
type Model struct {
	TableName         string   `yaml:"table_name" json:"table_name"`
	Columns           []Column `yaml:"columns" json:"columns"`
	AutoMigration     bool     `yaml:"auto_migration" json:"auto_migration"` // auto_migration if table not exist or misisng column
	DisplayName       string   `yaml:"display_name" json:"display_name"`
	NameDisplayColumn string   `yaml:"name_display_column" json:"name_display_column"`
}

// Models array of model
type Models []Model

// ColumnsByTableName create map columns by table name from array of column
func (ms Models) ColumnsByTableName() map[string][]Column {
	res := make(map[string][]Column)
	for _, m := range ms {
		if _, ok := res[m.TableName]; ok {
			res[m.TableName] = append(res[m.TableName], m.Columns...)
		} else {
			res[m.TableName] = m.Columns
		}
	}

	return res
}

// Column store information of a column
type Column struct {
	Name         string `yaml:"name" json:"name"`
	Type         string `yaml:"type" json:"type"`
	Tags         string `yaml:"tags" json:"tags"`
	IsNullable   bool   `yaml:"is_nullable" json:"is_nullable"`
	IsPrimary    bool   `yaml:"is_primary" json:"is_primary"`
	DefaultValue string `yaml:"default_value" json:"default_value"`
}

// Columns array of column
type Columns []Column

// GroupByName group column by name
func (cols Columns) GroupByName() map[string][]Column {
	res := make(map[string][]Column)
	for _, col := range cols {
		if _, ok := res[col.Name]; ok {
			res[col.Name] = append(res[col.Name], col)
		} else {
			res[col.Name] = []Column{col}
		}
	}

	return res
}

// Names return names of all columns
func (cols Columns) Names() []string {
	res := []string{}
	for _, col := range cols {
		res = append(res, col.Name)

	}

	return res
}