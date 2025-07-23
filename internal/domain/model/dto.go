package model

import "time"

type TraitRarityRequestDTO struct {
	CollectionID string               `json:"collectionId"`
	Properties   []TraitPropertyInput `json:"properties"`
}

type TraitPropertyInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TraitRarityResponseDTO struct {
	Continuation string                  `json:"continuation,omitempty"`
	Traits       []ExtendedTraitProperty `json:"traits"`
	Code         string                  `json:"code,omitempty"`
	Message      string                  `json:"message,omitempty"`
	StatusCode   int                     `json:"-"`
}

type ExtendedTraitProperty struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Rarity string `json:"rarity"`
}

type CreatorDTO struct {
	Account string `json:"account"`
	Value   int    `json:"value"`
}

type OwnershipDTO struct {
	ID            string       `json:"id"`
	Blockchain    string       `json:"blockchain"`
	ItemID        string       `json:"itemId"`
	Contract      string       `json:"contract"`
	Collection    string       `json:"collection"`
	TokenID       string       `json:"tokenId"`
	Owner         string       `json:"owner"`
	Value         string       `json:"value"`
	CreatedAt     time.Time    `json:"createdAt"`
	LastUpdatedAt time.Time    `json:"lastUpdatedAt"`
	Creators      []CreatorDTO `json:"creators"`
	LazyValue     string       `json:"lazyValue"`
	Code          string       `json:"code,omitempty"`
	Message       string       `json:"message,omitempty"`
	StatusCode    int          `json:"-"`
}
