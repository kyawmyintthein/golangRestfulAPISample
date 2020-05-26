package clconfigmanager

import (
	"bitbucket.org/libertywireless/circles-framework/cllogging"
	"bitbucket.org/libertywireless/circles-framework/cloption"
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/coreos/etcd/clientv3"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	LocaleLabel         string = "locale"
	DefaultLocaleRegexp string = `{{loc_([\.\d\w]+)}}`
)

type ConfigurationManager interface {
	/*
	 * LoadConfig
	 *     Load configuration files and localization files according to CCM config and
	 *     store in memory
	 */
	LoadConfig(context.Context, ...cloption.Option) error

	/*
	 * Scan
	 *     Unmarshal configuration data to struct or map
	 */
	Scan(interface{}) error

	/*
	 * Sync
	 *     Read configuration and localization data from ETCD and listen realtime changes.
	 *	   Then, apply the realtime changes in memory and also automatically unmarshall to struct or map
	 *     which is used before
	 */
	Sync(context.Context) error

	// Getter
	/*
	 * GetValue
	 *     Get value (type Value) by key from configuration storage. The nested key can be defined by dot (.)
	 *	   For example; key1.key2.key3
	 *     Value type can be convert to any specific type such as string, int, bool and etc...
	 *     For example;
	 *        value := GetValue("key1.key2.key3")
	 *        strValue := value.String()
	 */
	GetValue(string) Value

	/*
	 * GetLocalizedValue
	 *     Get value (type Value) by key (with {{loc_ }}) and locale from locale storage. The nested key can be defined by dot (.)
	 *	   For example; key1.key2.key3
	 *     Value type can be convert to any specific type such as string, int, bool and etc...
	 *     For example;
	 *        value := GetValue("key1.key2.key3", "en-US")
	 *        strValue := value.String()
	 */
	GetLocalizedValue(string, string) Value

	/*
	 * GetLocalizedConfig
	 *     Get value (type string) by key and locale from configuration and locale storage. The nested key can be defined by dot (.)
	 *     The key should be configuration key and the value should be locale key. This function will find the key in configuration storage and
	 *     get value (which is locale key). Then, it is find locale value in locale storage.
	 *	   For example:
	 *			Configuration storage
	 *          {
	 *				"example": "{{loc_example.name}}"
	 *			}
	 *
	 *
	 *			Locale storage
	 *          {
	 *				"en-US":
	 *				{
	 *					example: {
	 *						"name": "Example"
	 *					}
	 *				}
	 *			}
	 *
	 * 			localeString := GetLocalizedConfig("example", "en-US")
	 *
	 */
	GetLocalizedConfig(string, string) string // look up

	/*
	 * GetLocalizedString
	 *     Get value (type string) by key (with {{loc_ }}) and locale from locale storage. The nested key can be defined by dot (.)
	 *	   For example:
	 * 			localeString := GetLocalizedString("example.name", "en-US")
	 *
	 */
	GetLocalizedString(string, string) string

	/*
	 * GetLocalizedMessage
	 *     Get value (type string) by key (without {{loc_ }}) and locale from locale storage. The nested key can be defined by dot (.)
	 *	   For example:
	 * 			localeString := GetLocalizedString("example.name", "en-US")
	 *
	 */
	GetLocalizedMessage(string, string) string

	/*
	 * GetLocalizedMessageAsValue
	 *     Get value (type Value) by key (without {{loc_ }}) and locale from locale storage. The nested key can be defined by dot (.)
	 *	   For example:
	 * 			localeString := GetLocalizedString("example.name", "en-US")
	 *
	 */
	GetLocalizedMessageAsValue(string, string) Value

	/*
	 * WatchKey
	 *     Watch specific key changes and register callback to execute
	 *
	 */
	WatchKey(...string) KeyWatcher

	/*
	 * WatchLabel
	 *     Watch specific label changes and register callback to execute
	 *
	 */
	WatchLabel(string) KeyWatcher

	GetAllLocale() map[string]interface{}

	GetAllConfig() map[string]interface{}

	// Setter
	SetLogger(cllogging.KVLogger)
	SetConfig(*CCMCfg)
}
type ConfigurationSource map[string]interface{}
type configManager struct {
	mux sync.RWMutex
	cfg *CCMCfg

	localeRegexp *regexp.Regexp
	// services
	logger            cllogging.KVLogger
	configFile        string
	appConfigPaths    []string
	sharedConfigPaths []string
	localePaths       []string
	sharedLocalePaths []string
	ignoreDirs        []interface{}

	etcdClient     *clientv3.Client
	configStruct   interface{}
	configStorage  *viper.Viper
	localeStorage  *viper.Viper
	eventListeners EventListeners
}

