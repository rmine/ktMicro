package ktMicro

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	ktimer "github.com/rmine/ktMicro/util/ktimer"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math/rand"
	"sync"
	"time"
)

//======== mysqlConf Model start ========
type MysqlConf struct {
	Cluster []MysqlCluster `yaml:"MysqlClusters"`
}

type MysqlCluster struct {
	Db       string        `yaml:"Db"`
	DBServer MysqlDBServer `yaml:"DBServer"`
}

type MysqlDBServer struct {
	DBName string            `yaml:"DBName"`
	Write  MysqlConnection   `yaml:"Write"`
	Reads  []MysqlConnection `yaml:"Reads"`
	Pool   MysqlPool         `yaml:"Pool"`
}

type MysqlConnection struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
}

type MysqlPool struct {
	MaxActive int           `yaml:"MaxActive"` //最大连接数据库连接数,设 0 为没有限制
	MaxIdle   int           `yaml:"MaxIdle"`   //最大等待连接中的数量,设 0 为没有限制
	MaxWait   time.Duration `yaml:"MaxWait"`   //最大等待毫秒数, 单位为 ms, 超过时间会出错误信息
}

//======== mysqlConf Model end ========

type MysqlManager struct {
	err   error
	timer *ktimer.KTimer
}

type Mysql struct {
	Database string `gorm:"-"`
	Table    string `gorm:"-"`
	Logger   *logrus.Logger
}

var (
	mysqlOnce         sync.Once
	shareMysqlManager *MysqlManager
	mysqlconf         MysqlConf
	writeDBMap        map[string]*gorm.DB
	readsDBMap        map[string][]*gorm.DB
	dbLockMap         map[string]*sync.RWMutex
)

func init() {
	writeDBMap = make(map[string]*gorm.DB)
	readsDBMap = make(map[string][]*gorm.DB)
	dbLockMap = make(map[string]*sync.RWMutex)
}

func SharedMysqlManager() *MysqlManager {
	mysqlOnce.Do(func() {
		shareMysqlManager = &MysqlManager{}
	})
	return shareMysqlManager
}

func (m *MysqlManager) startMonitor() {
	if m.timer == nil {
		m.timer = ktimer.NewTimer(60, 0, m.verifyDBConfig)
		m.timer.Start()
		m.timer.Fire()
	}
}

//数据库超时重连
func (m *MysqlManager) verifyDBConfig(t *ktimer.KTimer) {
	path := fmt.Sprintf("configFiles/%v/mysql.yml", AppEnv())
	mysqlconf = getMysqlConf(path)
	for _, cluster := range mysqlconf.Cluster {
		//write
		writeDB := writeDBMap[cluster.Db]
		if writeDB == nil || writeDB.DB().Ping() != nil {
			if writeDB != nil {
				writeDB.Close()
				log.Fatalln("Write数据库异常断开")
			}
			if db, err := initDataPool(cluster.DBServer.Write, cluster.DBServer.Pool, cluster.Db); err == nil {
				writeDBMap[cluster.Db] = db
			} else {
				m.err = err
				writeDBMap[cluster.Db] = nil
				log.Fatalln("Write数据库启动失败", err)
			}
		}
		//read
		if len(cluster.DBServer.Reads) > 0 {
			readDBList := readsDBMap[cluster.Db]
			if readDBList == nil || cap(readDBList) < len(cluster.DBServer.Reads) {
				readDBList = make([]*gorm.DB, len(cluster.DBServer.Reads))
				readDBList[0] = nil
				readsDBMap[cluster.Db] = readDBList
			}
			for idx, conn := range cluster.DBServer.Reads {
				readDB := readDBList[idx]
				if readDB == nil || readDB.DB().Ping() != nil {
					if readDB != nil {
						readDB.Close()
						log.Fatalln("Read数据库异常断开")
					}
					if db, err := initDataPool(conn, cluster.DBServer.Pool, cluster.Db); err == nil {
						readDBList[idx] = db
					} else {
						m.err = err
						readDBList[idx] = nil
						log.Fatalln("Read数据库启动失败", err.Error())
					}
				}
			}
		}
	}
}

func (m *MysqlManager) InitDB() error {
	mysqlEnable, err := GetConfigSwitch("mysql")
	if err != nil {
		return err
	}

	if !mysqlEnable {
		return errors.New("configSwitch mysql is false!")
	}

	m.startMonitor()
	return m.err
}

