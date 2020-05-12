package fa

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"package/locales"
	ut "package/universal-translator"
	"package/validator"
)

// RegisterDefaultTranslations registers a set of default translations
// for all built in tag's in validator; you may add your own as desired.
func RegisterDefaultTranslations(v *validator.Validate, trans ut.Translator) (err error) {

	translations := []struct {
		tag             string
		translation     string
		override        bool
		customRegisFunc validator.RegisterTranslationsFunc
		customTransFunc validator.TranslationFunc
	}{
		{
			tag:         "required",
			translation: "{0} نمیتونه خالی باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				fld, _ := ut.T(fe.Field())
				t, err := ut.T(fe.Tag(), fld)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag: "len",
			customRegisFunc: func(ut ut.Translator) (err error) {
				if err = ut.Add("len-string", "{0} باید {1} کاراکتر داشته باشه", false); err != nil {
					return
				}

				if err = ut.AddCardinal("len-string-character", "{0} ", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("len-string-character", "{0} ", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("len-number", "تعداد کاراکتر های {0} باید {1} باشه", false); err != nil {
					return
				}

				if err = ut.Add("len-items", "{0} باید شامل {1} آیتم باشه", false); err != nil {
					return
				}
				if err = ut.AddCardinal("len-items-item", "{0} ", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("len-items-item", "{0} ", locales.PluralRuleOther, false); err != nil {
					return
				}

				return
			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string

				var digits uint64
				var kind reflect.Kind

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err := strconv.ParseFloat(fe.Param(), 64)
				if err != nil {
					goto END
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					c, err = ut.C("len-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("len-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					c, err = ut.C("len-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("len-items", fe.Field(), c)

				default:
					t, err = ut.T("len-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "min",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("min-string", "طول {0} حداقل باید {1} کاراکتر باشه!", false); err != nil {
					return
				}

				if err = ut.AddCardinal("min-string-character", "{0}", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("min-string-character", "{0}", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("min-number", " {0} باید بزرگتر از {1} باشه!", false); err != nil {
					return
				}

				if err = ut.Add("min-items", "{0} باید شامل {1} باشه!", false); err != nil {
					return
				}
				if err = ut.AddCardinal("min-items-item", "{0} ", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("min-items-item", "{0} ", locales.PluralRuleOther, false); err != nil {
					return
				}

				return

			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				var err error
				var t string

				var digits uint64
				var kind reflect.Kind

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err := strconv.ParseFloat(fe.Param(), 64)
				if err != nil {
					goto END
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:
					var c string
					c, err = ut.C("min-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}
					t, err = ut.T("min-string", f, c)
				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					c, err = ut.C("min-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("min-items", f, c)

				default:
					t, err = ut.T("min-number", f, ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "max",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("max-string", "{0} حداکثر میتونه {1} کاراکتر داشته باشه!", false); err != nil {
					return
				}

				if err = ut.AddCardinal("max-string-character", "{0} ", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("max-string-character", "{0} ", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("max-number", "{0} باید {1} یا کمتر باشه!", false); err != nil {
					return
				}

				if err = ut.Add("max-items", "{0} حداکثر باید شامل {1} باشه!", false); err != nil {
					return
				}
				if err = ut.AddCardinal("max-items-item", "{0} آیتم", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("max-items-item", "{0} آیتم", locales.PluralRuleOther, false); err != nil {
					return
				}

				return

			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				var err error
				var t string

				var digits uint64
				var kind reflect.Kind

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err := strconv.ParseFloat(fe.Param(), 64)
				if err != nil {
					goto END
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					c, err = ut.C("max-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("max-string", f, c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					c, err = ut.C("max-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("max-items", f, c)

				default:
					t, err = ut.T("max-number", f, ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "eq",
			translation: "{0} برابر با {1} نیست!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f1, _ := ut.T(fe.Field())
				if f1 == "" {
					f1 = fe.Field()
				}
				f2, _ := ut.T(fe.Param())
				if f2 == "" {
					f2 = fe.Param()
				}
				t, err := ut.T(fe.Tag(), f1, f2)
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ne",
			translation: "{0} نباید برابر با{1} باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f1, _ := ut.T(fe.Field())
				if f1 == "" {
					f1 = fe.Field()
				}
				f2, _ := ut.T(fe.Param())
				if f2 == "" {
					f2 = fe.Param()
				}
				t, err := ut.T(fe.Tag(), f1, f2)
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "lt",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("lt-string", "{0} باید کمتر از {1} در مقدار کاراکتر باشه!", false); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-string-character", "{0} کاراکتر", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-string-character", "{0} کاراکتر", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lt-number", "{0} باید کمتر از {1} باشه!", false); err != nil {
					return
				}

				if err = ut.Add("lt-items", "{0} باید کمتر {1} باشه!", false); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-items-item", "{0} آیتم", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-items-item", "{0} آیتم", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lt-datetime", "{0} باید کمتر از روز و ساعت الان باشه!", false); err != nil {
					return
				}

				return

			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string
				var f64 float64
				var digits uint64
				var kind reflect.Kind

				fn := func() (err error) {

					if idx := strings.Index(fe.Param(), "."); idx != -1 {
						digits = uint64(len(fe.Param()[idx+1:]))
					}

					f64, err = strconv.ParseFloat(fe.Param(), 64)

					return
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("lt-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("lt-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("lt-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("lt-items", fe.Field(), c)

				case reflect.Struct:
					if fe.Type() != reflect.TypeOf(time.Time{}) {
						err = fmt.Errorf("tag '%s' cannot be used on a struct type", fe.Tag())
						goto END
					}

					t, err = ut.T("lt-datetime", fe.Field())

				default:
					err = fn()
					if err != nil {
						goto END
					}

					t, err = ut.T("lt-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "lte",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("lte-string", "{0} باید حداکثر {1} در طول کاراکتر باشه!", false); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-string-character", "{0} کاراکتر", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-string-character", "{0} کاراکتر", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lte-number", "{0} باید {1} یا کمتر باشه!", false); err != nil {
					return
				}

				if err = ut.Add("lte-items", "{0} نهایتا باید شامل {1} باشه!", false); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-items-item", "{0} آیتم", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-items-item", "{0} آیتم", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lte-datetime", "{0} باید کمتر یا برابر ساعت و تاریخ الان باشه!", false); err != nil {
					return
				}

				return
			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string
				var f64 float64
				var digits uint64
				var kind reflect.Kind

				fn := func() (err error) {

					if idx := strings.Index(fe.Param(), "."); idx != -1 {
						digits = uint64(len(fe.Param()[idx+1:]))
					}

					f64, err = strconv.ParseFloat(fe.Param(), 64)

					return
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("lte-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("lte-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("lte-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("lte-items", fe.Field(), c)

				case reflect.Struct:
					if fe.Type() != reflect.TypeOf(time.Time{}) {
						err = fmt.Errorf("tag '%s' cannot be used on a struct type", fe.Tag())
						goto END
					}

					t, err = ut.T("lte-datetime", fe.Field())

				default:
					err = fn()
					if err != nil {
						goto END
					}

					t, err = ut.T("lte-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "gt",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("gt-string", "{0} باید بیشتر از {1} در طول کاراکتر باشه!", false); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-string-character", "{0} کاراکتر", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-string-character", "{0} کاراکتر", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gt-number", "{0} باید بزرگتر از {1} باشه!", false); err != nil {
					return
				}

				if err = ut.Add("gt-items", "{0} باید شامل بیش از {1} باشه", false); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-items-item", "{0} آیتم", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-items-item", "{0} آیتم", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gt-datetime", "{0} باید بیشتر یا مساوی تاریخ و ساعت فعلی باشه!", false); err != nil {
					return
				}

				return
			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string
				var f64 float64
				var digits uint64
				var kind reflect.Kind

				fn := func() (err error) {

					if idx := strings.Index(fe.Param(), "."); idx != -1 {
						digits = uint64(len(fe.Param()[idx+1:]))
					}

					f64, err = strconv.ParseFloat(fe.Param(), 64)

					return
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("gt-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("gt-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("gt-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("gt-items", fe.Field(), c)

				case reflect.Struct:
					if fe.Type() != reflect.TypeOf(time.Time{}) {
						err = fmt.Errorf("tag '%s' cannot be used on a struct type", fe.Tag())
						goto END
					}

					t, err = ut.T("gt-datetime", fe.Field())

				default:
					err = fn()
					if err != nil {
						goto END
					}

					t, err = ut.T("gt-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "gte",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("gte-string", "{0} must be at least {1} in length", false); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-string-character", "{0} character", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-string-character", "{0} characters", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gte-number", "{0} must be {1} or greater", false); err != nil {
					return
				}

				if err = ut.Add("gte-items", "{0} must contain at least {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-items-item", "{0} item", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-items-item", "{0} items", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gte-datetime", "{0} must be greater than or equal to the current Date & Time", false); err != nil {
					return
				}

				return
			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string
				var f64 float64
				var digits uint64
				var kind reflect.Kind

				fn := func() (err error) {

					if idx := strings.Index(fe.Param(), "."); idx != -1 {
						digits = uint64(len(fe.Param()[idx+1:]))
					}

					f64, err = strconv.ParseFloat(fe.Param(), 64)

					return
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("gte-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("gte-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("gte-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("gte-items", fe.Field(), c)

				case reflect.Struct:
					if fe.Type() != reflect.TypeOf(time.Time{}) {
						err = fmt.Errorf("tag '%s' cannot be used on a struct type", fe.Tag())
						goto END
					}

					t, err = ut.T("gte-datetime", fe.Field())

				default:
					err = fn()
					if err != nil {
						goto END
					}

					t, err = ut.T("gte-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "eqfield",
			translation: "{0} باید با {1} یکی باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f1, _ := ut.T(fe.Field())
				if f1 == "" {
					f1 = fe.Field()
				}
				f2, _ := ut.T(fe.Param())
				if f2 == "" {
					f2 = fe.Param()
				}
				t, err := ut.T(fe.Tag(), f1, f2)
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "eqcsfield",
			translation: "{0} باید با {1} یکی باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f1, _ := ut.T(fe.Field())
				if f1 == "" {
					f1 = fe.Field()
				}
				f2, _ := ut.T(fe.Param())
				if f2 == "" {
					f2 = fe.Param()
				}
				t, err := ut.T(fe.Tag(), f1, f2)
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "necsfield",
			translation: "{0} نیمتونه با {1} یکی باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f1, _ := ut.T(fe.Field())
				if f1 == "" {
					f1 = fe.Field()
				}
				f2, _ := ut.T(fe.Param())
				if f2 == "" {
					f2 = fe.Param()
				}
				t, err := ut.T(fe.Tag(), f1, f2)
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtcsfield",
			translation: "{0} باید بزرگتر از {1} باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f1, _ := ut.T(fe.Field())
				if f1 == "" {
					f1 = fe.Field()
				}
				f2, _ := ut.T(fe.Param())
				if f2 == "" {
					f2 = fe.Param()
				}
				t, err := ut.T(fe.Tag(), f1, f2)
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtecsfield",
			translation: "{0} باید بزرگتر یا مساوری {1} باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f1, _ := ut.T(fe.Field())
				if f1 == "" {
					f1 = fe.Field()
				}
				f2, _ := ut.T(fe.Param())
				if f2 == "" {
					f2 = fe.Param()
				}
				t, err := ut.T(fe.Tag(), f1, f2)
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltcsfield",
			translation: "{0} باید کوچکتر از {1} باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f1, _ := ut.T(fe.Field())
				if f1 == "" {
					f1 = fe.Field()
				}
				f2, _ := ut.T(fe.Param())
				if f2 == "" {
					f2 = fe.Param()
				}
				t, err := ut.T(fe.Tag(), f1, f2)
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltecsfield",
			translation: "{0} باید کوچکتر یا مساوی {1} باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f1, _ := ut.T(fe.Field())
				if f1 == "" {
					f1 = fe.Field()
				}
				f2, _ := ut.T(fe.Param())
				if f2 == "" {
					f2 = fe.Param()
				}
				t, err := ut.T(fe.Tag(), f1, f2)
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "nefield",
			translation: "{0} نمی تونه برابر با {1} باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f1, _ := ut.T(fe.Field())
				if f1 == "" {
					f1 = fe.Field()
				}
				f2, _ := ut.T(fe.Param())
				if f2 == "" {
					f2 = fe.Param()
				}
				t, err := ut.T(fe.Tag(), f1, f2)
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtfield",
			translation: "{0} باید بزرگتر از {1} باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f1, _ := ut.T(fe.Field())
				if f1 == "" {
					f1 = fe.Field()
				}
				f2, _ := ut.T(fe.Param())
				if f2 == "" {
					f2 = fe.Param()
				}
				t, err := ut.T(fe.Tag(), f1, f2)
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtefield",
			translation: "{0} باید بزرگتر یا مساوی {1} باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f1, _ := ut.T(fe.Field())
				if f1 == "" {
					f1 = fe.Field()
				}
				f2, _ := ut.T(fe.Param())
				if f2 == "" {
					f2 = fe.Param()
				}
				t, err := ut.T(fe.Tag(), f1, f2)
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltfield",
			translation: "{0} باید کوچکتراز  {1} باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f1, _ := ut.T(fe.Field())
				if f1 == "" {
					f1 = fe.Field()
				}
				f2, _ := ut.T(fe.Param())
				if f2 == "" {
					f2 = fe.Param()
				}
				t, err := ut.T(fe.Tag(), f1, f2)
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltefield",
			translation: "{0} باید کوجکتر یا مساوی {1} باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f1, _ := ut.T(fe.Field())
				if f1 == "" {
					f1 = fe.Field()
				}
				f2, _ := ut.T(fe.Param())
				if f2 == "" {
					f2 = fe.Param()
				}
				t, err := ut.T(fe.Tag(), f1, f2)
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "alpha",
			translation: "{0} فقط میتونه شامل حروف باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "alphanum",
			translation: "{0} فقط میتونه شامل حروف و عدد باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "numeric",
			translation: "{0} باید یک مقدار عددی معتبر باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "number",
			translation: "{0} باید یک عدد معتبر باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "hexadecimal",
			translation: "{0} باید یک مقدار هگزا دسیمال معتبر باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "hexcolor",
			translation: "{0} باید یک هگز رنگ معتبر باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "rgb",
			translation: "{0} باید یک مقدار RGB معتبر باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "rgba",
			translation: "{0} باید یک مقدار RGBA معتبر باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "hsl",
			translation: "{0} باید یک مقدار رنگ HSL معتبر باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "hsla",
			translation: "{0} باید یک مقدار رنگ HSLA معتبر باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "email",
			translation: "{0} باید یک مقدار معتبر باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "url",
			translation: "{0} باید یک مقدار معتبر باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "uri",
			translation: "{0} باید یک مقدار معتبر باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "base64",
			translation: "{0} باید یک رشته Base64 معتبر باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "contains",
			translation: "{0} باید شامل کلمه '{1}' باشد",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f1, _ := ut.T(fe.Field())
				if f1 == "" {
					f1 = fe.Field()
				}
				f2, _ := ut.T(fe.Param())
				if f2 == "" {
					f2 = fe.Param()
				}
				t, err := ut.T(fe.Tag(), f1, f2)
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "containsany",
			translation: "{0} حداقل باید شامل یکی از '{1}' باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f1, _ := ut.T(fe.Field())
				if f1 == "" {
					f1 = fe.Field()
				}
				f2, _ := ut.T(fe.Param())
				if f2 == "" {
					f2 = fe.Param()
				}
				t, err := ut.T(fe.Tag(), f1, f2)
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "excludes",
			translation: "{0} نمیتونه شامل متن '{1}' باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f1, _ := ut.T(fe.Field())
				if f1 == "" {
					f1 = fe.Field()
				}
				f2, _ := ut.T(fe.Param())
				if f2 == "" {
					f2 = fe.Param()
				}
				t, err := ut.T(fe.Tag(), f1, f2)
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "excludesall",
			translation: "{0} نمیتونه شامل کاراکترهای '{1}' باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f1, _ := ut.T(fe.Field())
				if f1 == "" {
					f1 = fe.Field()
				}
				f2, _ := ut.T(fe.Param())
				if f2 == "" {
					f2 = fe.Param()
				}
				t, err := ut.T(fe.Tag(), f1, f2)
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "excludesrune",
			translation: "{0} نمیتونه شامل '{1}' باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f1, _ := ut.T(fe.Field())
				if f1 == "" {
					f1 = fe.Field()
				}
				f2, _ := ut.T(fe.Param())
				if f2 == "" {
					f2 = fe.Param()
				}
				t, err := ut.T(fe.Tag(), f1, f2)
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "isbn",
			translation: "{0} باید یک شماره ISBN معتبر باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "isbn10",
			translation: "{0} باید یک ISBN-10 باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "isbn13",
			translation: "{0} باید یک شماره ISBN-13 معتبر باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "uuid",
			translation: "{0} باید یک UUID معتبر باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "uuid3",
			translation: "{0} باید یک UUID ورژن 3 معتبر باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "uuid4",
			translation: "{0} باید یک مقدار UUID ورژن 4 معتبر باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "uuid5",
			translation: "{0} باید یک مقدار UUID ورژن 5 معتبر باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "ascii",
			translation: "{0} فقط میتونه شامل کارکترهای ascii معتبر باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "printascii",
			translation: "{0} فقط میتونه شامل کدهای ascii قابل چاپ باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "multibyte",
			translation: "{0} فقط میتونه شامل کارکترهای مولتی بایت باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "datauri",
			translation: "{0} فقط میتونه شامل URI باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "latitude",
			translation: "{0} فقط میتونه شامل مختصات latitude باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "longitude",
			translation: "{0} فقط میتونه شامل مختصات longitude باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "ssn",
			translation: "{0} باید یک شماره SSN معتبر باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "ipv4",
			translation: "{0} باید یک آدرس IPv4 معتبر باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "ipv6",
			translation: "{0} باید یک آدرس IPv6 معتبر باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "ip",
			translation: "{0} باید یک آدرس IP معتبر باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "cidr",
			translation: "{0} باید شامل CIDR notation باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "cidrv4",
			translation: "{0} باید شامل CIDR notation معتبر برای آدرس IPv4 باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "cidrv6",
			translation: "{0} باید شامل CIDR notation برای آدرس IPv6 باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "tcp_addr",
			translation: "{0} باید یک آدرس TCP معتبر باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "tcp4_addr",
			translation: "{0} باید یک آدرس IPv4 TCP باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "tcp6_addr",
			translation: "{0} باید یک شامل آدرس معتبر IPv6 TCP باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "udp_addr",
			translation: "{0} باید شامل آدرس UDP معتبر باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "udp4_addr",
			translation: "{0} باید شامل آدرس  IPv4 UDP معتبر باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "udp6_addr",
			translation: "{0} باید شامل آدرس معتبر IPv6 UDP باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "ip_addr",
			translation: "{0} باید یک IP قابل دسترس باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "ip4_addr",
			translation: "{0} باید IPv4 address معتبر باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "ip6_addr",
			translation: "{0} باید IPv6 address معتبر باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "unix_addr",
			translation: "{0} باید UNIX address باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "mac",
			translation: "{0} باید MAC address معتبر باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "iscolor",
			translation: "{0} باید یک رنگ معتبر باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "oneof",
			translation: "{0} باید یکی از [{1}] باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				s, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}
				return s
			},
		},
		{
			tag:         "username",
			translation: "{0} فقط میتونه شامل حروف انگلیسی و _ باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "localchars",
			translation: "{0} فقط میتونه شامل حروف فارسی و اعداد باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "localphone",
			translation: "شماره تلفن همراه باید مقدار معتبر باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "boolean",
			translation: "باید یک مقدار بولین باشه",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "charsonly",
			translation: "{0} فقط میتونه شامل حروف باشه!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
		{
			tag:         "unique",
			translation: "{0} تکراریه و قبلا انتخاب شده!",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				f, _ := ut.T(fe.Field())
				if f == "" {
					f = fe.Field()
				}
				t, err := ut.T(fe.Tag(), f)
				if err != nil {
					return fe.(error).Error()
				}
				return t
			},
		},
	}

	for _, t := range translations {

		if t.customTransFunc != nil && t.customRegisFunc != nil {

			err = v.RegisterTranslation(t.tag, trans, t.customRegisFunc, t.customTransFunc)

		} else if t.customTransFunc != nil && t.customRegisFunc == nil {

			err = v.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), t.customTransFunc)

		} else if t.customTransFunc == nil && t.customRegisFunc != nil {

			err = v.RegisterTranslation(t.tag, trans, t.customRegisFunc, translateFunc)

		} else {
			err = v.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), translateFunc)
		}

		if err != nil {
			return
		}
	}

	return
}

func registrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {

	return func(ut ut.Translator) (err error) {

		if err = ut.Add(tag, translation, override); err != nil {
			return
		}

		return

	}

}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {

	t, err := ut.T(fe.Tag(), fe.Field())
	if err != nil {
		log.Printf("warning: error translating FieldError: %#v", fe)
		return fe.(error).Error()
	}

	return t
}