func NewConfigManager(ccmCfg *CCMCfg, configFile string, opts ...cloption.Option) (ConfigurationManager, error) {
	configManager, err := newConfigManager(ccmCfg, configFile, opts...)
	if err != nil {
		return configManager, err
	}
	configManager.cfg = ccmCfg
	return configManager, nil
}

func newConfigManager(configStruct interface{}, configFile string, opts ...cloption.Option) (*configManager, error) {
	options := cloption.NewOptions(opts...)
	configManager := configManager{
		configFile:     configFile,
		localeRegexp:   regexp.MustCompile(DefaultLocaleRegexp),
		eventListeners: NewEventListeners(),
		cfg:            &CCMCfg{},
	}

	configManager.mux.Lock()
	defer configManager.mux.Unlock()

	// locale
	configManager.localeStorage = viper.New()

	// config
	configManager.configStorage = viper.New()
	configManager.configStorage.SetConfigFile(configFile)
	err := configManager.configStorage.ReadInConfig()
	if err != nil {
		return &configManager, err
	}

	err = configManager.scan(configStruct)
	if err != nil {
		return &configManager, err
	}

	// set logger
	lgr, ok := options.Context.Value(loggerKey{}).(cllogging.KVLogger)
	if lgr != nil && ok {
		configManager.logger = lgr
	} else {
		configManager.logger = cllogging.DefaultKVLogger()
	}

	return &configManager, nil
}

func (configManager *configManager) LoadConfig(ctx context.Context, opts ...cloption.Option) error {
	options := cloption.NewOptions(opts...)
	configManager.mux.Lock()
	defer configManager.mux.Unlock()

	// set App Config Paths
	ignoreDirs, ok := options.Context.Value(ignoreDirs{}).([]interface{})
	if len(ignoreDirs) != 0 && ok {
		configManager.ignoreDirs = ignoreDirs
	} else {
		configManager.ignoreDirs = []interface{}{}
	}

	// set cfg
	cfg, ok := options.Context.Value(ccmCfgKey{}).(*CCMCfg)
	if cfg == nil && !ok {
		cfg = configManager.cfg
	}

	// load files from CCMCfg
	configManager.loadConfigFilesFromCCMConfig(ctx, cfg)
	// load files from optional parameters
	configManager.loadConfigFilesFromOptionalParameters(options, cfg)

	for _, sharedConfigFile := range configManager.sharedConfigPaths {
		f, err := os.Open(sharedConfigFile)
		if err != nil {
			return err
		}
		err = configManager.configStorage.MergeConfig(f)
		if err != nil {
			return err
		}
		f.Close()
	}

	for _, appConfigFile := range configManager.appConfigPaths {
		f, err := os.Open(appConfigFile)
		if err != nil {
			return err
		}
		err = configManager.configStorage.MergeConfig(f)
		if err != nil {
			return err
		}
		f.Close()
	}

	primaryLocaleFilePath := ""
	if len(configManager.localePaths) != 0 {
		primaryLocaleFilePath = configManager.localePaths[0]
	} else if len(configManager.sharedLocalePaths) != 0 {
		primaryLocaleFilePath = configManager.sharedLocalePaths[0]
	} else {
		return nil
	}

	configManager.localeStorage.SetConfigFile(primaryLocaleFilePath)
	err := configManager.localeStorage.ReadInConfig()
	if err != nil {
		return err
	}

	for _, sharedLocalePath := range configManager.sharedLocalePaths {
		f, err := os.Open(sharedLocalePath)
		if err != nil {
			return err
		}
		err = configManager.localeStorage.MergeConfig(f)
		if err != nil {
			return err
		}
		f.Close()
	}

	for _, localeConfigFile := range configManager.localePaths {
		f, err := os.Open(localeConfigFile)
		if err != nil {
			return err
		}
		err = configManager.localeStorage.MergeConfig(f)
		if err != nil {
			return err
		}
		f.Close()
	}

	return nil
}

