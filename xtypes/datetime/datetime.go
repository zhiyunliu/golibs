package datetime

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/zhiyunliu/golibs/bytesconv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	val, err := time.ParseInLocation(fmt.Sprintf(`"%s"`, d.opts.Format), bytesconv.BytesToString(bytes), time.Local)
	if err != nil {
		return err
	}
	*d = DateTime{Time: val, opts: d.opts}
	return nil
}

// Format 默认2006-01-02 15:04:05
func (d DateTime) Format() string {
	if d.opts == nil {
		return DefaultTimeformat
	}
	return d.opts.Format
}

// String String
func (d DateTime) String() string {
	return d.Time.Format(d.Format())
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
		tmp := New(transferTolocal(vt))
		*t = *tmp
	case *time.Time:
		// 字符串转成 time.Time 类型
		// 切换时区
		tmp := New(transferTolocal(*vt))
		*t = *tmp
	case string:
		tmpDate, err := time.ParseInLocation(DefaultTimeformat, vt, time.Local)
		if err != nil {
			return err
		}
		tmp := New(tmpDate)
		*t = *tmp
	case *string:
		tmpDate, err := time.ParseInLocation(DefaultTimeformat, *vt, time.Local)
		if err != nil {
			return err
		}
		tmp := New(tmpDate)
		*t = *tmp
	default:
		return fmt.Errorf("类型处理错误:%+v", v)
	}
	return nil
}

// UTC时间转换为北京时间
func TransferUtcToCts8(t time.Time) time.Time {
	// 解析数据库时间相关的字段没有时区，默认转换成了UTC时间
	cstTime := t.In(time.Local)
	cstTime = cstTime.Add(-time.Hour * 8)
	return cstTime
}

func transferTolocal(t time.Time) time.Time {
	timeStr := t.Format("2006-01-02 15:04:05")
	t1, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
	return t1
}

func (t DateTime) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(time.Time(t.Time))
}

func (t *DateTime) UnmarshalBSONValue(bt bsontype.Type, data []byte) error {
	if bt != bson.TypeDateTime {
		return errors.New("DateTime UnmarshalBSONValue type error, want DateTime but get " + bt.String())
	}
	var tm primitive.DateTime
	err := bson.UnmarshalValue(bt, data, &tm)
	if err != nil {
		return err
	}
	*t = *New(tm.Time())
	return nil
}
