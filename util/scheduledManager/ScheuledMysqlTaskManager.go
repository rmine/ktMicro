package ktMicro

import (
	"github.com/rmine/ktMicro/util/ipUtil"
	ktimer "github.com/rmine/ktMicro/util/ktimer"
	"github.com/satori/go.uuid"
	"math/rand"
	"sort"
	"strconv"
	"sync"
)

type ScheduledMysqlManager struct {
	identifier     string
	taskInfo       map[string]*ktimer.KTimer
	heartbeatTimer *ktimer.KTimer
}

var syncMysqlOnce sync.Once
var shareMysqlManager *ScheduledMysqlManager
var scheuledMysqlTaskDao *ScheduledMysqlTaskModel

func DefaultScheduledMysqlManager() *ScheduledMysqlManager {
	syncMysqlOnce.Do(func() {
		shareMysqlManager = &ScheduledMysqlManager{
			taskInfo: map[string]*ktimer.KTimer{},
		}
		scheuledMysqlTaskDao = NewScheuledTaskDao()
		shareMysqlManager.genUniqueIdentifier()
		shareMysqlManager.initHeartbeat()
	})
	return shareMysqlManager
}

func (m *ScheduledMysqlManager) genUniqueIdentifier() {
	if len(m.identifier) == 0 {
		if ip, err1 := ipUtil.ExternalIP(); err1 == nil {
			m.identifier = ip
		} else if uuid := uuid.NewV4(); uuid.String() != "" {
			m.identifier = uuid.String()
		} else {
			m.identifier = strconv.FormatInt(rand.Int63(), 10)
		}
	}
}

func (m *ScheduledMysqlManager) initHeartbeat() {
	if m.heartbeatTimer == nil {
		m.heartbeatTimer = ktimer.NewTimer(60, 0, m.heartbeat)
		m.heartbeatTimer.Start()
	}
}

func (m *ScheduledMysqlManager) heartbeat(timer *ktimer.KTimer) {
	unit, _ := scheuledMysqlTaskDao.QueryTaskWithIdentifier(m.identifier)
	if unit == nil || unit.Id == 0 {
		(&ScheduledMysqlTaskModel{UUid: m.identifier}).CreateNewTask()
	} else {
		unit.UpdateTaskStatus(m.identifier)
	}
}
/*
taskName 任务名称
spec cron
interval 时间间隔 单位秒
singleUnit 是否单机运行
block block
*/
func (m *ScheduledMysqlManager) addTask(taskName string, spec string, interval float64, singleUnit bool, block ktimer.KTimerBlock) bool {
	if len(taskName) > 0 && m.taskInfo[taskName] != nil {
		return false
	}
	if singleUnit {
		block = m.wrapperBlock(block)
	}
	var timer *ktimer.KTimer = nil
	if len(spec) > 0 {
		timer = ktimer.NewCronTimer(spec, block)
	} else {
		timer = ktimer.NewTimer(interval, 0, block)
	}
	timer.Start()
	if len(taskName) > 0 {
		m.taskInfo[taskName] = timer
	}
	return true
}

func (m *ScheduledMysqlManager) wrapperBlock(block ktimer.KTimerBlock) ktimer.KTimerBlock {
	return func(timer *ktimer.KTimer) {
		if m.ValidateExecuteAuthority() {
			block(timer)
		}
	}
}

func (m *ScheduledMysqlManager) ValidateExecuteAuthority() bool {
	unitList, err := scheuledMysqlTaskDao.QueryAliveTaskList(300)
	if err != nil || len(unitList) == 0 {
		return false
	}
	identifierList := make([]string, len(unitList))
	for idx, unit := range unitList {
		identifierList[idx] = unit.UUid
	}
	sort.Strings(identifierList)
	return identifierList[0] == m.identifier
}

func (m *ScheduledMysqlManager) Start() {
	m.initHeartbeat()
	m.heartbeatTimer.Fire()
}

/*
taskName 任务名称
spec cron
singleUnit 是否单机运行
block block
*/
func (m *ScheduledMysqlManager) AddCronTask(taskName string, spec string, singleUnit bool, block ktimer.KTimerBlock) bool {
	return m.addTask(taskName, spec, 0, singleUnit, block)
}

/*
taskName 任务名称
interval 时间间隔 单位秒
singleUnit 是否单机运行
block block
*/
func (m *ScheduledMysqlManager) AddTask(taskName string, interval float64, singleUnit bool, block ktimer.KTimerBlock) bool {
	return m.addTask(taskName, "", interval, singleUnit, block)
}

func (m *ScheduledMysqlManager) RemoveTask(taskName string) {
	if len(taskName) > 0 {
		timer := m.taskInfo[taskName]
		if timer != nil {
			timer.Stop(true)
			m.taskInfo[taskName] = nil
		}
	}
}

func (m *ScheduledMysqlManager) InvokeTask(taskName string) {
	if len(taskName) > 0 {
		timer := m.taskInfo[taskName]
		if timer != nil {
			timer.Fire()
		}
	}
}