func getMysqlConf(path string) (mysqlconf MysqlConf) {
	yamlfile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("read mysqlDriver yaml file failed:", path)
	}

	umerr := yaml.Unmarshal(yamlfile, &mysqlconf)
	if umerr != nil {
		log.Fatalln("mysqlDriver yaml Unmarshal failed:", umerr)
	}

	return mysqlconf
}

func initDataPool(conn MysqlConnection, pool MysqlPool, dbname string) (db *gorm.DB, err error) {
	dnsconf := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4,utf8&parseTime=True&loc=Local", conn.Username,
		conn.Password, conn.Host, conn.Port, dbname)

	db, err = gorm.Open("mysql", dnsconf)
	if err != nil {
		log.Fatalln("Mysql initDataPool open error:", err)
		return db, err
	}
	db.DB().SetMaxOpenConns(pool.MaxActive)
	db.DB().SetMaxIdleConns(pool.MaxIdle)
	db.DB().SetConnMaxLifetime(pool.MaxWait * time.Second)

	pingerr := db.DB().Ping()

	if pingerr != nil {
		log.Fatalln("Mysql initDataPool ping error:", err)
		return db, err
	}

	if AppIsDev() {
		db.LogMode(true)
	} else {
		db.LogMode(false)
	}

	return db, err
}

func CloseDB(db *gorm.DB) error {
	return db.Close()
}

//#################### public api ####################
func NewMysql(table string) *Mysql {
	return &Mysql{Table: table, Logger: MysqlLogger()}
}

func (m *Mysql) MysqlLogger() *logrus.Logger {
	return MysqlLogger()
}

func (m *Mysql) getDBName() string {
	if len(m.Database) > 0 {
		return m.Database
	} else {
		for _, cluster := range mysqlconf.Cluster {
			if len(cluster.Db) > 0 {
				return cluster.Db
			}
		}
	}
	return ""
}

func (m *Mysql) WriteDB() *gorm.DB {
	if len(writeDBMap) == 0 {
		return nil
	}
	dbname := m.getDBName()
	lock := m.getLock()
	lock.Lock()
	db := writeDBMap[dbname]
	lock.Unlock()
	return db
}

func (m *Mysql) ReadDB() *gorm.DB {
	dbname := m.getDBName()
	readDBSlice := readsDBMap[dbname]
	if len(readDBSlice) == 0 {
		//如果没有配置读写分离,统一按照write进行读写
		return m.WriteDB()
	}
	var readIndex int
	if len(readDBSlice) > 1 {
		rand.Seed(time.Now().Unix())
		readIndex = rand.Intn(len(readDBSlice))
	} else {
		readIndex = 0
	}
	lock := m.getLock()
	lock.RLock()
	readDB := readDBSlice[readIndex]
	lock.RUnlock()
	if readDB != nil {
		return readDB
	}
	//如果对应的没有取到, 顺序取第一个有值的
	for _, db := range readDBSlice {
		if db != nil {
			return db
		}
	}
	//如果没有配置读写分离,统一按照write进行读写
	return m.WriteDB()
}

func (m *Mysql) getLock() *sync.RWMutex {
	//key := m.getDBName()+"_"+m.Table
	key := m.getDBName() //每个数据库一个锁
	lock := dbLockMap[key]
	if lock == nil {
		lock = &sync.RWMutex{}
		dbLockMap[key] = lock
	}
	return lock
}

func (m *Mysql) Write(f func(db *gorm.DB) *gorm.DB) (*gorm.DB, error) {
	writeDB := m.WriteDB()
	if writeDB == nil {
		return nil, errors.New("WriteDB获取失败")
	}
	db := f(writeDB)
	var err error = nil
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			return db, nil
		}
		err = db.Error
	}
	return db, err
}

func (m *Mysql) Read(f func(db *gorm.DB) *gorm.DB) (*gorm.DB, error) {
	readDB := m.ReadDB()
	if readDB == nil {
		return nil, errors.New("ReadDB获取失败")
	}
	db := f(readDB)
	var err error = nil
	if db.Error != nil {
		if db.Error == gorm.ErrRecordNotFound {
			return db, nil
		}
		err = db.Error
	}
	return db, err
}
