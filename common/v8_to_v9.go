package common

import (
	"reflect"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"
)

var (
	validatorManager = make(map[string]*validator.Validate)
	managerMutex     = &sync.Mutex{}
	transManager     = make(map[string]ut.Translator)
)

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ binding.StructValidator = &defaultValidator{}

// NewValidator 返回 defaultValidator 实例
func NewValidator() *defaultValidator {
	return &defaultValidator{}
}

// SingletonValidator 返回 tagName 为 tagName 变量的单例 Validate
func SingletonValidator(tagName string) (*validator.Validate, error) {
	if v, ok := validatorManager[tagName]; ok {
		return v, nil
	}
	managerMutex.Lock()
	defer managerMutex.Unlock()
	// 中文 trans
	v := validator.New()
	v.SetTagName(tagName)
	t := newZHTrans()
	if err := zh_translations.RegisterDefaultTranslations(v, t); err != nil {
		return nil, err
	}

	transManager[tagName] = t
	validatorManager[tagName] = v
	return v, nil
}

// GetZHTrans 返回 中文 translator
func GetZHTrans(tagName string) ut.Translator {
	return transManager[tagName]
}

func (v *defaultValidator) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {

		v.lazyinit()

		if err := v.validate.Struct(obj); err != nil {
			return error(err)
		}
	}

	return nil
}

func (v *defaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate, _ = SingletonValidator("binding")

		// add any custom validations etc. here
	})
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

func newZHTrans() ut.Translator {
	zh := zh.New()
	uni := ut.New(zh, zh)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	t, _ := uni.GetTranslator("zh")
	return t
}
