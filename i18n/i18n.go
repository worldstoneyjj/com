package i18n

import (
	"embed"
	"encoding/json"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var localeFiles embed.FS

var (
	bundle     *i18n.Bundle
	once       sync.Once
	localizers map[string]*i18n.Localizer
)

// InitI18n 初始化并加载多语言文件
func initI18n() {
	once.Do(func() {
		bundle = i18n.NewBundle(language.English)
		bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

		localizers = make(map[string]*i18n.Localizer)

		// 读取嵌入的语言文件列表
		files, err := localeFiles.ReadDir("locales")
		if err != nil {
			panic(err)
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			// 获取语言代码，例如从 "en.json" 提取 "en"
			langCode := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			filePath := path.Join("locales", file.Name())
			data, err := localeFiles.ReadFile(filePath)
			if err != nil {
				panic(err)
			}
			// 解析语言文件
			bundle.MustParseMessageFileBytes(data, file.Name())
			// 创建对应的 Localizer
			localizers[langCode] = i18n.NewLocalizer(bundle, langCode)
		}
	})
}

func Translate(lang, messageID string, templateData map[string]interface{}) string {
	initI18n()
	localizer, ok := localizers[lang]
	if !ok {
		return messageID
	}
	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})
	if err != nil {
		return messageID // 如果未找到翻译，返回消息 ID 本身
	}
	return message
}
