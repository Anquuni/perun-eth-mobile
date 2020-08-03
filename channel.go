// Copyright (c) 2020 Chair of Applied Cryptography, Technische Universität
// Darmstadt, Germany. All rights reserved. This file is part of
// perun-eth-mobile. Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package prnm

import (
	"github.com/pkg/errors"

	ethwallet "perun.network/go-perun/backend/ethereum/wallet"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
)

// State wraps a go-perun/channel.State
// ref https://pkg.go.dev/perun.network/go-perun/channel?tab=doc#State
type State struct {
	s *channel.State
}

// GetID returns the immutable id of the channel this state belongs to.
func (s *State) GetID() []byte {
	return s.s.ID[:]
}

// GetVersion returns the version counter.
func (s *State) GetVersion() int64 {
	return int64(s.s.Version)
}

// GetBalances returns a BigInts with length two containing the current
// balances.
func (s *State) GetBalances() *BigInts {
	return &BigInts{values: s.s.Balances[0]}
}

// IsFinal indicates that the channel is in its final state.
// Such a state can immediately be settled on the blockchain.
// A final state cannot be further progressed.
func (s *State) IsFinal() bool {
	return s.s.IsFinal
}

// Params wraps a go-perun/channel.Params
// ref https://pkg.go.dev/perun.network/go-perun/channel?tab=doc#Params
type Params struct {
	params *channel.Params
}

// GetID returns the channelID of this channel.
// ref https://pkg.go.dev/perun.network/go-perun/channel?tab=doc#Params.ID
func (p *Params) GetID() []byte {
	id := p.params.ID()
	return id[:]
}

// GetChallengeDuration how many seconds an on-chain dispute of a
// non-final channel can be refuted.
func (p *Params) GetChallengeDuration() int64 {
	return int64(p.params.ChallengeDuration)
}

// GetParts returns the channel participants.
func (p *Params) GetParts() *Addresses {
	addrs := make([]ethwallet.Address, len(p.params.Parts))
	for i := range addrs {
		addrs[i] = *p.params.Parts[i].(*ethwallet.Address)
	}
	return &Addresses{values: addrs}
}

// PaymentChannel is a convenience wrapper for go-perun/client.Channel
// which provides all necessary functionality of a two-party payment channel.
// ref https://pkg.go.dev/perun.network/go-perun/client?tab=doc#Channel
type PaymentChannel struct {
	ch *client.Channel
}

// Watch starts the channel watcher routine. It subscribes to RegisteredEvents
// on the adjudicator. If an event is registered, it is handled by making sure
// the latest state is registered and then all funds withdrawn to the receiver
// specified in the adjudicator that was passed to the channel.
//
// If handling failed, the watcher routine returns the respective error. It is
// the user's job to restart the watcher after the cause of the error got fixed.
// ref https://pkg.go.dev/perun.network/go-perun/client?tab=doc#Channel.Watch
func (c *PaymentChannel) Watch() error {
	return c.ch.Watch()
}

// Send pays `amount` to the counterparty. Only positive amounts are supported.
func (c *PaymentChannel) Send(ctx *Context, amount *BigInt) error {
	if amount.i.Sign() < 1 {
		return errors.New("Only positive amounts supported in send")
	}

	state := c.ch.State().Clone()
	my := c.ch.Idx()
	other := 1 - my
	bals := state.Allocation.Balances[0]
	bals[my].Sub(bals[my], amount.i)
	bals[other].Add(bals[other], amount.i)
	state.Version++

	return c.ch.Update(ctx.ctx, client.ChannelUpdate{
		State:    state,
		ActorIdx: c.ch.Idx(),
	})
}

// GetIdx returns our index in the channel.
// ref https://pkg.go.dev/perun.network/go-perun/client?tab=doc#Channel.Idx
func (c *PaymentChannel) GetIdx() int {
	return int(c.ch.Idx())
}

// Finalize finalizes the channel with the current state.
func (c *PaymentChannel) Finalize(ctx *Context) error {
	state := c.ch.State().Clone()
	state.IsFinal = true
	state.Version++
	return c.ch.Update(ctx.ctx, client.ChannelUpdate{
		State:    state,
		ActorIdx: c.ch.Idx(),
	})
}

// Settle settles the channel: it is made sure that the current state is
// registered and the final balance withdrawn. This call blocks until the
// channel has been successfully withdrawn.
// Call Finalize before settling a channel to avoid waiting a full
// challenge duration.
// ref https://pkg.go.dev/perun.network/go-perun/client?tab=doc#Channel.Settle
func (c *PaymentChannel) Settle(ctx *Context) error {
	return c.ch.Settle(ctx.ctx)
}

// GetState returns the current state. Do not modify it.
// ref https://pkg.go.dev/perun.network/go-perun/client?tab=doc#Channel.State
func (c *PaymentChannel) GetState() *State {
	return &State{c.ch.State()}
}

// GetParams returns the channel parameters.
// ref https://pkg.go.dev/perun.network/go-perun/client?tab=doc#Channel.Params
func (c *PaymentChannel) GetParams() *Params {
	return &Params{c.ch.Params()}
}