func (configManager *configManager) loadConfigFilesFromCCMConfig(ctx context.Context, cfg *CCMCfg) {
	// config
	extraConfigFiles := configManager.registerAppConfigPaths(ctx, &cfg.Config.AppConfig, cfg.Config.AppConfig.AppFilePaths)
	sharedConfigFiles := configManager.registerConfigPaths(ctx, &cfg.Config.CirclesConfig, cfg.Config.CirclesConfig.SharedFilePaths)
	appConfigFiles := configManager.registerConfigPaths(ctx, &cfg.Config.CirclesConfig, cfg.Config.CirclesConfig.AppFilePaths)

	sharedConfigFilesFromDir := configManager.registerConfigFilesFromDirectories(ctx, &cfg.Config.CirclesConfig, cfg.Config.CirclesConfig.SharedDirectories)
	appConfigFilesFromDir := configManager.registerConfigFilesFromDirectories(ctx, &cfg.Config.CirclesConfig, cfg.Config.CirclesConfig.AppDirectories)

	configManager.appConfigPaths = append(configManager.appConfigPaths, extraConfigFiles...)
	configManager.sharedConfigPaths = append(configManager.sharedConfigPaths, sharedConfigFilesFromDir...)
	configManager.sharedConfigPaths = append(configManager.sharedConfigPaths, sharedConfigFiles...)
	configManager.appConfigPaths = append(configManager.appConfigPaths, appConfigFilesFromDir...)
	configManager.appConfigPaths = append(configManager.appConfigPaths, appConfigFiles...)

	// locale
	sharedLocaleFiles := configManager.registerConfigPaths(ctx, &cfg.Config.CirclesLocale, cfg.Config.CirclesLocale.SharedFilePaths)
	appLocaleFiles := configManager.registerConfigPaths(ctx, &cfg.Config.CirclesLocale, cfg.Config.CirclesLocale.AppFilePaths)

	sharedLocaleFilesFromDir := configManager.registerConfigFilesFromDirectories(ctx, &cfg.Config.CirclesLocale, cfg.Config.CirclesLocale.SharedDirectories)
	appLocaleFilesFromDir := configManager.registerConfigFilesFromDirectories(ctx, &cfg.Config.CirclesLocale, cfg.Config.CirclesLocale.AppDirectories)

	configManager.sharedLocalePaths = append(configManager.sharedLocalePaths, sharedLocaleFilesFromDir...)
	configManager.sharedLocalePaths = append(configManager.sharedLocalePaths, sharedLocaleFiles...)
	configManager.localePaths = append(configManager.localePaths, appLocaleFilesFromDir...)
	configManager.localePaths = append(configManager.localePaths, appLocaleFiles...)

}

