package app

// Create 创建单条或者多条记录
//
// Example
//
//		u := User{Name:"abc"}
//		db.DB(ctx).Create(&u)
//
//	 us := []*User{
//	   {Name: "Jinzhu", Age: 18, Birthday: time.Now()},
//	   {Name: "Jackson", Age: 19, Birthday: time.Now()},
//	 }
//	 db.DB(ctx).Create(us)
//
func (d *Database) Create(v any) Result {
	res := d.db().Create(v)
	res.Error = wrapCaller(res.Error, 1)
	return Result{
		RowsAffected: res.RowsAffected,
		Error:        res.Error,
	}
}

func (d *Database) FirstOrCreate(dest interface{}, conds ...interface{}) Result {
	res := d.db().FirstOrCreate(dest, conds...)
	res.Error = wrapCaller(res.Error, 1)
	return Result{
		RowsAffected: res.RowsAffected,
		Error:        res.Error,
	}
}
