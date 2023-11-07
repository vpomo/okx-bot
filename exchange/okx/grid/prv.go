package grid

import (
	"okx-bot/exchange/model"
	"okx-bot/exchange/okx/common"
)

type PrvApi struct {
	*common.Prv
}

func (prv *PrvApi) GetGridAlgoOrderDetails(req model.GridAlgoOrderDetailsRequest) (model.GridAlgoOrderDetailsResponse, []byte, error) {
	details, respBody, err := prv.Prv.GetGridAlgoOrderDetails(req)
	return details, respBody, err
}
