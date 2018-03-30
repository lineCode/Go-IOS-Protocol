package protocol

import (
	"IOS/src/iosbase"
	"fmt"
)

//go:generate mockgen -destination mocks/mock_dataholder.go -package protocol -source dataholder.go

type DataHolder interface {
	Init(rd *RuntimeData, network *NetworkFilter) error
	OnNewBlock(block *iosbase.Block)
}

func DataHolderFactory(kind string) (DataHolder, error) {
	switch kind {
	case "base1.0":
		holder := DataHolderImpl{}
		return &holder, nil
	}
	return nil, fmt.Errorf("target recorder not found")
}

type DataHolderImpl struct {
	*RuntimeData
	network *NetworkFilter

	blkHashCounter map[string]int
	blkCache       map[string]*iosbase.Block
}

func (d *DataHolderImpl) Init(rd *RuntimeData, network *NetworkFilter) error {
	d.RuntimeData = rd
	d.network = network

	d.blkCache = make(map[string]*iosbase.Block)
	d.blkHashCounter = make(map[string]int)
	return nil
}

func (d *DataHolderImpl) OnNewBlock(block *iosbase.Block) {
	b58Hash := iosbase.Base58Encode(block.HeadHash())
	d.blkHashCounter[b58Hash]++
	for key, val := range d.blkHashCounter {
		if val > d.view.ByzantineTolerance() {
			d.blockChain.Push(*d.blkCache[key])

		}
	}
}