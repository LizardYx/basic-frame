package entity

// GetTables 自动映射数据表
func GetTables() []interface{} {
	return []interface{}{
		new(User),
		new(UserExtendInfo),
		new(Organization),
		new(Position),
		new(Role),
		new(Menu),
		new(Button),
		new(RestfulApi),
		new(DisabledField),
		new(SecurityLevel),
		new(TagManage),
		new(SystemConfig),
	}
}
