package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type RangeByTag struct {
	from int
	to   int
	tag  string
}

var data = []RangeByTag{
	RangeByTag{17001, 17065, "江苏"},
	RangeByTag{17066, 17125, "四川"},
	RangeByTag{17126, 17155, "浙江"},
	RangeByTag{17156, 17170, "上元"},
	RangeByTag{17171, 17190, "武汉"},
	RangeByTag{17191, 17200, "青岛"},
	RangeByTag{17201, 17206, "河南"},
	RangeByTag{17207, 17211, "江西"},
	RangeByTag{17212, 17219, "长沙"},
	RangeByTag{17220, 17225, "辽宁"},
	RangeByTag{17226, 17232, "广西"},
	RangeByTag{17233, 17240, "福建"},
	RangeByTag{17241, 17245, "广州"},
	RangeByTag{17246, 17249, "石家庄"},
	RangeByTag{17250, 17253, "京津"},
	RangeByTag{17254, 17256, "台元"},
	RangeByTag{17257, 17302, "梦果子"},
	RangeByTag{17303, 17374, "总部"}}

//var config = DBConfig{"560f33a30c4cb.sh.cdb.myqcloud.com:6322", "facewall-ganso", "root", "zealioniPLUS!"}

var config = DBConfig{"127.0.0.1:3306", "facewall-ganso", "root", "root"}
var ctrl Controller
var sc *ShuffleController

func init() {
	db, err := sql.Open("mysql", config.connectionString())
	if err != nil {
		log.Fatal(err)
	}

	ctrl = Controller{db}

	sc, err = NewShuffleController(nil, db)
}

func TestImport(t *testing.T) {
	expected := 0
	actual := 0
	for i := 0; i < len(data); i++ {
		from, to, tag := data[i].from, data[i].to, data[i].tag
		for j := from; j <= to; j++ {
			expected += 1
			var g = Guest{Code: strconv.Itoa(j), Tag: tag}
			k, err := ctrl.InsertGuest(g)
			if err == nil && k > 0 {
				actual += 1
			} else {
				fmt.Println(err)
			}
		}
	}

	if expected != actual {
		t.Errorf("expected: %d got: %d", expected, actual)
	}
}

func TestGetAssignment(t *testing.T) {
	var table = []struct {
		tag    string
		prize  string
		maxWin int
	}{
		{"浙江", "四等奖", 3},
		{"石家庄", "一等奖", 0},
		{"不存在", "一等奖", 0},
	}

	var max = 0
	for i := 0; i < len(table); i++ {
		max = ctrl.MaxWinFor(table[i].tag, table[i].prize)
		if max != table[i].maxWin {
			t.Errorf("expected: %d got: %d", table[i].maxWin, max)
		}
	}
}

func TestShuffle(t *testing.T) {
	var table = []struct {
		prize  string
		numWin int
	}{
		{"五等奖", 10},
		{"五等奖", 10},
		{"五等奖", 10},
		{"五等奖", 10},
		{"五等奖", 10},
		{"五等奖", 10},
		{"五等奖", 10},
		{"五等奖", 10},
		{"五等奖", 0},
		{"四等奖", 10},
		{"四等奖", 10},
		{"四等奖", 10},
		{"四等奖", 10},
		{"四等奖", 0},
		{"三等奖", 5},
		{"三等奖", 5},
		{"三等奖", 5},
		{"三等奖", 5},
		{"三等奖", 5},
		{"三等奖", 5},
		{"三等奖", 0},
		{"二等奖", 5},
		{"二等奖", 5},
		{"二等奖", 5},
		{"二等奖", 5},
		{"二等奖", 0},
		{"一等奖", 1},
		{"一等奖", 1},
		{"一等奖", 1},
		{"一等奖", 1},
		{"一等奖", 1},
		{"一等奖", 1},
		{"一等奖", 1},
		{"一等奖", 1},
		{"一等奖", 1},
		{"一等奖", 1},
		{"一等奖", 0},
	}
	for i := 0; i < len(table); i++ {
		step, winners, err := sc.stepWinnersForPrize(table[i].prize)
		if err != nil {
			t.Error(err)
		}
		if len(winners) != table[i].numWin {
			t.Errorf("step: %d, prize: %s, expected: %d, got: %d", step, table[i].prize, table[i].numWin, len(winners))
		}
	}
}

func TestReset(t *testing.T) {
	err := sc.reset()
	if err != nil {
		t.Error(err)
	}
}
