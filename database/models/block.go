package models

import (
	"database/sql"
	"fmt"
)

type Block struct {
	Id        int64  `json:"id"`
	Height    int64  `json:"height"`
	Hash      string `json:"hash"`
	Timestamp int64  `json:"timestamp"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func GetBlocks(db *sql.DB, table string, limit int64) (blocks []Block, err error) {
	const logHead = "[Block Model] Get - "
	var query string

	if limit == 0 {
		query = fmt.Sprintf("SELECT * FROM %s order by height desc", table)
	} else {
		query = fmt.Sprintf("SELECT * FROM %s order by height desc limit %d", table, limit)
	}

	fmt.Printf("%sStarts\n", logHead)

	rows, queryError := db.Query(query)

	if queryError != nil {
		return []Block{}, fmt.Errorf("%sFailed to get all blocks from database: %s", logHead, queryError.Error())
	}

	var finalBlocks []Block

	for rows.Next() {
		var nextBlock Block
		// var id int64
		// var height int64
		// var hash string
		// var timestamp int64
		// var createdAt string
		// var updatedAt string

		scanError := rows.Scan(&nextBlock)

		if scanError != nil {
			return []Block{}, fmt.Errorf("%sFailed to parse records into block interface: %s", logHead, scanError.Error())
		}

		finalBlocks = append(finalBlocks, nextBlock)
	}

	fmt.Printf("%sSuccessfully get all blocks from database\n", logHead)

	return finalBlocks, nil
}
