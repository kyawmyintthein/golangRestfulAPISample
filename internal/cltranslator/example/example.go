package main

import (
	clconfigmanager "bitbucket.org/libertywireless/circles-framework/clconfigmanager"
	"bitbucket.org/libertywireless/circles-framework/cltranslator"
	"context"
	"fmt"
)

func main() {

	cfg := &cltranslator.TranslatorCfg{
		Enabled:       true,
		DefaultLocale: "en-US",
	}

	ccmCfg := clconfigmanager.CCMCfg{}
	configManager, err := clconfigmanager.NewConfigManager(&ccmCfg, "./cltranslator/example/config/config.yml")
	if err != nil {
		panic(fmt.Sprintf("failed to start config manager : %v", err))
	}

	err = configManager.LoadConfig(context.Background())
	if err != nil {
		panic(fmt.Sprintf("failed to load configuration data : %v", err))
	}

	translator := cltranslator.NewTranslator(cfg, configManager)
	translatedMsg := translator.Translate(cltranslator.WithTranslationContext(context.Background(), "en-US"), "email_title", "first_name", "Kyaw", "last_name", "Myint Thein")
	fmt.Println("Translated message :", translatedMsg)

}
