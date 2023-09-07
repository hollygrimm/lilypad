package solver

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/bacalhau-project/lilypad/pkg/server"
	"github.com/bacalhau-project/lilypad/pkg/system"
	"github.com/bacalhau-project/lilypad/pkg/web3"
	"github.com/bacalhau-project/lilypad/pkg/web3/bindings/token"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

type SolverOptions struct {
	Web3   web3.Web3Options
	Server server.ServerOptions
}

type Solver struct {
	web3SDK    *web3.ContractSDK
	web3Events *web3.EventChannels
	server     *solverServer
}

func NewSolver(
	options SolverOptions,
) (*Solver, error) {
	web3SDK, err := web3.NewContractSDK(options.Web3)
	if err != nil {
		return nil, err
	}
	web3Events, err := web3.NewEventChannels()
	if err != nil {
		return nil, err
	}
	server, err := NewSolverServer(options.Server)
	if err != nil {
		return nil, err
	}
	solver := &Solver{
		web3SDK:    web3SDK,
		web3Events: web3Events,
		server:     server,
	}
	return solver, nil
}

func (solver *Solver) Start(ctx context.Context, cm *system.CleanupManager) error {
	ticker := time.NewTicker(1 * time.Second)
	err := solver.web3Events.Start(solver.web3SDK, ctx, cm)
	if err != nil {
		return err
	}
	err = solver.server.ListenAndServe(ctx, cm)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				log.Info().Msgf("sending tx")
				tx, err := solver.web3SDK.Contracts.Token.Transfer(
					solver.web3SDK.Auth,
					common.HexToAddress("0x2546BcD3c84621e976D8185a91A922aE77ECEc30"),
					big.NewInt(1),
				)
				if err != nil {
					fmt.Printf("error sending tx: %s\n", err.Error())
					break
				}
				fmt.Printf("tx sent: %s\n", tx.Hash())
			case <-ctx.Done():
				return
			}
		}
	}()

	solver.web3Events.Token.SubscribeTransfer(func(event *token.TokenTransfer) {
		log.Debug().Msgf("New MyEvent. From: %s, Value: %d", event.From.Hex(), event.Value)
	})

	return nil
}