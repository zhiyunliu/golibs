package datetime

import (
	"fmt"
	"time"

	"github.com/zhiyunliu/golibs/bytesconv"
)

//DateTime DateTime
type DateTime struct {
	opts *options
	sys  *time.Time
}

//NewDateTime 构建新的DateTime
func New(time time.Time, opts ...Option) *DateTime {
	opt := &options{
		Format: DefaultTimeformat,
	}
	val := &DateTime{opts: opt, sys: &time}

	for i := range opts {
		opts[i](opt)
	}

	return val
}

//MarshalJSON MarshalJSON
func (d *DateTime) MarshalJSON() (bytes []byte, err error) {
	val := d.sys
	tmpV := fmt.Sprintf(`"%s"`, val.Format(d.Format()))
	return bytesconv.StringToBytes(tmpV), nil
}

//UnmarshalJSON UnmarshalJSON
func (d *DateTime) UnmarshalJSON(bytes []byte) error {
	if d.opts == nil {
		d.opts = &options{
			Format: DefaultTimeformat,
		}
	}

	val, err := time.Parse(fmt.Sprintf(`"%s"`, d.opts.Format), bytesconv.BytesToString(bytes))
	*d = DateTime{sys: &val, opts: d.opts}
	return err
}

//Format 默认2006-01-02 15:04:05
func (d *DateTime) Format() string {
	return d.opts.Format
}

//String String
func (d *DateTime) String() string {
	return d.sys.Format(d.opts.Format)
}
