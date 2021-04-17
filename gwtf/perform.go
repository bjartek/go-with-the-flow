package gwtf

import (
	"context"
	"log"
	"time"

	"github.com/enescakir/emoji"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"
)

var emptyEvents []flow.Event

func (f *GoWithTheFlow) performTransaction(
	tx *flow.Transaction,
	mainSigner *GoWithTheFlowAccount,
	signers []*GoWithTheFlowAccount,
	arguments []cadence.Value) ([]flow.Event, error) {

	ctx := context.Background()

	c, err := client.New(f.Address, grpc.WithInsecure(), grpc.WithMaxMsgSize(maxGRPCMessageSize))
	if err != nil {
		return emptyEvents, err
	}

	//Always need to fetch this signer anew since the sequenceNumber will change
	mainSigner, err = mainSigner.EnrichWithAccountSignerAndKey(c)
	if err != nil {
		return emptyEvents, err
	}

	// everything from here and almost down is EXACTLY the same as transaction
	blockHeader, err := c.GetLatestBlockHeader(ctx, true)
	if err != nil {
		return emptyEvents, err
	}
	tx.SetReferenceBlockID(blockHeader.ID)

	tx.SetProposalKey(mainSigner.Address, mainSigner.Key.Index, mainSigner.Key.SequenceNumber)
	tx.SetPayer(mainSigner.Address)
	tx.SetGasLimit(f.Gas)
	if len(tx.Authorizers) == 0 {
		tx.AddAuthorizer(mainSigner.Address)
	}

	for _, arg := range arguments {
		tx.AddArgument(arg)
	}

	for _, signer := range signers {
		tx.AddAuthorizer(signer.Address)
	}

	for _, signer := range signers {
		//The first time the service has not fetched the account or the signer
		if signer.Account == nil {
			signer, err := signer.EnrichWithAccountSignerAndKey(c)
			if err != nil {
				return emptyEvents, err
			}
			tx.SignPayload(signer.Address, signer.Key.Index, signer.Signer)
		}
	}

	err = tx.SignEnvelope(mainSigner.Address, mainSigner.Key.Index, mainSigner.Signer)
	if err != nil {
		return emptyEvents, err
	}

	err = c.SendTransaction(ctx, *tx)
	if err != nil {
		return emptyEvents, err
	}

	id := tx.ID()
	log.Printf("%v transaction with id %s\n", emoji.Scroll, id)
	result, err := waitForSeal(ctx, c, id)
	if err != nil {
		return emptyEvents, err
	}

	if result.Error != nil {
		return emptyEvents, result.Error
	}

	return result.Events, nil
}

// WaitForSeal wait fot the process to seal
func waitForSeal(ctx context.Context, c *client.Client, id flow.Identifier) (*flow.TransactionResult, error) {
	result, err := c.GetTransactionResult(ctx, id)
	if err != nil {
		return nil, err
	}

	for result.Status != flow.TransactionStatusSealed {
		time.Sleep(time.Second)
		result, err = c.GetTransactionResult(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
