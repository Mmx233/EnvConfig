package EnvConfig

import (
	"os"
	"reflect"
	"strconv"
	"strings"
)

func parse(prefix string, c reflect.Value) {
	for i := 0; i < c.NumField(); i++ {
		v := c.Field(i)
		name := prefix + c.Type().Field(i).Name
		if v.Kind() == reflect.Struct {
			parse(name, v)
			continue
		}
		env := os.Getenv(name)
		if env == "" {
			if strings.Contains(c.Type().Field(i).Tag.Get("config"), "omitempty") {
				continue
			} else {
				panic("config v " + name + " not found")
			}
		}
		switch v.Kind() {
		case reflect.String:
			v.SetString(env)
		case reflect.Bool:
			v.SetBool(env == "true")
		case reflect.Uint:
			d, e := strconv.ParseUint(env, 10, 64)
			if e != nil {
				panic(e)
			}
			v.SetUint(d)
		case reflect.Int:
			d, e := strconv.ParseInt(env, 10, 64)
			if e != nil {
				panic(e)
			}
			v.SetInt(d)
		default:
			panic("config type " + v.Kind().String() + " not supported")
		}
	}
}

func Load(prefix string, config interface{}) {
	parse(prefix, reflect.ValueOf(config).Elem())
}
