package schema

import "time"

type UserExtendInfo struct {
	ID          uint64    `json:"id"`           // 唯一标识
	UserID      uint64    `json:"user_id"`      // 用户ID
	RealName    string    `json:"real_name"`    // 用户昵称
	MobilePhone int       `json:"mobile_phone"` // 移动手机
	QQAccount   int       `json:"qq_account"`   // QQ账号
	Email       string    `json:"email"`        // 邮箱账号
	Creator     uint64    `json:"creator"`      // 创建者
	CreatedAt   time.Time `json:"created_at"`   // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`   // 更新时间
}
