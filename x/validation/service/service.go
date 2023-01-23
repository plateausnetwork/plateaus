package service

import (
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"io"
	"net/http"
	"time"
)

//TODO: check if will be moved to config
const urlPlateausValidation = "https://vyd0yyst26.execute-api.us-east-1.amazonaws.com/development/nodes"

type Validations map[string]bool

func GetValidations(valAddr sdk.ValAddress) (Validations, error) {
	accAddr := sdk.AccAddress(valAddr.Bytes())

	c := &http.Client{
		Timeout: time.Millisecond * 5000,
	}
	r, err := c.Get(fmt.Sprintf("%s/%s", urlPlateausValidation, accAddr))

	if err != nil {
		return nil, errors.New("could not check validator permission")
	}

	defer r.Body.Close()

	if r.StatusCode == http.StatusForbidden {
		return nil, errors.New("validator isn't able to check node")
	}

	b, err := io.ReadAll(r.Body)
	resp := map[string]bool{}

	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, errors.New("json.Unmarshal node validation")
	}

	return resp, nil
}
