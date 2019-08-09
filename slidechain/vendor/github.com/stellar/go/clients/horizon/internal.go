package equator

import (
	"encoding/json"
	"net/http"

	"github.com/zioncoin/go/support/errors"
)

func decodeResponse(resp *http.Response, object interface{}) (err error) {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		equatorError := &Error{
			Response: resp,
		}
		decodeError := decoder.Decode(&equatorError.Problem)
		if decodeError != nil {
			return errors.Wrap(decodeError, "error decoding equator.Problem")
		}
		return equatorError
	}

	err = decoder.Decode(&object)
	if err != nil {
		return
	}
	return
}

func loadMemo(p *Payment) error {
	res, err := http.Get(p.Links.Transaction.Href)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(&p.Memo)
}
