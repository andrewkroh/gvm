package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/pkg/errors"
)

type envFormatter interface {
	Set(name, val string) string
	Prepend(name, val string) string
	Append(name, val string) string
}

type (
	bashFormatter       struct{}
	batchFormatter      struct{}
	powershellFormatter struct{}
)

var (
	_batchFormatter      envFormatter = (*batchFormatter)(nil)
	_bashFormatter       envFormatter = (*bashFormatter)(nil)
	_powershellFormatter envFormatter = (*powershellFormatter)(nil)
)

// Output formats.
const (
	BashFormat       = "bash"
	BatchFormat      = "batch"
	PowershellFormat = "powershell"
)

func defaultFormat() string {
	if runtime.GOOS == "windows" {
		return BatchFormat
	}
	return BashFormat
}

func getEnvFormatter(format string) (envFormatter, error) {
	switch format {
	case BashFormat:
		return _bashFormatter, nil
	case BatchFormat:
		return _batchFormatter, nil
	case PowershellFormat:
		return _powershellFormatter, nil
	default:
		return nil, errors.Errorf("invalid format option '%v'", format)
	}
}

func (f *bashFormatter) Set(name, val string) string {
	return fmt.Sprintf(`export %v="%v"`, name, val)
}

func (f *bashFormatter) Prepend(name, val string) string {
	return fmt.Sprintf(`export %v="%v;$%v"`, name, val, name)
}
func (f *bashFormatter) Append(name, val string) string {
	return fmt.Sprintf(`export %v="$%v;%v"`, name, name, val)
}

func (f *batchFormatter) Set(name, val string) string {
	return fmt.Sprintf(`set %v=%v`, name, val)
}

func (f *batchFormatter) Prepend(name, val string) string {
	return fmt.Sprintf(`set %v=%v;$%v`, name, val, os.Getenv(name))
}
func (f *batchFormatter) Append(name, val string) string {
	return fmt.Sprintf(`set %v=$%v;%v`, name, os.Getenv(name), val)
}

func (f *powershellFormatter) Set(name, val string) string {
	return fmt.Sprintf(`$env:%v = "%v"`, name, val)
}

func (f *powershellFormatter) Prepend(name, val string) string {
	return fmt.Sprintf(`$env:%v = "%v;$env:%v"`, name, val, name)
}
func (f *powershellFormatter) Append(name, val string) string {
	return fmt.Sprintf(`$env:%v="$env:%v;%v"`, name, name, val)
}
