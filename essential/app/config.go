package app

import (
	"poseidon/essential/endpoint"
	"poseidon/essential/storage"

	"encoding/json"

	"time"

	"context"

	"bytes"

	"io/ioutil"

	"github.com/coreos/etcd/clientv3"
	"github.com/fsnotify/fsnotify"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type ServerConfig struct {
	Host     string
	GrpcPort int
	HttpPort int
	PrometheusPort int
}

type KafkaConfig struct {
	BrokerList []string
}

type GrpcConfig struct {
	ServiceRegistry endpoint.ServiceRegistryOption
}
type MysqlConfig struct {
	ConnectionString map[string]string
}

type FtxConfig struct {
	Address string
	Mid     string
}

type RedisConfig struct {
	Addr     string
	Password string
	PoolSize int
}

type PublicConfig struct {
	ExcelSavePath string
	DownloadPath  string
}

type CORSConfig struct {
	AllowOrigins string
}

type TokenConfig struct {
	Secret       string
	ExpireTime   int64
	HcCheckPath  []string
	HcIgnorePath []string
	JwCheckPath  []string
	JwIgnorePath []string
}

type OauthConfig struct {
	ClientId       string
	ClientSecret   string
	ClientCallback string
}

type PasswordConfig struct {
	Secret string
}

type ZipkinConfig struct {
	Url                    string
	Debug                  bool
	Sampler                string
	Mod                    int
	ClientServerSameSample bool
}

type ReleaseConfig struct {
	Release string
}

type SessionConfig struct {
	ExpireTime int64
}

type LoggerConfig interface {
	GetLogLevel() string
}

type ServiceRegistryConfig interface {
	GetServiceRegistry() endpoint.ServiceRegistryOption
}

type TCPServerConfig interface {
	GetTCPServerConfig() ServerConfig
}

type MysqlServerConfig interface {
	GetMysqlConfig() MysqlConfig
}

type JwtTokenConfig interface {
	GetTokenConfig() TokenConfig
}

type KafkaServerConfig interface {
	GetKafkaConfig() KafkaConfig
}

type PublicServiceConfig interface {
	GetPublicConfig() PublicConfig
}

type UserPasswordConfig interface {
	GetPasswordConfig() PasswordConfig
}

type ZipkinServiceConfig interface {
	GetZipkinConfig() *ZipkinConfig
}

type RedisServiceConfig interface {
	GetRedisConfig() *storage.RedisOption
}

type SessionServiceConfig interface {
	GetSessionConfig() *SessionConfig
}

type UploadConfig struct {
	UploadPath        string
	UploadFieldName   string
	UploadType        int32
	UploadMaxFileSize string
}

type OssConfig struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	Bucket          string
	BindDomain      string
	Schema          string
	CacheLiveDays   int32
}

type ConfigWatcher func(config interface{}) error

type ConfigLoader struct {
	Name          string
	FilePath      string
	EtcdEndpoint  []string
	Config        interface{}
	ConfigWatcher ConfigWatcher
	rawData       string
}

func (loader *ConfigLoader) Load() error {

	if err := loader.LoadEtcdConfig(); err == nil {
		return nil
	}

	return loader.LoadConfig()
}

func (loader *ConfigLoader) LoadConfig() error {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("ConfigLoader ReadConfig failed, err: %s", err)
		return err
	}

	if err := viper.Unmarshal(loader.Config); err != nil {
		log.Errorf("ConfigLoader ReadConfig failed %s", err)
		return err
	}

	configData, _ := json.MarshalIndent(viper.AllSettings(), "", "  ")
	log.Debugf("ConfigLoader ReadRemoteConfig success, \n%s", configData)

	//viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Infof("ConfigLoader OnConfigChange op:%s, name:%s", in.Op, in.Name)

		if err := viper.Unmarshal(loader.Config); err != nil {
			log.Errorf("ConfigLoader OnConfigChange failed %s", err)
		} else {
			if loader.ConfigWatcher != nil {
				loader.ConfigWatcher(loader.Config)
			}
		}
	})

	return nil
}

func (loader *ConfigLoader) LoadEtcdConfig() error {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   loader.EtcdEndpoint,
		DialTimeout: time.Second * 3,
	})

	if err != nil {
		log.Errorf("ConfigLoader LoadEtcdConfig err:%s", err)
		return err
	}

	key := "/poseidon/config/" + loader.Name
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	response, err := etcdClient.KV.Get(ctx, key)
	if err != nil {
		log.Errorf("ConfigLoader LoadEtcdConfig err:%s", err)
		return err
	}

	if response.Count == 0 {
		log.Errorf("ConfigLoader LoadEtcdConfig response empty")
		return errors.Errorf("ConfigLoader LoadEtcdConfig response empty")

	}

	err = toml.Unmarshal(response.Kvs[0].Value, loader.Config)
	if err != nil {
		log.Errorf("ConfigLoader LoadEtcdConfig err:%s", err)
		return err
	}

	viper.SetConfigType("toml")
	if err := viper.ReadConfig(bytes.NewBuffer(response.Kvs[0].Value)); err != nil {
		log.Errorf("ConfigLoader LoadEtcdConfig err:%s", err)
		return err
	}

	configData, _ := json.MarshalIndent(viper.AllSettings(), "", "  ")
	log.Debugf("ConfigLoader ReadRemoteConfig success, \n%s", configData)

	watch := etcdClient.Watcher.Watch(context.Background(), key)

	go func() {
		for resp := range watch {
			log.Infof("ConfigLoader ReadRemoteConfig watcher resp on [%s]", resp)
			value := resp.Events[0].Kv.Value

			err = toml.Unmarshal(value, loader.Config)
			if err != nil {
				log.Errorf("ConfigLoader ReadRemoteConfig watcher err:%s", err)
				continue
			}

			configData, _ := json.MarshalIndent(loader.Config, "", "  ")
			log.Debugf("ConfigLoader ReadRemoteConfig watcher success, \n%s", configData)

			if err := viper.ReadConfig(bytes.NewBuffer(value)); err != nil {
				log.Errorf("ConfigLoader ReadRemoteConfig watcher err:%s", err)
				continue
			}

			if loader.ConfigWatcher != nil {
				loader.ConfigWatcher(loader.Config)
			}
		}
	}()
	return nil
}

func (loader *ConfigLoader) PushEtcdConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("ConfigLoader PushEtcdConfig read file err:%s", err)
		return err
	}

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   loader.EtcdEndpoint,
		DialTimeout: time.Second * 3,
	})
	if err != nil {
		log.Errorf("ConfigLoader PushEtcdConfig err:%s", err)
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	_, err = etcdClient.KV.Put(ctx, "/poseidon/config/"+loader.Name, string(data))
	if err != nil {
		log.Errorf("ConfigLoader LoadEtcdConfig err:%s", err)
		return err
	}

	return nil
}
