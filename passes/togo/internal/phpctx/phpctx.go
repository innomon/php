package phpctx

import (
	"errors"
	"flag"
	"io"
	"net/http"
	"os/exec"
	"reflect"
	"strings"
)

var (
	ErrNoStruct     = errors.New("receiver is not a struct")
	ErrMissingField = errors.New("struct is missing field")
	ErrNotSet       = errors.New("dynamic value is not set")
)

var zero = reflect.Value{}
var addr = flag.String("http-address", "localhost:8080", "host and port for serving http")

type PHPVar interface{}

type PHPContext struct {
	Echo          io.Writer
	dynamicValues map[string]*interface{}
}

func (ctx PHPContext) SetDynamic(name string, value interface{}) {
	ctx.dynamicValues[name] = &value
}

func (ctx PHPContext) GetDynamic(name string) (interface{}, error) {
	v, ok := ctx.dynamicValues[name]
	if !ok {
		return nil, ErrNotSet
	}

	if v == nil {
		return nil, nil
	}

	return *v, nil
}

func (ctx PHPContext) Write(b []byte) (int, error) {
	return ctx.Echo.Write(b)
}

func ListenAndServe() {
	if !flag.Parsed() {
		flag.Parse()
	}
	http.ListenAndServe(*addr, nil)
}

func GetDynamicProperty(rcvr interface{}, field string) (interface{}, error) {
	v := reflect.ValueOf(rcvr)
	if v.Kind() == reflect.Struct {
		f := v.FieldByName(field)
		if f == zero {
			return nil, ErrMissingField
		}
		return f.Interface(), nil
	}
	return nil, ErrNoStruct
}

func Shell(cmd string) ([]byte, error) {
	cmdParts := strings.SplitN(cmd, " ", 2)
	cmdName := cmdParts[0]
	args := cmdParts[1]
	c := exec.Command(cmdName, args)
	err := c.Run()
	if err != nil {
		return nil, err
	}

	return c.Output()
}

type Server interface {
	ServePHP(ctx PHPContext)
}
