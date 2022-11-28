package entity

// GetTables 自动映射数据表
func GetTables() []interface{} {
	return []interface{}{
		new(RestfulApi),
		new(SystemConfig),
	}
}
