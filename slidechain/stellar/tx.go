package zioncoin

import (
	"log"

	"github.com/chain/txvm/errors"
	b "github.com/zioncoin/go/build"
	"github.com/zioncoin/go/clients/equator"
	"github.com/zioncoin/go/xdr"
)

// SignAndSubmitTx signs and submits a transaction to the Zioncoin network. If there is
// an error, SubmitTx will log the Result string to the console and return the error.
func SignAndSubmitTx(hclient equator.ClientInterface, tx *b.TransactionBuilder, seeds ...string) (*equator.TransactionSuccess, error) {
	txenv, err := tx.Sign(seeds...)
	if err != nil {
		return nil, errors.Wrap(err, "signing tx")
	}
	txstr, err := xdr.MarshalBase64(txenv.E)
	if err != nil {
		return nil, errors.Wrap(err, "marshaling pre-export txenv")
	}
	resp, submitErr := hclient.SubmitTransaction(txstr)
	if submitErr != nil {
		// Attempt to extract more detailed result information
		log.Printf("error submitting tx: %s\ntx: %s", submitErr, txstr)
		var (
			resultStr string
			err       error
		)
		if herr, ok := submitErr.(*equator.Error); ok {
			resultStr, err = herr.ResultString()
			if err != nil {
				log.Print(err, "extracting result string from equator.Error")
				resultStr = ""
			}
		}
		if resultStr == "" {
			resultStr = resp.Result
			if resultStr == "" {
				log.Print("cannot locate result string from failed tx submission")
			}
		}
		log.Println("result string: ", resultStr)
	}
	return &resp, submitErr
}
