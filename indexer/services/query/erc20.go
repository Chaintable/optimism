package query

import (
	"github.com/ethereum-optimism/optimism/indexer/db"
<<<<<<< HEAD:indexer/services/l2/query.go
<<<<<<< HEAD:indexer/services/l2/query.go
	"github.com/ethereum-optimism/optimism/l2geth/accounts/abi/bind"
	l2common "github.com/ethereum-optimism/optimism/l2geth/common"
	l2ethclient "github.com/ethereum-optimism/optimism/l2geth/ethclient"
=======

=======
>>>>>>> @eth-optimism/l2geth@0.5.27:indexer/services/query/erc20.go
	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
>>>>>>> v0.5.23:go/indexer/services/l2/query.go
)

func NewERC20(address common.Address, client *ethclient.Client) (*db.Token, error) {
	contract, err := bindings.NewERC20(address, client)
	if err != nil {
		return nil, err
	}

	name, err := contract.Name(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	symbol, err := contract.Symbol(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	decimals, err := contract.Decimals(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	return &db.Token{
		Name:     name,
		Symbol:   symbol,
		Decimals: decimals,
	}, nil
}
