package commond

import (
	"gorm.io/gorm"
	"time"
)

//这个还是不错的，不错 抄了
type Create_Model struct {
	ID        uint           `gorm:"primarykey"` // 主键ID
	CreatedAt time.Time      // 创建时间
	UpdatedAt time.Time      // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}
