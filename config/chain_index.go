package config

import (
	"github.com/ekival-labs/ekival-canvas/constants"

	"github.com/Salvionied/apollo/txBuilding/Backend/MaestroChainContext"
)

var (
	// OKC    OgmiosChainContext.OgmiosChainContext
	// CHAIN_CTX BlockFrostChainContext.BlockFrostChainContext
	CHAIN_CTX MaestroChainContext.MaestroChainContext
)

func ChainCTXSetup() error {
	// OKC = OgmiosChainContext.NewOgmiosChainContext(*ogmigo.New(ogmigo.WithEndpoint(constants.BLINKLABS_OGMIOS_ENDPOINT)), *kugo.New(kugo.WithEndpoint(constants.BLINKLABS_KUPO_ENDPOINT)))

	// BFC, err := BlockFrostChainContext.NewBlockfrostChainContext(
	// 	constants.BFC_API_URL,
	// 	constants.BFC_NETWORK_ID,
	// 	constants.BFC_API_KEY,
	// )

	MC, err := MaestroChainContext.NewMaestroChainContext(
		constants.MAESTRO_NETWORK_ID,
		constants.MAESTRO_API_KEY,
	)

	if err != nil {
		return err
	} else {
		// CHAIN_CTX = OKC
		CHAIN_CTX = MC
		// CHAIN_CTX = BFC
	}
	return nil
}
