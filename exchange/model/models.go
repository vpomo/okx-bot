package model

import (
	"time"
)

type OrderType string
type OrderSide string
type KlinePeriod string

type OrderStatus int

func (s OrderStatus) String() string {
	switch s {
	case 1:
		return "pending"
	case 2:
		return "finished"
	case 3:
		return "canceled"
	case 4:
		return "part-finished"
	}
	return "unknown-status"
}

// OptionParameter 可选参数
type OptionParameter struct {
	Key   string
	Value string
}

type CurrencyPair struct {
	Symbol               string  `json:"symbol,omitempty"`          //交易对
	BaseSymbol           string  `json:"base_symbol,omitempty"`     //币种
	QuoteSymbol          string  `json:"quote_symbol,omitempty"`    //交易区：usdt/usdc/btc ...
	PricePrecision       int     `json:"price_precision,omitempty"` //价格小数点位数
	QtyPrecision         int     `json:"qty_precision,omitempty"`   //数量小数点位数
	MinQty               float64 `json:"min_qty,omitempty"`
	MaxQty               float64 `json:"max_qty,omitempty"`
	MarketQty            float64 `json:"market_qty,omitempty"`
	ContractVal          float64 `json:"contract_val,omitempty"`           //1张合约价值
	ContractValCurrency  string  `json:"contract_val_currency,omitempty"`  //合约面值计价币
	SettlementCurrency   string  `json:"settlement_currency,omitempty"`    //结算币
	ContractAlias        string  `json:"contract_alias,omitempty"`         //交割合约alias
	ContractDeliveryDate int64   `json:"contract_delivery_date,omitempty"` //合约交割日期
}

//func (pair CurrencyPair) String() string {
//	return pair.Symbol
//}

//type FuturesCurrencyPair struct {
//	CurrencyPair
//	DeliveryDate int64   //结算日期
//	OnboardDate  int64   //上线日期
//	MarginAsset  float64 //保证金资产
//}

type Ticker struct {
	Pair      CurrencyPair `json:"pair"`
	Last      float64      `json:"l"`
	Buy       float64      `json:"b"`
	Sell      float64      `json:"s"`
	High      float64      `json:"h"`
	Low       float64      `json:"lw"`
	Vol       float64      `json:"v"`
	Percent   float64      `json:"percent"`
	Timestamp int64        `json:"t"`
}

type DepthItem struct {
	Price  float64 `json:"price"`
	Amount float64 `json:"amount"`
}

type DepthItems []DepthItem

func (dr DepthItems) Len() int {
	return len(dr)
}

func (dr DepthItems) Swap(i, j int) {
	dr[i], dr[j] = dr[j], dr[i]
}

func (dr DepthItems) Less(i, j int) bool {
	return dr[i].Price < dr[j].Price
}

type Depth struct {
	Pair  CurrencyPair `json:"pair"`
	UTime time.Time    `json:"ut"`
	Asks  DepthItems   `json:"asks"`
	Bids  DepthItems   `json:"bids"`
}

type Kline struct {
	Pair      CurrencyPair `json:"pair"`
	Timestamp int64        `json:"t"`
	Open      float64      `json:"o"`
	Close     float64      `json:"s"`
	High      float64      `json:"h"`
	Low       float64      `json:"l"`
	Vol       float64      `json:"v"`
}

type Order struct {
	Pair        CurrencyPair `json:"pair,omitempty"`
	Id          string       `json:"id,omitempty"`       //订单ID
	CId         string       `json:"c_id,omitempty"`     //客户端自定义ID
	Side        OrderSide    `json:"side,omitempty"`     //交易方向: sell,buy
	OrderTy     OrderType    `json:"order_ty,omitempty"` //类型: limit , market , ...
	Status      OrderStatus  `json:"status,omitempty"`   //状态
	Price       float64      `json:"price,omitempty"`
	Qty         float64      `json:"qty,omitempty"`
	ExecutedQty float64      `json:"executed_qty,omitempty"`
	PriceAvg    float64      `json:"price_avg,omitempty"`
	Fee         float64      `json:"fee,omitempty"`
	FeeCcy      string       `json:"fee_ccy,omitempty"` //收取交易手续费币种
	CreatedAt   int64        `json:"created_at,omitempty"`
	FinishedAt  int64        `json:"finished_at,omitempty"` //订单完成时间
	CanceledAt  int64        `json:"canceled_at,omitempty"`
}

type Account struct {
	Coin             string  `json:"coin,omitempty"`
	Balance          float64 `json:"balance,omitempty"`
	AvailableBalance float64 `json:"available_balance,omitempty"`
	FrozenBalance    float64 `json:"frozen_balance,omitempty"`
}

type FuturesPosition struct {
	Pair     CurrencyPair `json:"pair,omitempty"`
	PosSide  OrderSide    `json:"pos_side,omitempty"`  //开仓方向
	Qty      float64      `json:"qty,omitempty"`       // 持仓数量
	AvailQty float64      `json:"avail_qty,omitempty"` //可平仓数量
	AvgPx    float64      `json:"avg_px,omitempty"`    //开仓均价
	LiqPx    float64      `json:"liq_px,omitempty"`    // 爆仓价格
	Upl      float64      `json:"upl,omitempty"`       //盈亏
	UplRatio float64      `json:"upl_ratio,omitempty"` // 盈亏率
	Lever    float64      `json:"lever,omitempty"`     //杠杆倍数
}

