package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
	"trrader/internal/adapter"
	"trrader/internal/adapter/bybit"
	"trrader/internal/adapter/bybit/domain/linear"
)

const (
	traidingSymbol   string  = "OPUSDT"
	traidindInterval         = 3
	exchangeCom              = 0.00055
	balance          float64 = 50                              // $
	riskPerDeal      float64 = 1                               //procent
	maxLoses         float64 = (balance * (riskPerDeal / 100)) //$ without exchange fee
	takeProfitSize   float64 = 2.0                             //profit Factor
)

type service struct {
	Lenprice float64
	ctx      context.Context
	Bybit    *adapter.Client
}

func New(cl *adapter.Client) *service {
	return &service{
		Lenprice: 0,
		ctx:      context.Background(),
		Bybit:    cl,
	}
}

// TODO сделать мониторинг сделки и перенос стопа в безубыток
func (s *service) StartTraiding(ticker, side string) error {
	cl := s.Bybit.Linear()
	p := &linear.GetPositionsBySymbolParams{
		Symbol: ticker,
	}
	resp, err := cl.Account().GetPositionsBySymbol(s.ctx, p)
	if err != nil {
		return err
	}
	checker := 0
	if side == "Sell" {
		checker = 1
	}
	if resp.Result[checker].EntryPrice != 0.0 {
		err = errors.New(fmt.Sprintf("position already open.\nPNL: %v",
			resp.Result[checker].UnrealisedPnl))
		return err
	} else {
		side := bybit.Side(side)
		order, _ := s.CreateOrder(ticker, side)
		time.Sleep(2 * time.Second)
		fmt.Println(order.OrderStatus)
	}
	return nil
}

// CreateLimitOrder
func (s *service) CreateOrder(ticker string, side bybit.Side) (*linear.PlaceActiveOrderResult, error) {
	cl := s.Bybit.Linear()
	param := &linear.PlaceActiveOrderParams{}
	param.Side = side
	param.Symbol = ticker
	params, err := s.SetActiveOrderParams(param)
	fmt.Println(params)
	orderResponse, err := cl.Account().PlaceActiveOrder(s.ctx, params)
	fmt.Println(orderResponse)
	if err != nil {
		fmt.Println(err)
	}

	s.Lenprice = 0
	fmt.Printf("Order created:\n Ticker:%s\n Side: %s\n Price: %v\n", orderResponse.Result.Symbol, orderResponse.Result.Side, orderResponse.Result.Price)
	return &orderResponse.Result, err
}

// GetOrderParams
func (s *service) SetActiveOrderParams(param *linear.PlaceActiveOrderParams) (*linear.PlaceActiveOrderParams, error) {
	cl := s.Bybit.Linear()
	//price
	inforesp, err := cl.Market().GetSymbolInformation(s.ctx, &linear.GetSymbolInformationParams{Symbol: param.Symbol})
	if err != nil {
		return &linear.PlaceActiveOrderParams{}, err
	}
	index := strings.Index(inforesp.Result[0].LastPrice, ".")
	value := len(inforesp.Result[0].LastPrice[index+1:])
	s.Lenprice = float64(value)
	price, err := strconv.ParseFloat(inforesp.Result[0].LastPrice, 8)
	if err != nil {
		return &linear.PlaceActiveOrderParams{}, err
	}
	param.Price = price
	//stoploss
	stoploss, err := s.GetStopLoss(param.Symbol, param.Side)
	if err != nil {
		return &linear.PlaceActiveOrderParams{}, err
	}
	param.StopLoss = stoploss
	//qty
	qty := s.GetQty(param.Symbol, param.Side, price, stoploss)
	param.Qty = qty
	//takeprofit
	take := s.GetTakeProfit(param.Symbol, param.Side, price, stoploss)
	param.TakeProfit = take
	//set params
	params, err := s.SetOtherParams(param)
	if err != nil {
		return &linear.PlaceActiveOrderParams{}, err
	}
	fmt.Println("params in SetActiveOrder:", params)
	return params, nil
}

