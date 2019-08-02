package data

import (
	"testing"
)

func TestTransaction(t *testing.T) {
	source := GetSource()
	src, ok := source.(*Source)
	if !ok {
		t.Error("Source is not implements SourceI")
		return
	}
	_, err := src.DB.Exec("CREATE TABLE IF NOT EXISTS test (a VARCHAR(20))")
	if err != nil {
		t.Error(err)
		return
	}
	tx, err := src.DB.Begin()
	if err != nil {
		t.Error(err)
		return
	}
	stmt, err := tx.Prepare("insert into test(a)values (?)")
	if err != nil {
		t.Error(err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec("1")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = tx.Exec("insert into test(a)values ('2')")
	if err != nil {
		t.Error(err)
		return
	}
	err = tx.Rollback()
	if err != nil {
		t.Error(err)
		return
	}
	rows, err := src.DB.Query("select a from test")
	if err != nil {
		t.Error(err)
		return
	}
	defer rows.Close()
	if rows.Next() {
		t.Error("rollback fail.")
	} else {
		t.Log("rollback success.")
	}
}
