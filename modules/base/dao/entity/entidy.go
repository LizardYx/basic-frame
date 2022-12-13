package entity

// GetTables 自动映射数据表
func GetTables() []interface{} {
	return []interface{}{
		new(Menu),
		new(Button),
		new(RestfulApi),
		new(DisabledField),
		new(SystemConfig),
	}
}
