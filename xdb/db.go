package xdb

import (
	"time"

	"github.com/zhiyunliu/xdb/tpl"
)

//IDB 数据库操作接口,安装可需能需要执行export LD_LIBRARY_PATH=/usr/local/lib
type IDB interface {
	Executer
	Begin() (ITrans, error)
	Close()
}

//ITrans 数据库事务接口
type ITrans interface {
	Executer
	Rollback() error
	Commit() error
}

//Executer 数据库操作对象集合
type Executer interface {
	Query(sql string, input map[string]interface{}) (data Rows, err error)
	First(sql string, input map[string]interface{}) (data Row, err error)
	Scalar(sql string, input map[string]interface{}) (data interface{}, err error)
	Exec(sql string, input map[string]interface{}) (r Result, err error)
	ExecSp(procName string, input map[string]interface{}) (r Result, err error)
}

//DB 数据库操作类
type DB struct {
	db  ISysDB
	tpl tpl.ITPLContext
}

//NewDB 创建DB实例
func NewDB(proto string, connString string, maxOpen int, maxIdle int, maxLifeTime int) (obj IDB, err error) {
	dbobj := &DB{}
	dbobj.tpl, err = tpl.GetDBContext(proto)
	if err != nil {
		return
	}
	dbobj.db, err = NewSysDB(proto, connString, maxOpen, maxIdle, time.Duration(maxLifeTime)*time.Second)
	obj = dbobj
	return
}

//GetTPL 获取模板翻译参数
func (db *DB) GetTPL() tpl.ITPLContext {
	return db.tpl
}

//Query 查询数据
func (db *DB) Query(sql string, input map[string]interface{}) (data Rows, err error) {
	query, args := db.tpl.GetSQLContext(sql, input)
	data, err = db.db.Query(query, args...)
	if err != nil {
		return nil, getDBError(err, query, args)
	}
	return
}

//Scalar 根据包含@名称占位符的查询语句执行查询语句
func (db *DB) Scalar(sql string, input map[string]interface{}) (data interface{}, err error) {
	query, args := db.tpl.GetSQLContext(sql, input)
	result, err := db.db.Query(query, args...)
	if err != nil {
		return nil, getDBError(err, query, args)
	}
	if result.Len() == 0 || len(result[0]) == 0 {
		return nil, nil
	}
	data, _ = result[0].Get(result[0].Keys()[0])
	return
}

//Execute 根据包含@名称占位符的语句执行查询语句
func (db *DB) Execute(sql string, input map[string]interface{}) (r Result, err error) {
	query, args := db.tpl.GetSQLContext(sql, input)
	r, err = db.db.Execute(query, args...)
	if err != nil {
		return nil, getDBError(err, query, args)
	}
	return
}

//ExecuteSP 根据包含@名称占位符的语句执行查询语句
func (db *DB) ExecSp(procName string, input map[string]interface{}, output ...interface{}) (r Result, err error) {
	query, args := db.tpl.GetSPContext(procName, input)
	ni := append(args, output...)
	r, err = db.db.Execute(query, ni...)
	if err != nil {
		return nil, getDBError(err, query, ni)
	}
	return
}

//Replace 替换SQL语句中的参数
func (db *DB) Replace(sql string, args []interface{}) string {
	return db.tpl.Replace(sql, args)
}

//Begin 创建事务
func (db *DB) Begin() (t ITrans, err error) {
	tt := &DBTrans{}
	tt.tx, err = db.db.Begin()
	if err != nil {
		return
	}
	tt.tpl = db.tpl
	return tt, nil
}

//Close  关闭当前数据库连接
func (db *DB) Close() {
	db.db.Close()
}
