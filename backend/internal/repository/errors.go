package repository

import "errors"

var (
	ErrInvalidPlayerCount = errors.New("invalid player count")
	ErrPlayerAlreadyExists = errors.New("player already exists")
	ErrGameFull = errors.New("game is full")
	ErrGameInProgress = errors.New("game is in progress")
	ErrPlayerNotFound = errors.New("player not found")
)