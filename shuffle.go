package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Prize struct {
	Prize      string
	MaxStepWin int
}

type Sum struct {
	Tag    string
	Prize  string
	MaxWin int
	SumWin int
}

/*
  shuffle states:
  - all
  - shuffling/{prize}
  - history/{step}
*/
type ShuffleController struct {
	state              string
	db                 *sql.DB
	router             *mux.Router
	maxStepWinForPrize map[string]int
}

func NewShuffleController(r *mux.Router, db *sql.DB) (*ShuffleController, error) {
	sc := new(ShuffleController)

	sc.state = "all"
	sc.db = db

	prizes, err := sc.getPrizes()
	if err != nil {
		return nil, err
	}

	sc.maxStepWinForPrize = make(map[string]int)
	for i := 0; i < len(prizes); i++ {
		sc.maxStepWinForPrize[prizes[i].Prize] = prizes[i].MaxStepWin
	}

	if r != nil {
		sc.router = r
		sc.router.StrictSlash(true)
		sc.router.HandleFunc("/start/{prize}", sc.start).Methods("POST")
		sc.router.HandleFunc("/end", sc.end).Methods("POST")
		sc.router.HandleFunc("/history/{step}", sc.history).Methods("GET")
		sc.router.HandleFunc("/all", sc.all).Methods("GET")
		sc.router.HandleFunc("/state", sc.getState).Methods("GET")
	}
	return sc, nil
}

func (sc *ShuffleController) start(w http.ResponseWriter, r *http.Request) {
	sc.state = "shuffling/一等奖"
	w.Write([]byte(sc.state))
}

func (sc *ShuffleController) end(w http.ResponseWriter, r *http.Request) {
	sc.state = "history/2"
	w.Write([]byte(sc.state))
}

func (sc *ShuffleController) history(w http.ResponseWriter, r *http.Request) {
	sc.state = "history/1"
	w.Write([]byte(sc.state))
}

func (sc *ShuffleController) all(w http.ResponseWriter, r *http.Request) {
	sc.state = "all"
	w.Write([]byte(sc.state))
}

func (sc *ShuffleController) getState(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(sc.state))
}

func (sc *ShuffleController) getPrizes() ([]Prize, error) {
	rows, err := sc.db.Query("SELECT prize, maxStepWin FROM prize")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var prizes []Prize
	for rows.Next() {
		var prize Prize
		err = rows.Scan(&prize.Prize, &prize.MaxStepWin)
		if err != nil {
			return nil, err
		}
		prizes = append(prizes, prize)
	}

	return prizes, nil
}

// select winners for prize one step, and save to database
func (sc *ShuffleController) stepWinnersForPrize(prize string) ([]int, error) {
	maxStepWin, found := sc.maxStepWinForPrize[prize]
	if !found {
		return nil, errors.New("missing prize config for " + prize)
	}

	// this step
	var nullStep sql.NullInt64
	var step int64
	err := sc.db.QueryRow("SELECT MAX(step) FROM guest").Scan(&nullStep)
	if err != nil {
		return nil, err
	}
	if nullStep.Valid {
		step = nullStep.Int64
	}
	step = step + 1

	fmt.Printf("== 第%d次 - %s ==\n", step, prize)

	// candidate tags
	rows, err := sc.db.Query(fmt.Sprintf("SELECT tag, prize, maxWin, sumWin FROM sum WHERE prize='%s'", prize))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sums []Sum
	for rows.Next() {
		var sum Sum
		err = rows.Scan(&sum.Tag, &sum.Prize, &sum.MaxWin, &sum.SumWin)
		switch {
		case err == sql.ErrNoRows:
			continue
		case err != nil:
			return nil, err
		}
		sums = append(sums, sum)
	}
	tags := make([]string, 0)
	for i := 0; i < len(sums); i++ {
		for j := 0; j < sums[i].MaxWin-sums[i].SumWin; j++ {
			tags = append(tags, sums[i].Tag)
		}
	}
	fmt.Println("有效:", len(tags), ":", tags)
	if len(tags) == 0 {
		return nil, nil
	}
	Shuffle(tags)
	num := maxStepWin
	if num > len(tags) {
		num = len(tags)
	}
	candidateTags := tags[:num]
	fmt.Println("候选:", len(candidateTags), ":", candidateTags)

	// select winner and save
	var winners []int
	for _, tag := range candidateTags {
		var gid sql.NullInt64
		var code sql.NullString
		sqlStr := fmt.Sprintf("SELECT Gid, Code FROM guest WHERE tag='%s' AND (prize IS NULL OR prize='') ORDER BY RAND() LIMIT 1", tag)
		err := sc.db.QueryRow(sqlStr).Scan(&gid, &code)
		switch {
		case err == sql.ErrNoRows:
			fmt.Println("cant pick winner for ", tag)
			continue
		case err != nil:
			return nil, err
		}
		if gid.Valid {
			sqlStr := fmt.Sprintf("UPDATE guest SET prize='%s', step=%d WHERE gid=%d", prize, step, gid.Int64)
			result, err := sc.db.Exec(sqlStr)
			if err != nil {
				return nil, err
			}
			n, err := result.RowsAffected()
			if err != nil {
				return nil, err
			}
			if n == 1 {
				winner := int(gid.Int64)
				winners = append(winners, winner)
			}
		}
	}
	fmt.Println("获奖:", len(winners), ":", winners)

	return winners, nil
}

func (sc *ShuffleController) reset() error {
	sqlStr := "UPDATE guest SET prize=NULL, step=NULL"
	_, err := sc.db.Exec(sqlStr)
	return err
}
