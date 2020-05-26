package cltranslator

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"strings"
)

const (
	X_LOCALE string = "X-LOCALE"
)

type Storage interface {
	GetLocalizedMessage(string, string) string
	GetLocalizedConfig(string, string) string
}

type TranslatorCfg struct {
	Enabled       bool   `mapstructure:"enabled" json:"enabled"`
	DefaultLocale string `mapstructure:"default_locale" json:"default_locale"`
}

type Translator interface {
	Translate(context.Context, string, ...string) string
	TranslateConfig(context.Context, string, ...string) string
	SetLocale(http.Handler) http.Handler
}

type translator struct {
	cfg     *TranslatorCfg
	storage Storage
}

func NewTranslator(cfg *TranslatorCfg, storage Storage) Translator {
	return &translator{
		cfg:     cfg,
		storage: storage,
	}
}

func (t *translator) Translate(ctx context.Context, messageID string, argKvs ...string) string {

	translatedString := messageID

	if t.cfg.Enabled {

		locale, _ := ctx.Value(X_LOCALE).(string)
		localizedString := t.storage.GetLocalizedMessage(messageID, locale)
		if localizedString != "" {
			translatedString = localizedString
		}
	}

	if len(argKvs) != 0 {
		argsMap := make(map[string]string)
		previousKey := ""
		for _, v := range argKvs {
			if previousKey != "" {
				argsMap[previousKey] = v
			}
			previousKey = v
		}

		for k, v := range argsMap {
			translatedString = strings.Replace(translatedString, fmt.Sprintf("{{var_%s}}", k), v, -1)
		}

	}
	return translatedString
}

func (t *translator) TranslateConfig(ctx context.Context, messageID string, argKvs ...string) string {
	if !t.cfg.Enabled {
		return messageID
	}

	locale, _ := ctx.Value(X_LOCALE).(string)
	localizedString := t.storage.GetLocalizedConfig(messageID, locale)
	if localizedString == "" {
		return messageID
	}
	if len(argKvs) != 0 {
		argsMap := make(map[string]string)
		previousKey := ""
		for _, v := range argKvs {
			if previousKey != "" {
				argsMap[previousKey] = v
			}
			previousKey = v
		}

		for k, v := range argsMap {
			localizedString = strings.Replace(localizedString, fmt.Sprintf("{{var_%s}}", k), v, -1)
		}
	}
	return localizedString
}

func (t *translator) SetLocale(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		var langSpec = ""
		langSpec = chi.URLParam(r, "locale")

		if langSpec == "" {
			langSpec = r.Header.Get(X_LOCALE)
		}

		if langSpec == "" {
			langSpec = t.cfg.DefaultLocale
		}

		ctx := context.WithValue(r.Context(), X_LOCALE, strings.ToLower(langSpec))
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func WithTranslationContext(ctx context.Context, locale string) context.Context {
	return context.WithValue(ctx, X_LOCALE, locale)
}
