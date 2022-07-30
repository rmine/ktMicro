package ktMicro

import (
	"github.com/jinzhu/gorm"
	ktMicro "github.com/rmine/ktMicro/config"
	"time"
)

type ScheduledMysqlTaskModel struct {
	ktMicro.Mysql
	Id        uint32    `json:"id" gorm:"column:id;type:int(10) unsigned auto_increment; NOT NULL; COMMENT:'id';primary_key;"`
	UUid      string    `json:"uuid" gorm:"column:uuid;type:varchar(255); NOT NULL; COMMENT:'唯一标识';"`
	TaskIP    string    `json:"task_ip" gorm:"column:task_ip;type:varchar(255); NOT NULL; COMMENT:'唯一标识';"`
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime; NOT NULL; DEFAULT:CURRENT_TIMESTAMP; COMMENT:'创建时间';"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:datetime; NOT NULL; DEFAULT: CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP; COMMENT:'更新时间';"`
}

func (ScheduledMysqlTaskModel) TableName() string {
	return "scheuled_task"
}

func NewScheuledTaskDao() *ScheduledMysqlTaskModel {
	m := new(ScheduledMysqlTaskModel)
	return m
}

func (m *ScheduledMysqlTaskModel) CreateNewTask() error {
	_, err := m.Write(func(db *gorm.DB) *gorm.DB {
		//插入使用数据库时间
		return db.Exec("INSERT INTO `scheuled_task` (`created_at`,`updated_at`,`deleted_at`,`uuid`,`task_ip`) VALUES (NOW(),NOW(),NULL,?,'')", m.UUid)
		//return db.Create(m)
	})
	return err
}

func (m *ScheduledMysqlTaskModel) QueryTask(query interface{}, args ...interface{}) (result ScheduledMysqlTaskModel, err error) {
	_, err = m.Read(func(db *gorm.DB) *gorm.DB {
		return db.Where(query, args...).First(&result)
	})
	return result, err
}

func (m *ScheduledMysqlTaskModel) QueryTaskWithIdentifier(identifier string) (*ScheduledMysqlTaskModel, error) {
	model := &ScheduledMysqlTaskModel{}
	_, err := m.Read(func(db *gorm.DB) *gorm.DB {
		return db.Where("uuid=?", identifier).First(model)
	})
	if err != nil || model.Id == 0 {
		model = nil
	}
	return model, err
}

func (m *ScheduledMysqlTaskModel) QueryAliveTaskList(seconds int) ([]ScheduledMysqlTaskModel, error) {
	//以数据库时间查询活跃机器
	return m.QueryTaskList("TIMESTAMPDIFF(SECOND,updated_at,NOW()) <= ?", seconds)
	//return m.QueryTaskList("updated_at >= ?", time.Now().Add(-time.Second * time.Duration(seconds)))
}

func (m *ScheduledMysqlTaskModel) QueryTaskList(query interface{}, args ...interface{}) (result []ScheduledMysqlTaskModel, err error) {
	_, err = m.Read(func(db *gorm.DB) *gorm.DB {
		return db.Where(query, args...).Find(&result)
	})
	return result, err
}

func (m *ScheduledMysqlTaskModel) UpdateTaskStatus(uuid string) (err error) {
	_, err = m.Read(func(db *gorm.DB) *gorm.DB {
		//更新为数据库时间
		return db.Exec("UPDATE "+m.TableName()+" SET updated_at = NOW() WHERE uuid = ?", uuid)
		//return db.Model(m).Where("uuid=?",uuid).Update("updated_at", time.Now())
	})
	return nil
}
