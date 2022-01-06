package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/nighostchris/tx-track-engine-go/blockchain"
	"github.com/nighostchris/tx-track-engine-go/database"
	"github.com/nighostchris/tx-track-engine-go/database/models"
)

func main() {
	dbUrl := os.Getenv("DATABASE_CONNECTION")

	var trackers int
	var lastProcessedBlock string

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

	// Determine whether the pipeline should start from latest block or last processed block in database
	gotLastProcessedBlock := false

	for !gotLastProcessedBlock {
		lastProcessedBlockInDatabase, getBlocksError := models.GetBlocks(db, "celo_blocks", 1)

		if getBlocksError != nil {
			fmt.Println(getBlocksError.Error())
		} else {
			if len(lastProcessedBlockInDatabase) == 0 {
				gotLatestBlock := false

				for !gotLatestBlock {
					blockNumber, getLatestBlockError := blockchain.GetLatestBlock(nodeUrl)

					if getLatestBlockError != nil {
						fmt.Println(getLatestBlockError.Error())
					} else {
						lastProcessedBlock = blockNumber
						fmt.Printf("Will be using latest block from blockchain node as starting point - %s\n", lastProcessedBlock)
						gotLastProcessedBlock = true
						gotLatestBlock = true
					}
				}
			} else {
				lastProcessedBlock = "0x" + strconv.FormatInt(lastProcessedBlockInDatabase[0].Height, 16)
				fmt.Printf("Will be using last processed block from database as starting point - %s\n", lastProcessedBlock)
				gotLastProcessedBlock = true
			}
		}
	}

	go releaseThread(threads, processed, false)

	for {
		var blockNumber string
		gotLatestBlock := false

		for !gotLatestBlock {
			latestBlock, getLatestBlockError := blockchain.GetLatestBlock(nodeUrl)

			if getLatestBlockError != nil {
				fmt.Println(getLatestBlockError.Error())
			} else {
				blockNumber = latestBlock
				gotLatestBlock = true
			}
		}

		if blockNumber >= lastProcessedBlock && len(threads) < trackers {
			blockNumberDec, _ := strconv.ParseInt(blockNumber[2:], 16, 64)
			lastProcessedBlockDec, _ := strconv.ParseInt(lastProcessedBlock[2:], 16, 64)

			fmt.Printf("[Main] Found block difference of %d\n", blockNumberDec-lastProcessedBlockDec+1)

			var poolCapacity = trackers - len(threads)

			for i := 0; i < poolCapacity && lastProcessedBlock <= blockNumber; i++ {
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
