package eth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/ethereum-optimism/optimism/l2geth/common"
)

func RequestTxTraceCache(ctx context.Context, endpoint string, chainId *big.Int, hash common.Hash) (interface{}, error) {
	requrl := fmt.Sprintf("%s/trace/tx?chainId=%s&hash=%s", endpoint, chainId.String(), hash.Hex())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("transaction %#x not found from cache server", hash)
	}

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		return json.RawMessage(raw), nil
	}

	return nil, fmt.Errorf("failed to get tx trace from cache server: %s", string(raw))
}
