// Copyright 2023 Blink Labs, LLC.
package config

import (
	"ekival-canvas/constants"

	"github.com/Salvionied/apollo/txBuilding/Backend/BlockFrostChainContext"
	"github.com/Salvionied/apollo/txBuilding/Backend/OgmiosChainContext"
	"github.com/SundaeSwap-finance/kugo"
	"github.com/SundaeSwap-finance/ogmigo/v6"
)

var (
	OGMIOS = OgmiosChainContext.NewOgmiosChainContext(*ogmigo.New(ogmigo.WithEndpoint(constants.OGMIGO_ENDPOINT)), *kugo.New(kugo.WithEndpoint(constants.KUGO_ENDPOINT)))

	BFC = BlockFrostChainContext.NewBlockfrostChainContext(constants.API_URL, constants.NETWORK_ID, constants.API_KEY)
)