func (configManager *configManager) loadConfigFilesFromOptionalParameters(options cloption.Options, cfg *CCMCfg) {
	// set App Config Paths
	appConfigPaths, ok := options.Context.Value(appConfigFilePaths{}).([]string)
	if len(appConfigPaths) != 0 && ok {
		for _, v := range appConfigPaths {
			if v == "" {
				continue
			}
			absPath := filepath.Join(cfg.Config.CirclesConfig.BaseDir, v)
			configManager.appConfigPaths = append(configManager.appConfigPaths, absPath)
		}
	}

	// set Shared Config Path
	sharedConfigPaths, ok := options.Context.Value(sharedConfigFilePaths{}).([]string)
	if len(sharedConfigPaths) != 0 && ok {
		for _, v := range sharedConfigPaths {
			if v == "" {
				continue
			}
			absPath := filepath.Join(cfg.Config.CirclesConfig.BaseDir, v)
			configManager.sharedConfigPaths = append(configManager.sharedConfigPaths, absPath)
		}
	}

	// set Locale Config Path
	localeConfigPaths, ok := options.Context.Value(localeConfigFilePaths{}).([]string)
	if len(localeConfigPaths) != 0 && ok {
		for _, v := range localeConfigPaths {
			if v == "" {
				continue
			}
			absPath := filepath.Join(cfg.Config.CirclesLocale.BaseDir, v)
			configManager.localePaths = append(configManager.localePaths, absPath)
		}
	}

	// set Locale Config Path
	sharedLocaleFilePaths, ok := options.Context.Value(sharedLocaleFilePaths{}).([]string)
	if len(sharedLocaleFilePaths) != 0 && ok {
		for _, v := range sharedLocaleFilePaths {
			if v == "" {
				continue
			}
			absPath := filepath.Join(cfg.Config.CirclesLocale.BaseDir, v)
			configManager.sharedLocalePaths = append(configManager.sharedLocalePaths, absPath)
		}
	}

	var sharedLocaleFiles, localeFiles []string
	// get shared config dirs
	sharedLocaleDirs, ok := options.Context.Value(sharedLocaleConfigDirs{}).([]string)
	if len(sharedLocaleDirs) != 0 && ok {
		for _, sharedLocaleDir := range sharedLocaleDirs {
			if sharedLocaleDir == "" {
				continue
			}
			absDir := path.Join(cfg.Config.CirclesLocale.BaseDir, sharedLocaleDir)
			filepath.Walk(absDir, func(path string, f os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if f.IsDir() && sliceContains(configManager.ignoreDirs, f.Name()) {
					return filepath.SkipDir
				}

				ext := filepath.Ext(path)
				if f.IsDir() || ext != fmt.Sprintf(".%s", cfg.Config.CirclesLocale.FileType) {
					return nil
				}

				sharedLocaleFiles = append(sharedLocaleFiles, path)
				return nil
			})
		}
	}
	configManager.sharedLocalePaths = append(configManager.sharedLocalePaths, sharedLocaleFiles...)

	localeConfigDirs, ok := options.Context.Value(localeConfigDirs{}).([]string)
	if len(localeConfigDirs) != 0 {
		for _, localeConfigDir := range localeConfigDirs {
			if localeConfigDir == "" {
				continue
			}
			absDir := path.Join(cfg.Config.CirclesLocale.BaseDir, localeConfigDir)
			filepath.Walk(absDir, func(path string, f os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if f.IsDir() && sliceContains(configManager.ignoreDirs, f.Name()) {
					return filepath.SkipDir
				}
				ext := filepath.Ext(path)
				if f.IsDir() || ext != fmt.Sprintf(".%s", cfg.Config.CirclesLocale.FileType) {
					return nil
				}

				localeFiles = append(localeFiles, path)
				return nil
			})
		}
	}
	configManager.localePaths = append(configManager.localePaths, localeFiles...)

}

func (configManager *configManager) registerAppConfigPaths(ctx context.Context, cfg *AppFileCfg, filePaths []string) []string {
	var configFilePaths []string
	if len(filePaths) != 0 {
		for _, v := range filePaths {
			if v == "" {
				continue
			}
			absPath := filepath.Join(cfg.BaseDir, v)
			configFilePaths = append(configFilePaths, absPath)
		}
	}
	return configFilePaths
}

func (configManager *configManager) registerConfigPaths(ctx context.Context, cfg *FileCfg, filePaths []string) []string {
	var configFilePaths []string
	if len(filePaths) != 0 {
		for _, v := range filePaths {
			if v == "" {
				continue
			}
			absPath := filepath.Join(cfg.BaseDir, v)
			configFilePaths = append(configFilePaths, absPath)
		}
	}
	return configFilePaths
}

