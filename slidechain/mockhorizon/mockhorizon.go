package mockequator

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/zioncoin/go/clients/equator"
	"github.com/zioncoin/go/network"
	"github.com/zioncoin/go/xdr"
)

// New returns a new instance of the mockequator
// struct, which implements equator.ClientInterface
func New() *Client {
	return &Client{
		mu:        new(sync.Mutex),
		submitted: sync.NewCond(new(sync.Mutex)),
	}
}

// Client is a mock Horizon client, implemengint the Horizon
// client interface.
//
// Our mock implementation assumes that the calling functions
// want to submit transactions and then stream to see if they
// have been successfully included in the ledger.
type Client struct {
	txs       []string
	mu        *sync.Mutex
	submitted *sync.Cond
}

// SubmitTransaction unmarshals the tx envelope string into a xdr.TransactionEnvelope,
// and then adds the transaction to the Client's internal record of transactions to
// "stream".
func (c *Client) SubmitTransaction(txeBase64 string) (equator.TransactionSuccess, error) {
	var txe xdr.TransactionEnvelope
	err := xdr.SafeUnmarshalBase64(txeBase64, &txe)
	if err != nil {
		return equator.TransactionSuccess{}, errors.Wrap(err, "submittx: unmarshaling tx envelope")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.txs = append(c.txs, txeBase64)
	c.submitted.Broadcast()
	return equator.TransactionSuccess{}, nil
}

// StreamTransactions "streams" all transactions that have been submitted to SubmitTransaction.
func (c *Client) StreamTransactions(ctx context.Context, accountID string, cursor *equator.Cursor, handler equator.TransactionHandler) error {
	txindex := 0
	ch := make(chan struct{})

	go func() {
		c.submitted.L.Lock()
		defer c.submitted.L.Unlock()
		for {
			if ctx.Err() != nil {
				return
			}
			c.submitted.Wait()
			ch <- struct{}{}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ch:
		}

		c.mu.Lock()
		txs := c.txs[txindex:]
		c.mu.Unlock()

		for _, tx := range txs {
			htx := equator.Transaction{EnvelopeXdr: tx}
			handler(htx)
			txindex++
		}
	}
}

// Unimplemented functions
func (*Client) Root() (equator.Root, error) {
	return equator.Root{
		NetworkPassphrase: network.TestNetworkPassphrase,
	}, nil
}

func (*Client) HomeDomainForAccount(aid string) (string, error) {
	return "", nil
}

func (*Client) LoadAccount(accountID string) (equator.Account, error) {
	return equator.Account{}, nil
}

func (*Client) LoadAccountOffers(accountID string, params ...interface{}) (equator.OffersPage, error) {
	return equator.OffersPage{}, nil
}

func (*Client) LoadTradeAggregations(baseAsset, counterAsset equator.Asset, resolution int64, params ...interface{}) (equator.TradeAggregationsPage, error) {
	return equator.TradeAggregationsPage{}, nil
}

func (*Client) LoadTrades(baseAsset, counterAsset equator.Asset, offerID, resolution int64, params ...interface{}) (equator.TradesPage, error) {
	return equator.TradesPage{}, nil
}

func (*Client) LoadAccountMergeAmount(p *equator.Payment) error {
	return nil
}

func (*Client) LoadMemo(p *equator.Payment) error {
	return nil
}

func (*Client) LoadOperation(operationID string) (equator.Payment, error) {
	return equator.Payment{}, nil
}

func (*Client) LoadOrderBook(selling, buying equator.Asset, params ...interface{}) (equator.OrderBookSummary, error) {
	return equator.OrderBookSummary{}, nil
}

func (*Client) LoadTransaction(transactionID string) (equator.Transaction, error) {
	return equator.Transaction{}, nil
}

func (*Client) SequenceForAccount(accountID string) (xdr.SequenceNumber, error) {
	return xdr.SequenceNumber(0), nil
}

func (*Client) StreamLedgers(ctx context.Context, cursor *equator.Cursor, handler equator.LedgerHandler) error {
	return nil
}

func (*Client) StreamPayments(ctx context.Context, accountID string, cursor *equator.Cursor, handler equator.PaymentHandler) error {
	return nil
}
