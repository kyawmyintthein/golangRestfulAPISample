package clconfigmanager

import (
	"bitbucket.org/libertywireless/circles-framework/cllogging"
	"bitbucket.org/libertywireless/circles-framework/cloption"
	"context"
)

type loggerKey struct{}

func WithLogger(a cllogging.KVLogger) cloption.Option {
	return func(o *cloption.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, loggerKey{}, a)
	}
}

type appConfigConfigDir struct{}

func WithAppConfigDir(dirs ...string) cloption.Option {
	return func(o *cloption.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, appConfigConfigDir{}, dirs)
	}
}

type appConfigFilePaths struct{}

func WithAppConfigPaths(filePaths ...string) cloption.Option {
	return func(o *cloption.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, appConfigFilePaths{}, filePaths)
	}
}

type sharedConfigFilePaths struct{}

func WithSharedConfigPaths(filePaths ...string) cloption.Option {
	return func(o *cloption.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, sharedConfigFilePaths{}, filePaths)
	}
}

type localeConfigFilePaths struct{}

func WithLocaleConfigPaths(filePaths ...string) cloption.Option {
	return func(o *cloption.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, localeConfigFilePaths{}, filePaths)
	}
}

type sharedLocaleFilePaths struct{}

func WithSharedLocalePaths(filePaths ...string) cloption.Option {
	return func(o *cloption.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, sharedLocaleFilePaths{}, filePaths)
	}
}

type sharedConfigDirs struct{}

func WithSharedDirs(dirs ...string) cloption.Option {
	return func(o *cloption.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, sharedConfigDirs{}, dirs)
	}
}

type localeConfigDirs struct{}

func WithLocaleDirs(dirs ...string) cloption.Option {
	return func(o *cloption.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, localeConfigDirs{}, dirs)
	}
}

type sharedLocaleConfigDirs struct{}

func WithSharedLocaleDirs(dirs ...string) cloption.Option {
	return func(o *cloption.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, sharedLocaleConfigDirs{}, dirs)
	}
}

type ignoreDirs struct{}

func WithIgnoreDirs(dirs ...interface{}) cloption.Option {
	return func(o *cloption.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, ignoreDirs{}, dirs)
	}
}

type localeFileTypeKey struct{}

func WithLocaleFileType(s string) cloption.Option {
	return func(o *cloption.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, localeFileTypeKey{}, s)
	}
}

type ccmCfgKey struct{}

func WithConfig(config *CCMCfg) cloption.Option {
	return func(o *cloption.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, ccmCfgKey{}, config)
	}
}
