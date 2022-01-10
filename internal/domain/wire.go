package domain

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewParserFactory,
	NewActivitiesParser,
)
