package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/nighostchris/tx-track-engine-go/blockchain"
	"github.com/nighostchris/tx-track-engine-go/database"
)

func main() {
	dbUrl := os.Getenv("DATABASE_CONNECTION")

	var trackers int

	if threadPool, getThreadPoolError := strconv.Atoi(os.Getenv("THREAD_POOL")); getThreadPoolError != nil {
		trackers = 5
	} else {
		trackers = threadPool
	}

	// Connect to the database
	db, err := database.Connect(dbUrl)

	if err != nil {
		fmt.Println(err.Error())
	}

	// Auto database migration
	migrate := database.Migrate(db)

	fmt.Println(migrate)

	defer db.Close()

	const nodeUrl = "https://alfajores-forno.celo-testnet.org"

	threads := make(chan bool, trackers)
	processed := make(chan string, trackers)

	var lastProcessedBlock = "0x8bc540"

	go releaseThread(threads, processed, false)

	for {
		// blockNumber, getLatestBlockError := blockchain.GetLatestBlock(nodeUrl)
		blockNumber := "0x8bc549"

		// if getLatestBlockError != nil {
		// 	fmt.Println(getLatestBlockError.Error())
		// }

		if blockNumber >= lastProcessedBlock && len(threads) < trackers {
			var poolCapacity = trackers - len(threads)

			for i := 0; i < poolCapacity; i++ {
				midwayBlock, _ := strconv.ParseUint(lastProcessedBlock[2:], 16, 64)
				targetBlock := "0x" + strconv.FormatUint(midwayBlock+1, 16)

				go blockchain.ProcessBlock(nodeUrl, targetBlock, []string{}, processed)

				threads <- true
				lastProcessedBlock = targetBlock
			}
		}

		time.Sleep(5 * time.Second)
	}
}

func releaseThread(threads chan bool, processed chan string, verbose bool) {
	const logHead = "[Release Thread] "
	if verbose {
		fmt.Printf("%sStarts\n", logHead)
	}

	for {
		if verbose {
			fmt.Printf("%s%d threads can be released\n", logHead, len(processed))
		}

		for i := 0; i < len(processed); i++ {
			<-threads
			<-processed
		}

		time.Sleep(1 * time.Second)
	}
}
