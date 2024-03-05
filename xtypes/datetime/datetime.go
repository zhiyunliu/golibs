package datetime

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/zhiyunliu/golibs/bytesconv"
)

// DateTime DateTime
type DateTime struct {
	opts *options
	time.Time
}

// NewDateTime 构建新的DateTime
func New(time time.Time, opts ...Option) *DateTime {
	opt := &options{
		Format: DefaultTimeformat,
	}
	val := &DateTime{opts: opt, Time: time}

	for i := range opts {
		opts[i](opt)
	}

	return val
}

// MarshalJSON MarshalJSON
func (d DateTime) MarshalJSON() (bytes []byte, err error) {
	tmpV := fmt.Sprintf(`"%s"`, d.Time.Format(d.Format()))
	return bytesconv.StringToBytes(tmpV), nil
}

// UnmarshalJSON UnmarshalJSON
func (d *DateTime) UnmarshalJSON(bytes []byte) error {
	if d.opts == nil {
		d.opts = &options{
			Format: DefaultTimeformat,
		}
	}

	val, err := time.Parse(fmt.Sprintf(`"%s"`, d.opts.Format), bytesconv.BytesToString(bytes))
	if err != nil {
		return err
	}
	val = val.Local()
	*d = DateTime{Time: val, opts: d.opts}
	return nil
}

// Format 默认2006-01-02 15:04:05
func (d DateTime) Format() string {
	return d.opts.Format
}

// String String
func (d DateTime) String() string {
	return d.Time.Format(d.opts.Format)
}

// Value String
func (d DateTime) Value() (driver.Value, error) {
	return d.String(), nil
}

func (t *DateTime) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		// 字符串转成 time.Time 类型
		// 切换时区
		tmp := New(vt.Local())
		*t = *tmp
	case string:
		tmpDate, err := time.Parse(DefaultTimeformat, vt)
		if err != nil {
			return err
		}
		tmp := New(tmpDate.Local())
		*t = *tmp
	default:
		return fmt.Errorf("类型处理错误:%+v", v)
	}
	return nil
}
