package config

import (
	"ekival-canvas/constants"

	"github.com/Salvionied/apollo/txBuilding/Backend/MaestroChainContext"
)

var (
	// OGMIOS    OgmiosChainContext.OgmiosChainContext
	// BFC       BlockFrostChainContext.BlockFrostChainContext
	MC        MaestroChainContext.MaestroChainContext
	CHAIN_CTX MaestroChainContext.MaestroChainContext
)

func ChainCTXSetup() error {
	// OGMIOS = OgmiosChainContext.NewOgmiosChainContext(*ogmigo.New(ogmigo.WithEndpoint(constants.OGMIGO_ENDPOINT)), *kugo.New(kugo.WithEndpoint(constants.KUGO_ENDPOINT)))

	// BFC = BlockFrostChainContext.NewBlockfrostChainContext(
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
		CHAIN_CTX = MC
	}
	return nil
}
