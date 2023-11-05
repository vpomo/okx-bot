package grid

import "okx-bot/exchange/model"

func (g *Grid) GetCompMinInvest(minInvestReq model.ComputeMinInvestmentRequest) (model.ComputeMinInvestmentResponse, []byte, error) {
	minInvestment, respBody, err := g.OKxV5.GetCompMinInvest(minInvestReq)
	return minInvestment, respBody, err
}