func (configManager *configManager) registerConfigFilesFromDirectories(ctx context.Context, fileCfg *FileCfg, directories []string) []string {
	var configFilePaths []string
	if len(directories) != 0 {
		for _, dir := range directories {
			if dir == "" {
				continue
			}
			absDir := path.Join(fileCfg.BaseDir, dir)
			filepath.Walk(absDir, func(path string, f os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if f.IsDir() && sliceContains(configManager.ignoreDirs, f.Name()) {
					return filepath.SkipDir
				}

				ext := filepath.Ext(path)
				if f.IsDir() || ext != fmt.Sprintf(".%s", fileCfg.FileType) {
					return nil
				}

				configFilePaths = append(configFilePaths, path)
				return nil
			})
		}
	}
	return configFilePaths
}

func (configManager *configManager) Sync(ctx context.Context) error {

	if configManager.cfg == nil {
		configManager.logger.Warn("[CONF-MANAGER] disabled ETCD.")
		return nil
	}

	if !configManager.cfg.ETCD.Enabled {
		configManager.logger.Warn("[CONF-MANAGER] disabled ETCD.")
		return nil
	}

	err := configManager.initETCDClient()
	if err != nil {
		configManager.logger.Error(err, "[CONF-MANAGER] failed to init ETCD connection.")
		if !configManager.cfg.ETCD.SkipErrorOnETCDConnFailed {
			return err
		}
		configManager.cfg.ETCD.Enabled = false
		configManager.logger.Error(err, "[CONF-MANAGER] skip ETCD connection error and disabled ETCD.")
	}

	if configManager.etcdClient == nil || !configManager.cfg.ETCD.Enabled {
		configManager.logger.Warn("[CONF-MANAGER] disabled runtime watcher: failed to init ETCD client.")
		return nil
	}

	if configManager.cfg.ETCD.WatcherPath == "" {
		configManager.logger.Warn("[CONF-MANAGER] disabled runtime watcher: missing watcher path.")
		return nil
	}

	// get configuration data from ETCD
	for _, label := range configManager.cfg.ETCD.ConfigurationLabels {
		// get configuration data for each label to store them in separate map
		labelPath := fmt.Sprintf("%s/%s", configManager.cfg.ETCD.WatcherPath, label)
		getResp, err := configManager.etcdClient.Get(ctx, labelPath, clientv3.WithPrefix())
		if err != nil {
			return err
		}

		for _, kv := range getResp.Kvs {
			// unmarshal value []byte to interface
			var val interface{}
			err = json.Unmarshal(kv.Value, &val)
			if err != nil {
				val = string(kv.Value)
			}

			// get actual configuration keys from ETCD key without watcher key
			configKeyStr := strings.Replace(string(kv.Key), fmt.Sprintf("%s/", labelPath), "", -1)
			nestedConfigKey := strings.Join(strings.Split(strings.Trim(configKeyStr, ""), "/"), ".")

			if label == LocaleLabel {
				configManager.localeStorage.Set(nestedConfigKey, val)
			} else {
				configManager.configStorage.Set(nestedConfigKey, val)
			}
		}
	}

	if configManager.cfg.ETCD.Enabled {
		configManager.watch(ctx, configManager.cfg.ETCD.WatcherPath)
	}

	return nil
}

func (configManager *configManager) initETCDClient() error {
	etcdCfg := configManager.cfg.ETCD
	// start etcd client service
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdCfg.Endpoints,
		DialTimeout: etcdCfg.DialTimeout * time.Second,
		Username:    etcdCfg.Username,
		Password:    etcdCfg.Password,
	})
	if err != nil {
		return err
	}
	configManager.etcdClient = etcdClient
	return nil
}

func (configManager *configManager) Scan(configStructPtr interface{}) error {
	configManager.mux.Lock()
	defer configManager.mux.Unlock()
	configManager.scan(configStructPtr)
	configManager.configStruct = configStructPtr
	return nil
}

