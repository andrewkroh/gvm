package shellfmt

import (
	"fmt"
	"io"
	"os"
	"runtime"
)

type Fmt struct {
	out io.Writer
	fmt EnvFormatter
}

type EnvFormatter interface {
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
	_batchFormatter      EnvFormatter = (*batchFormatter)(nil)
	_bashFormatter       EnvFormatter = (*bashFormatter)(nil)
	_powershellFormatter EnvFormatter = (*powershellFormatter)(nil)
)

// Output formats.
const (
	BashFormat       = "bash"
	BatchFormat      = "batch"
	PowershellFormat = "powershell"
)

func New(format string) (*Fmt, error) {
	formatter, err := GetEnvFormatter(format)
	if err != nil {
		return nil, err
	}
	return &Fmt{out: os.Stdout, fmt: formatter}, nil
}

func (f *Fmt) Set(name, val string) {
	fmt.Println(f.fmt.Set(name, val))
}

func (f *Fmt) Prepend(name, val string) {
	fmt.Println(f.fmt.Prepend(name, val))
}

func (f *Fmt) Append(name, val string) {
	fmt.Println(f.fmt.Append(name, val))
}

func DefaultFormat() string {
	if runtime.GOOS == "windows" {
		return BatchFormat
	}
	return BashFormat
}

func GetEnvFormatter(format string) (EnvFormatter, error) {
	if format == "" {
		format = DefaultFormat()
	}

	switch format {
	case BashFormat:
		return _bashFormatter, nil
	case BatchFormat:
		return _batchFormatter, nil
	case PowershellFormat:
		return _powershellFormatter, nil
	default:

		return nil, fmt.Errorf("invalid format option: %q", format)
	}
}

func (f *bashFormatter) Set(name, val string) string {
	return fmt.Sprintf(`export %v="%v"`, name, val)
}

func (f *bashFormatter) Prepend(name, val string) string {
	return fmt.Sprintf(`export %v="%v:$%v"`, name, val, name)
}

func (f *bashFormatter) Append(name, val string) string {
	return fmt.Sprintf(`export %v="$%v:%v"`, name, name, val)
}

func (f *batchFormatter) Set(name, val string) string {
	return fmt.Sprintf(`set %v=%v`, name, val)
}

func (f *batchFormatter) Prepend(name, val string) string {
	return fmt.Sprintf(`set %v=%v;%v`, name, val, os.Getenv(name))
}

func (f *batchFormatter) Append(name, val string) string {
	return fmt.Sprintf(`set %v=%v;%v`, name, os.Getenv(name), val)
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
