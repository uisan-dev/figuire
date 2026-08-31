package bot

import (
	"figuire/figuire/amiami"
	"figuire/figuire/store"
)

type Bot struct {
	Client *amiami.Client
	Store  *store.Store
}

var FiguireBot *Bot = &Bot{
	Client: amiami.NewClient(),
	Store:  store.FiguireStore,
}