func (configManager *configManager) scan(configStructPtr interface{}) error {
	valueOfIStructPointer := reflect.ValueOf(configStructPtr)

	if k := valueOfIStructPointer.Kind(); k != reflect.Ptr {
		return fmt.Errorf("config should be pointer type.")
	}

	valueOfIStructPointerElem := valueOfIStructPointer.Elem()

	// Below is a further (and definitive) check regarding settability in addition to checking whether it is a pointer earlier.
	if !valueOfIStructPointerElem.CanSet() {
		return fmt.Errorf("unable to set value to type!")
	}

	err := configManager.configStorage.Unmarshal(configStructPtr)
	if err != nil {
		return err
	}

	return nil
}

// convert keys to etcd key
// register key in cache
//
func (configManager *configManager) WatchKey(keys ...string) KeyWatcher {
	keyPath := strings.Join(keys, "/")
	keyWatcher := NewKeyWatcher(keyPath)
	configManager.eventListeners.Store(keyPath, keyWatcher)
	return keyWatcher
}

// convert keys to etcd key
// register key in cache
//
func (configManager *configManager) WatchLabel(label string) KeyWatcher {
	keyWatcher := NewKeyWatcher(label)
	configManager.eventListeners.Store(label, keyWatcher)
	return keyWatcher
}

// WatchfromEtcd watch etcd key and sync into config struct.
// It wil also call reload callback function to reinitalize the module.
func (configManager *configManager) watch(ctx context.Context, watcherPath string) {
	configManager.logger.Info("[CONF-MANAGER] ETCD runtime changes watcher is stared")
	go func() {
		watchChan := configManager.etcdClient.Watch(ctx, watcherPath, clientv3.WithPrefix())
		for true {
			select {
			case result := <-watchChan:
				updatedConfig := make(map[string][]byte)
				for _, ev := range result.Events {
					configManager.logger.Infof("Event : %s:%s \n", string(ev.Kv.Key), string(ev.Kv.Value))
					updatedConfig[string(ev.Kv.Key)] = ev.Kv.Value
				}
				err := configManager.updateRuntimeChanges(ctx, updatedConfig)
				if err != nil {
					configManager.logger.Errorf(err, "failed to update values on")
				}

			}
		}
	}()
}

func (configManager *configManager) updateRuntimeChanges(ctx context.Context, eventKeyValue map[string][]byte) error {
	configManager.mux.Lock()
	defer configManager.mux.Unlock()

	updatedConfig := make(map[string][]byte)
	for k, v := range eventKeyValue {
		// unmarshal value []byte to interface
		var val interface{}
		err := json.Unmarshal(v, &val)
		if err != nil {
			val = string(v)
		}

		// get actual configuration keys from ETCD key without watcher key
		configKeyStr := strings.Replace(k, fmt.Sprintf("%s/", configManager.cfg.ETCD.WatcherPath), "", -1)
		configKeysWithLabel := strings.Split(strings.Trim(configKeyStr, ""), "/")
		if len(configKeysWithLabel) < 2 {
			err = errors.New("invalid configuration key")
			configManager.logger.Errorf(err, "invalid configuration key format : %v", k)
			continue
		}
		configKeys := configKeysWithLabel[1:len(configKeysWithLabel)]
		nestedConfigKey := strings.Join(configKeys, ".")
		updatedConfig[nestedConfigKey] = v

		// execute callbacks for label
		label := configKeysWithLabel[0]
		if label == LocaleLabel {
			configManager.localeStorage.Set(nestedConfigKey, val)
		} else {
			configManager.configStorage.Set(nestedConfigKey, val)
		}

		keyWatcher, ok := configManager.eventListeners.Get(label)
		if ok {
			keyWatcher.Execute(ctx, nestedConfigKey, NewValue(v))
		}
	}

	err := configManager.configStorage.Unmarshal(configManager.configStruct)
	if err != nil {
		return err
	}

	// execute callbacks for updated key
	for k, v := range updatedConfig {
		keyWatcher, ok := configManager.eventListeners.Get(k)
		if ok {
			keyWatcher.Execute(ctx, k, NewValue(v))
		}
	}

	return nil
}