// get order parametr- stoploss
func (s *service) GetStopLoss(symbol string, side bybit.Side) (stop float64, err error) {
	cl := s.Bybit.Linear()
	limit := 5
	interval := traidindInterval
	timestamp := ((int(time.Now().UTC().UnixNano() / int64(time.Millisecond))) - 3000) / 1000
	from := timestamp - (interval * limit * 60)
	resp, err := cl.Market().QueryKline(s.ctx, &linear.QueryKlineParams{Symbol: symbol, Interval: interval, From: from, Limit: limit})
	if err != nil {
		return 0.0, err
	}
	klines := resp.Result
	if side == "Sell" {
		var forsort linear.HigherByPrice
		forsort = klines
		sort.Sort(forsort)
		stop = forsort[0].High + (forsort[0].High * 0.0005)
	} else if side == "Buy" {
		var forsort linear.LowByPrice
		forsort = klines
		sort.Sort(forsort)
		stop = forsort[0].Low - (forsort[0].Low * 0.0005)
	}
	stop = math.Floor(stop*math.Pow(10, s.Lenprice)) / math.Pow(10, s.Lenprice)
	return stop, nil

}

// get order parametr- qty (with exchange commision)
func (s *service) GetQty(symbol string, side bybit.Side, price, stop float64) (qty float64) {
	qty = balance / price
	loses := 0.0
	for {
		loses = s.CalculateLoses(loses, qty, price, stop, side)
		loses = math.Floor(loses*math.Pow(10, s.Lenprice)) / math.Pow(10, s.Lenprice)
		if loses <= maxLoses {
			break
		}
		qty = qty - (qty * 0.1)
	}
	if qty > 1.0 {
		qty = math.Floor(qty*math.Pow(10, 2)) / math.Pow(10, 2)
		return qty
	}
	qty = math.Floor(qty*math.Pow(10, s.Lenprice)) / math.Pow(10, s.Lenprice)
	return qty
}

// get order parametr- takeProfit
func (s *service) GetTakeProfit(symbol string, side bybit.Side, price, stop float64) (take float64) {
	difference := 0.0
	if side == "Sell" {
		difference = (stop - price) * takeProfitSize
		take = price - difference

	} else if side == "Buy" {
		difference = (price - stop) * takeProfitSize
		take = price + difference
	}
	take = math.Floor(take*math.Pow(10, s.Lenprice)) / math.Pow(10, s.Lenprice)
	return take
}

// set other params
func (s *service) SetOtherParams(p *linear.PlaceActiveOrderParams) (*linear.PlaceActiveOrderParams, error) {
	p.OrderType = "Limit"
	p.TimeInForce = bybit.TimeInForceGoodTillCancel
	p.ReduceOnly = false
	p.CloseOnTrigger = false
	positionidx := 0
	if p.Side == "Sell" {
		positionidx = 2
	} else if p.Side == "Buy" {
		positionidx = 1
	}
	p.PositionIDx = positionidx
	return p, nil
}

