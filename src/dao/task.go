package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"img-crawler/src/log"
	"img-crawler/src/utils"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	sq "gopkg.in/Masterminds/squirrel.v1"
)

type Task struct {
	ID          uint64         `db:"id,PRIMARY_KEY,AUTO_INCREMENT"`
	Name        string         `db:"name"`
	Seeds       string         `db:"seeds"`
	Desc        sql.NullString `db:"desci"`
	Status      int            `db:"status"`
	CreatedTime time.Time      `db:"create_time"`
	UpdatedTime time.Time      `db:"modify_time"`

	items []*TaskItem
}

type TaskItem struct {
	ID          uint64         `db:"id,PRIMARY_KEY,AUTO_INCREMENT"`
	TaskID      uint64         `db:"task_id"`
	Name        string         `db:"name"`
	Desc        sql.NullString `db:"desci"`
	Url         string         `db:"url"`
	FilePath    sql.NullString `db:"filepath"`
	Digest      sql.NullString `db:"digest"`
	Status      int            `db:"status"`
	Effective   int            `db:"effective"`
	CreatedTime time.Time      `db:"create_time"`
	UpdatedTime time.Time      `db:"modify_time"`
}

type TaskDAO interface {
	CreateTask(*Task) (uint64, error)
	CreateTaskItem(item *TaskItem, taskID uint64) (uint64, error)
	CreateItemTable(uint64)
	Get(map[string]interface{}) (*Task, error)
	ListItems(status, num uint64) ([]*TaskItem, error)
	List(bool, map[string]interface{}) ([]*Task, error)
	Update(bool, map[string]interface{}, map[string]interface{}) (int64, error)
}

type TaskDAOImpl struct {
	pool *Pool
	tb   string // table name
	tb_n string // nested table name
	Tb_r string
}

var _ TaskDAO = (*TaskDAOImpl)(nil)

func NewTaskDAO(pool *Pool) *TaskDAOImpl {
	return &TaskDAOImpl{
		pool: pool,
		tb:   "tasks",
		tb_n: "task_items",
		Tb_r: "task_items_"}
}

func (dao *TaskDAOImpl) CreateItemTable(id uint64) {
	db := sqlx.NewDb(dao.pool.Master().GetDB(), "mysql")
	dao.Tb_r += strconv.FormatUint(id, 10)
	schema := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s LIKE %s", dao.Tb_r, dao.tb_n)
	db.MustExec(schema)
}

func (dao *TaskDAOImpl) Get(conditions map[string]interface{}) (*Task, error) {

	db := sqlx.NewDb(dao.pool.Slave().GetDB(), "mysql")

	sql, args, err := sq.Select("*").From(dao.tb).Where(sq.Eq(conditions)).ToSql()

	log.Infof(sql, args...)

	task := Task{}
	err = db.Get(&task, sql, args...)

	if err != nil {
		log.Warnf("Task %s Get Failed %s", conditions, err)
		return nil, err
	}

	return &task, nil
}

func (dao *TaskDAOImpl) ListItems(status, num uint64) ([]*TaskItem, error) {
	db := sqlx.NewDb(dao.pool.Slave().GetDB(), "mysql")

	sql, args, err := sq.Select("*").From(dao.Tb_r).Where(sq.Eq{"status": status}).Limit(num).ToSql()

	log.Infof(sql, args...)

	items := []*TaskItem{}
	err = db.Select(&items, sql, args...)

	if err != nil {
		log.Warnf("Task Item Has no Records {%s}", err)
		return make([]*TaskItem, 0), err
	}

	return items, nil
}

func (dao *TaskDAOImpl) CreateTask(task *Task) (id uint64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			log.Errorf("DAO CreateTask failed: %s", err)
		}
	}()

	db := sqlx.NewDb(dao.pool.Master().GetDB(), "mysql")

	clauses := GetMapping(*task)
	sql, args, err := sq.Insert(dao.tb).SetMap(clauses).ToSql()
	utils.CheckError(err)

	log.Infof(sql, args...)

	res := db.MustExec(sql, args...)

	id2, err := res.LastInsertId()
	utils.CheckError(err)

	return uint64(id2), err
}

func (dao *TaskDAOImpl) CreateTaskItem(item *TaskItem, taskID uint64) (id uint64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			log.Errorf("DAO CreateTaskItem failed: %s", err)
		}
	}()

	db := sqlx.NewDb(dao.pool.Master().GetDB(), "mysql")

	clauses := GetMapping(*item)
	clauses["task_id"] = taskID

	sql, args, err := sq.Insert(dao.Tb_r).SetMap(clauses).ToSql()
	utils.CheckError(err)

	log.Infof(sql, args...)

	res := db.MustExec(sql, args...)

	id2, err := res.LastInsertId()
	utils.CheckError(err)

	return uint64(id2), err
}

/* Dump all fileds, Rewrite one row */
func (dao *TaskDAOImpl) Update(items bool, conditions, clauses map[string]interface{}) (int64, error) {
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			log.Errorf("DAO Update failed: %s", err)
		}
	}()

	if len(clauses) == 0 || len(conditions) == 0 {
		return 0, errors.New("UpdateTask Arguments Error")
	}

	db := sqlx.NewDb(dao.pool.Master().GetDB(), "mysql")

	var ctx = context.Background()
	tx := db.MustBeginTx(ctx, nil)

	tb_name := dao.tb
	if items {
		tb_name = dao.Tb_r
	}

	sql, args, err := sq.Update(tb_name).SetMap(clauses).Where(conditions).ToSql()
	utils.CheckError(err)

	log.Infof(sql, args...)

	res := tx.MustExec(sql, args...)

	num, _ := res.RowsAffected()

	log.Infof("Update Task with %d rows affected", num)

	err = tx.Commit()
	utils.CheckError(err)

	return num, nil
}

func (dao *TaskDAOImpl) List(items bool, conditions map[string]interface{}) ([]*Task, error) {
	db := sqlx.NewDb(dao.pool.Slave().GetDB(), "mysql")

	tb_name := dao.tb
	if items {
		tb_name = dao.Tb_r
	}

	sql, args, err := sq.Select("*").From(tb_name).Where(sq.Eq(conditions)).ToSql()

	log.Infof(sql, args...)

	tasks := []*Task{}
	err = db.Select(&tasks, sql, args...)

	if err != nil {
		log.Warnf("no task found {%s}", err)
		return make([]*Task, 0), err
	}

	return tasks, nil
}