func (configManager *configManager) GetValue(key string) Value {
	return NewValueViaInterface(configManager.configStorage.Get(key))
}

func (configManager *configManager) GetLocalizedValue(key string, locale string) Value {
	localizedValue := key
	if !configManager.localeRegexp.MatchString(localizedValue) {
		return NewValueViaInterface(localizedValue)
	}

	localeKey := strings.TrimPrefix(localizedValue, "{{loc_")
	localeKey = strings.TrimSuffix(localeKey, "}}")
	return NewValueViaInterface(configManager.localeStorage.Get(fmt.Sprintf("%s.%s", locale, localeKey)))
}

func (configManager *configManager) getLocalizedString(key string, locale string) string {
	localeValue := configManager.localeStorage.GetString(fmt.Sprintf("%s.%s", locale, key))
	if localeValue == "" {
		localeValue = key
	}
	return localeValue
}

func (configManager *configManager) GetLocalizedMessage(key string, locale string) string {
	localeValue := configManager.localeStorage.GetString(fmt.Sprintf("%s.%s", locale, key))
	if localeValue == "" {
		localeValue = key
	}
	return localeValue
}

func (configManager *configManager) GetLocalizedMessageAsValue(key string, locale string) Value {
	return NewValueViaInterface(configManager.localeStorage.Get(fmt.Sprintf("%s.%s", locale, key)))
}

func (configManager *configManager) GetLocalizedString(key string, locale string) string {
	localizedValue := key
	if !configManager.localeRegexp.MatchString(localizedValue) {
		return localizedValue
	}

	localeKey := strings.TrimPrefix(localizedValue, "{{loc_")
	localeKey = strings.TrimSuffix(localeKey, "}}")
	localizedValue = configManager.getLocalizedString(localeKey, locale)
	return localizedValue
}

func (configManager *configManager) GetLocalizedConfig(key string, locale string) string {
	configValue := configManager.configStorage.GetString(key)
	if !configManager.cfg.EnabledConfigLocalization {
		return configValue
	}

	if !configManager.localeRegexp.MatchString(configValue) {
		return configValue
	}

	localeKey := strings.TrimPrefix(configValue, "{{loc_")
	localeKey = strings.TrimSuffix(localeKey, "}}")
	localizedValue := configManager.getLocalizedString(localeKey, locale)
	if localizedValue == "" {
		return configValue
	}
	return localizedValue
}

func Unmarshal(ctx context.Context, filepath string, configStructPtr interface{}) error {
	valueOfIStructPointer := reflect.ValueOf(configStructPtr)

	if k := valueOfIStructPointer.Kind(); k != reflect.Ptr {
		return fmt.Errorf("config should be pointer type.")
	}

	valueOfIStructPointerElem := valueOfIStructPointer.Elem()

	// Below is a further (and definitive) check regarding settability in addition to checking whether it is a pointer earlier.
	if !valueOfIStructPointerElem.CanSet() {
		return fmt.Errorf("unable to set value to type!")
	}

	viperConfig := viper.New()
	viperConfig.SetConfigFile(filepath)
	err := viperConfig.ReadInConfig()
	if err != nil {
		return err
	}
	err = viperConfig.Unmarshal(configStructPtr)
	if err != nil {
		return err
	}
	return nil
}

func (configManager *configManager) GetAllConfig() map[string]interface{} {
	return configManager.configStorage.AllSettings()
}

func (configManager *configManager) GetAllLocale() map[string]interface{} {
	return configManager.localeStorage.AllSettings()
}

func (configManager *configManager) SetLogger(logger cllogging.KVLogger) {
	configManager.logger = logger
}

func (configManager *configManager) SetConfig(cfg *CCMCfg) {
	configManager.cfg = cfg
}

func GetBytes(val interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(val)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func sliceContains(slice interface{}, item interface{}) bool {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return false // given a non-slice type
	}
	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return true
		}
	}
	return false
}
