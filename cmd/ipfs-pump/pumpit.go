package main

import (
	"fmt"
	"log"

	"github.com/pkg/errors"

	pump "github.com/INFURA/ipfs-pump"
)

func PumpIt(enumerator pump.Enumerator, collector pump.Collector, drain pump.Drain) {
	CIDIn := make(chan pump.CID)
	CIDOut := make(chan pump.CID)
	blocks := make(chan pump.Block)

	err := enumerator.CIDs(CIDIn)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		count := 0
		for cid := range CIDIn {
			count++
			total := enumerator.TotalCount()
			if total > 0 {
				ratio := 100. * float32(count) / float32(total)
				fmt.Printf("[%d/%d - %.2f%%] %s\n", count, total, ratio, cid)
			} else {
				fmt.Printf("[%d/%d] %s\n", count, total, cid)
			}
			CIDOut <- cid
		}
		close(CIDOut)
	}()

	err = collector.Blocks(CIDOut, blocks)
	if err != nil {
		log.Fatal(err)
	}

	for block := range blocks {
		err = drain.Drain(block)
		if err != nil {
			log.Println(errors.Wrapf(err, "failed to push block %s", block.CID))
		}
	}
}
