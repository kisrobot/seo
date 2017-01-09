package seo

import (
	"reflect"
	"text/template"

	"github.com/qor/admin"
)

func seoSections(context *admin.Context, collection *Collection) []interface{} {
	settings := []interface{}{}
	for _, seo := range collection.registeredSeo {
		s := collection.SettingResource.NewStruct()
		db := context.GetDB()
		db.Where("name = ?", seo.Name).First(s)
		if db.NewRecord(s) {
			s.(QorSeoSettingInterface).SetName(seo.Name)
			s.(QorSeoSettingInterface).SetSeoType(seo.Name)
			db.Save(s)
		}
		s.(QorSeoSettingInterface).SetCollection(collection)
		settings = append(settings, s)
	}
	return settings
}

func seoSettingMetas(collection *Collection) []*admin.Section {
	return collection.SettingResource.EditAttrs()
}

func seoGlobalSetting(context *admin.Context, collection *Collection) interface{} {
	s := collection.SettingResource.NewStruct()
	db := context.GetDB()
	db.Where("is_global_seo = ? and name = ?", true, collection.GlobalSettingName).First(s)
	if db.NewRecord(s) {
		s.(QorSeoSettingInterface).SetName(collection.GlobalSettingName)
		s.(QorSeoSettingInterface).SetSeoType(collection.GlobalSettingName)
		s.(QorSeoSettingInterface).SetIsGlobalSeo(true)
		db.Save(s)
	}
	return s
}

func seoGlobalSettingValue(collection *Collection, setting map[string]string) interface{} {
	value := reflect.Indirect(reflect.ValueOf(collection.globalSetting))
	for i := 0; i < value.NumField(); i++ {
		fieldName := value.Type().Field(i).Name
		if setting[fieldName] != "" {
			value.Field(i).SetString(setting[fieldName])
		}
	}
	return value.Interface()
}

func seoGlobalSettingMetas(collection *Collection) []*admin.Section {
	return collection.globalResource.NewAttrs()
}

func seoTagsByType(seo *SEO) (tags []string) {
	if seo == nil {
		return []string{}
	}
	value := reflect.Indirect(reflect.ValueOf(seo.collection.globalSetting))
	for i := 0; i < value.NumField(); i++ {
		tags = append(tags, value.Type().Field(i).Name)
	}
	for _, s := range seo.Varibles {
		tags = append(tags, s)
	}
	return tags
}

func seoAppendDefaultValue(context *admin.Context, seo *SEO, resourceSeoValue interface{}) interface{} {
	db := context.GetDB()
	globalInteface := seo.collection.SettingResource.NewStruct()
	db.Where("name = ?", seo.Name).First(globalInteface)
	globalSetting := globalInteface.(QorSeoSettingInterface)
	setting := resourceSeoValue.(Setting)
	if !setting.EnabledCustomize && setting.Title == "" && setting.Description == "" && setting.Keywords == "" {
		setting.Title = globalSetting.GetTitle()
		setting.Description = globalSetting.GetDescription()
		setting.Keywords = globalSetting.GetKeywords()
	}
	return setting
}

func seoURL(collection *Collection, name string) string {
	return collection.SeoSettingURL(name)
}

func registerFuncMap(a *admin.Admin) {
	funcMaps := template.FuncMap{
		"seo_sections":             seoSections,
		"seo_setting_metas":        seoSettingMetas,
		"seo_global_setting_value": seoGlobalSettingValue,
		"seo_global_setting_metas": seoGlobalSettingMetas,
		"seo_global_setting":       seoGlobalSetting,
		"seo_tags_by_type":         seoTagsByType,
		"seo_append_default_value": seoAppendDefaultValue,
		"seo_url_for":              seoURL,
	}

	for key, value := range funcMaps {
		a.RegisterFuncMap(key, value)
	}
}
