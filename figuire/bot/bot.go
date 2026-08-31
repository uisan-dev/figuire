package bot

import (
	"figuire/figuire/amiami"
	"figuire/figuire/store"
)

type Bot struct {
	Client *amiami.Client
	Store  *store.Store
}