type FuturesAccount struct {
	Coin      string  `json:"coin,omitempty"` //币种
	Eq        float64 `json:"eq,omitempty"`   //总权益
	AvailEq   float64 `json:"avail_eq,omitempty"`
	FrozenBal float64 `json:"frozen_bal,omitempty"`
	MgnRatio  float64 `json:"mgn_ratio,omitempty"`
	Upl       float64 `json:"upl,omitempty"`
	RiskRate  float64 `json:"risk_rate,omitempty"`
}

// https://www.okx.com/docs-v5/en/#order-book-trading-grid-trading-post-compute-min-investment-public
type ComputeMinInvestmentRequest struct {
	InstId         string           `json:"inst_id,omitempty"`
	AlgoOrdType    string           `json:"algo_ord_type,omitempty"`
	MaxPx          string           `json:"max_px,omitempty"`
	MinPx          string           `json:"min_px,omitempty"`
	GridNum        string           `json:"grid_num,omitempty"`
	RunType        string           `json:"run_type,omitempty"`
	Direction      string           `json:"direction,omitempty"`
	Lever          string           `json:"lever,omitempty"`
	BasePos        bool             `json:"base_pos,omitempty"`
	InvestmentData []InvestmentData `json:"investment_data,omitempty"`
}

type InvestmentData struct {
	Amt string
	Ccy string
}

type ComputeMinInvestmentResponse struct {
	InvestmentData []InvestmentData `json:"investment_data,omitempty"`
	SingleAmt      string           `json:"single_amt,omitempty"`
}

// https://www.okx.com/docs-v5/en/#order-book-trading-grid-trading-get-grid-algo-order-details
// GET /api/v5/tradingBot/grid/orders-algo-details
type GridAlgoOrderDetailsRequest struct {
	AlgoOrdType string `json:"algo_ord_type,omitempty"`
	AlgoId      string `json:"algo_id,omitempty"`
}

type GridAlgoOrderDetailsResponse struct {
	AlgoId              string          `json:"algo_id,omitempty"`
	AlgoClOrdId         string          `json:"algo_cl_ord_id,omitempty"`
	InstType            string          `json:"inst_type,omitempty"`
	InstId              string          `json:"inst_id,omitempty"`
	CTime               string          `json:"c_time,omitempty"`
	UTime               string          `json:"u_time,omitempty"`
	AlgoOrdType         string          `json:"algo_ord_type,omitempty"`
	State               string          `json:"state,omitempty"`
	RebateTrans         []RebateTrans   `json:"rebate_trans,omitempty"`
	TriggerParams       []TriggerParams `json:"trigger_params,omitempty"`
	MaxPx               string          `json:"max_px,omitempty"`
	MinPx               string          `json:"min_px,omitempty"`
	GridNum             string          `json:"grid_num,omitempty"`
	RunType             string          `json:"run_type,omitempty"`
	TpTriggerPx         string          `json:"tp_trigger_px,omitempty"`
	SlTriggerPx         string          `json:"sl_trigger_px,omitempty"`
	TradeNum            string          `json:"trade_num,omitempty"`
	ArbitrageNum        string          `json:"arbitrage_num,omitempty"`
	SingleAmt           string          `json:"single_amt,omitempty"`
	PerMinProfitRate    string          `json:"per_min_profit_rate,omitempty"`
	PerMaxProfitRate    string          `json:"per_max_profit_rate,omitempty"`
	RunPx               string          `json:"run_px,omitempty"`
	TotalPnl            string          `json:"total_pnl,omitempty"`
	PnlRatio            string          `json:"pnl_ratio,omitempty"`
	Investment          string          `json:"investment,omitempty"`
	GridProfit          string          `json:"grid_profit,omitempty"`
	FloatProfit         string          `json:"float_profit,omitempty"`
	TotalAnnualizedRate string          `json:"total_annualized_rate,omitempty"`
	AnnualizedRate      string          `json:"annualized_rate,omitempty"`
	CancelType          string          `json:"cancel_type,omitempty"`
	StopType            string          `json:"stop_type,omitempty"`
	ActiveOrdNum        string          `json:"active_ord_num,omitempty"`
	QuoteSz             string          `json:"quote_sz,omitempty"`
	BaseSz              string          `json:"base_sz,omitempty"`
	CurQuoteSz          string          `json:"cur_quote_sz,omitempty"`
	CurBaseSz           string          `json:"cur_base_sz,omitempty"`
	Profit              string          `json:"profit,omitempty"`
	StopResult          string          `json:"stop_result,omitempty"`
	Direction           string          `json:"direction,omitempty"`
	BasePos             string          `json:"base_pos,omitempty"`
	Sz                  string          `json:"sz,omitempty"`
	Lever               string          `json:"lever,omitempty"`
	ActualLever         string          `json:"actual_lever,omitempty"`
	LiqPx               string          `json:"liq_px,omitempty"`
	Uly                 string          `json:"uly,omitempty"`
	InstFamily          string          `json:"inst_family,omitempty"`
	OrdFrozen           string          `json:"ord_frozen,omitempty"`
	AvailEq             string          `json:"avail_eq,omitempty"`
	Eq                  string          `json:"eq,omitempty"`
	Tag                 string          `json:"tag,omitempty"`
	ProfitSharingRatio  string          `json:"profit_sharing_ratio,omitempty"`
	CopyType            string          `json:"copy_type,omitempty"`
}

type RebateTrans struct {
	Rebate    string
	RebateCcy string
}

type TriggerParams struct {
	TriggerAction   string
	TriggerStrategy string
	DelaySeconds    string
	TriggerTime     string
	TriggerType     string
	Timeframe       string
	Thold           string
	TriggerCond     string
	TimePeriod      string
	TriggerPx       string
	StopType        string
}
