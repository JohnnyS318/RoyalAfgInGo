package bank

import (
	"errors"

	"github.com/Rhymond/go-money"
)

//Checks if the specified amount qualifies as a raise action. It has to be greater than the original maximum bet
func (b *Bank) isRaise(amount *money.Money) bool {
	//Error can neglected because it will return false if an error occurs.
	res, _ := amount.GreaterThan(b.MaxBet)
	return res
}

//MustAllIn determines whether a player has to bet everything in because the maximum bet is already past his wallet amount
func (b *Bank) MustAllIn(id string) (bool, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	p, ok := b.PlayerWallet[id]
	if !ok {
		return false, errors.New("The player was not found")
	}
	bet, ok := b.PlayerBets[id]
	if !ok {
		return false, errors.New("The player was not found")
	}
	add, err := bet.Add(p)
	if err != nil {
		return false, err
	}
	return b.MaxBet.GreaterThanOrEqual(add)
}

//IsAllIn determines whether a given player has already placed all his wallet. He can be excluded from the blocking list
func (b *Bank) IsAllIn(id string) bool {
	b.lock.RLock()
	defer b.lock.RUnlock()
	w, ok := b.PlayerWallet[id]
	if !ok {
		return true
	}

	bet, ok := b.PlayerBets[id]
	if !ok {
		return true
	}

	return w.IsZero() && bet.IsPositive()
}