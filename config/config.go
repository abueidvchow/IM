package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var Conf = new(AppConfig)

type AppConfig struct {
	Name              string `mapstructure:"name"`
	Mode              string `mapstructure:"mode"`
	Version           string `mapstructure:"version"`
	StartTime         string `mapstructure:"start_time"`         // 项目开始时间
	MachineID         int64  `mapstructure:"machine_id"`         // 雪花算法会使用到
	IP                string `mapstructure:"ip"`                 // 运行地址
	Port              int    `mapstructure:"port"`               // 运行端口
	HeartbeatTimeout  int    `mapstructure:"heartbeat_timeout"`  // 心跳超时时间（秒）
	HeartbeatInterval int    `mapstructure:"heartbeat_interval"` // 心跳检测时间间隔（秒）
	WorkerPoolSize    int32  `mapstructure:"worker_pool_size"`   // 队列数量
	MaxWorkerTask     int    `mapstructure:"max_worker_task"`    // 任务队列最大任务存储数量

	WebSocketConfig *WebSocketConfig `mapstructure:"websocket"`
	RabbitMQConfig  *RabbitMQConfig  `mapstructure:"rabbitmq"`
	LogConfig       *LogConfig       `mapstructure:"log"`
	MySQLConfig     *MySQLConfig     `mapstructure:"mysql"`
	RedisConfig     *RedisConfig     `mapstructure:"redis"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	UserName     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Addr         string `mapstructure:"addr"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	Compress   bool   `mapstructure:"compress"`
}

type WebSocketConfig struct {
	Port string `mapstructure:"port"`
}

type RabbitMQConfig struct {
	Url string `mapstructure:"url"`
}

func Init(configFile string) error {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	//把读取到的配置文件信息反序列化到Conf结构体
	if err = viper.Unmarshal(Conf); err != nil {
		fmt.Println("viper反序列化失败：", err)
		return err
	}
	return nil
}