// follow position after creating order
// TODO изменять stoploss
func (s *service) FollowPosition(position linear.PositionsResult, stoplossTicker int) int {
	cl := s.Bybit.Linear()
	loses := 0.0
	loses = s.CalculateLoses(loses, position.Size, position.EntryPrice, position.StopLoss, bybit.Side(position.Side))
	loses = math.Floor(loses*math.Pow(10, 4)) / math.Pow(10, 4)
	//loses = loses - 0.47
	fmt.Println(loses)
	if position.UnrealisedPnl > loses {
		// orderresp, _ := cl.Account().GetActiveOrder(s.ctx, &linear.GetActiveOrderParams{Symbol: traidingSymbol})
		// fmt.Println("active order:", orderresp)
		condresp, _ := cl.Account().GetConditionalOrder(s.ctx, &linear.GetConditionalOrderParams{Symbol: traidingSymbol, Order: "Untriggered"})

		//Всегда ли stoporder первый в массиве conditionalOrderResp??????????????
		stopOrder := condresp.Result.Data[0]
		takeOrder := condresp.Result.Data[1]
		fmt.Println("stoporder:", stopOrder)
		fmt.Println("takeOrder:", takeOrder)
		fmt.Println("stopTrigger:", stopOrder.TriggerPrice)
		fmt.Println("posStop:", position.StopLoss)

		stopresp, err := cl.Account().ReplaceConditionalOrder(s.ctx, &linear.ReplaceConditionalOrderParams{Symbol: position.Symbol, StopOrderID: stopOrder.StopOrderID, PRTriggerPrice: position.EntryPrice})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("stopresp:", stopresp.RetMsg)
		fmt.Println("replace answer:", stopresp)
		if stopresp.RetMsg == "OK" {
			fmt.Println("stop loss changed in follow position")
			stoplossTicker = 1

		}

		// if takeOrder.TriggerPrice == position.TakeProfit {
		// 	stopresp, err := cl.Account().ReplaceConditionalOrder(s.ctx, &linear.ReplaceConditionalOrderParams{Symbol: position.Symbol, StopOrderID: takeOrder.StopOrderID, PRTriggerPrice: 1.7})
		// 	fmt.Println("replace answer:", stopresp)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}
		// }

		// fmt.Println("in change stopLoss")
		// resp, err := cl.Account().ReplaceActiveOrder(s.ctx, &linear.ReplaceActiveOrderParams{Symbol: position.Symbol, StopLoss: position.EntryPrice, OrderID: activeOrder.Result.Data[0].OrderID})
		// fmt.Println(resp)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// if stopresp.RetMsg == "ok" {
		// 	fmt.Println("stop loss changed in follow position")
		// 	stoplossTicker = 1

		// }
	}
	return stoplossTicker
}

// check position
func (s *service) CheckPosition(symbol string) (position linear.PositionsResult, err error) {
	cl := s.Bybit.Linear()
	posresp, err := cl.Account().GetPositionsBySymbol(s.ctx, &linear.GetPositionsBySymbolParams{Symbol: symbol})
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println(posresp)
	if posresp.Result[0].EntryPrice == 0.0 {
		return linear.PositionsResult{}, errors.New("no open position")
	}
	switch {
	case posresp.Result[0].Size != 0:
		position = posresp.Result[0]
	case posresp.Result[1].Size != 0:
		position = posresp.Result[1]
	default:
		return linear.PositionsResult{}, errors.New("no open position")
	}
	fmt.Println("UnrealisedPnl:", position.UnrealisedPnl)
	return
}

// manage position in gorutine

// TODO после установки ытопа массив не читается
func (s *service) ManagePosition(symbol string) {
	stoplossTicker := 0
	for {
		position, err := s.CheckPosition(symbol)
		if err != nil {
			fmt.Println(err)
			time.Sleep(5 * time.Second)
		}
		if stoplossTicker == 0 {
			ticker := s.FollowPosition(position, stoplossTicker)
			stoplossTicker = ticker
		} else {
			fmt.Println("stoploss already changed")

		}
		time.Sleep(5 * time.Second)
	}

}

// calculate loses
func (s *service) CalculateLoses(loses, qty, price, stop float64, side bybit.Side) float64 {
	commission := (qty * price * exchangeCom) + (qty * (price * 2) * exchangeCom)
	if side == "Sell" {
		loses = (qty * stop) - (qty * price)
		loses = loses + commission
	} else if side == "Buy" {
		loses = (qty * price) - (qty * stop)
		loses = loses + commission
	}
	return loses
}
func (s *service) PumpAndDump() {
	cl := s.Bybit.Linear()
	resp, err := cl.Market().GetSymbolInformation(s.ctx, &linear.GetSymbolInformationParams{})
	ticker := 0
	for i, v := range resp.Result {
		ticker = i
		fmt.Println("Ticker:", v.Symbol)
		//fmt.Println("1hPrice:", v.Price1HPcnt)
		fmt.Println("Volume:", v.Volume24H)
	}
	fmt.Println(ticker)
	if err != nil {
		fmt.Println(err)
	}
}
