package tariff

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/tariff/energinet"
	"github.com/evcc-io/evcc/util"
	"github.com/evcc-io/evcc/util/request"
)

type Energinet struct {
	*embed
	log            *util.Logger
	region         string
	chargeowner    *energinet.AdditionalChargeOwner
	electricitytax float64
	vat            float64
	data           *util.Monitor[api.Rates]
}

var _ api.Tariff = (*Energinet)(nil)

func init() {
	registry.Add("energinet", NewEnerginetFromConfig)
}

func NewEnerginetFromConfig(other map[string]interface{}) (api.Tariff, error) {
	var cc struct {
		embed          `mapstructure:",squash"`
		Region         string
		ChargeOwner    string
		ElectricityTax float64
		VAT            float64
	}

	if err := util.DecodeOther(other, &cc); err != nil {
		return nil, err
	}

	if cc.Region == "" {
		return nil, errors.New("missing region")
	}

	t := &Energinet{
		embed:  &cc.embed,
		log:    util.NewLogger("energinet"),
		region: strings.ToLower(cc.Region),
		data:   util.NewMonitor[api.Rates](2 * time.Hour),
	}

	if cc.ChargeOwner == "" {
		t.log.INFO.Println("No ChargeOwner configured - skipping")
	} else {
		t.chargeowner = energinet.AdditionalChargeOwners[cc.ChargeOwner]
	}

	done := make(chan error)
	go t.run(done)
	err := <-done

	return t, err
}

func (t *Energinet) run(done chan error) {
	var once sync.Once
	client := request.NewHelper(t.log)

	tick := time.NewTicker(time.Hour)
	for ; true; <-tick.C {

		additionalCharges := make(map[int64]float64)

		if t.chargeowner != nil {
			var additionalCharge energinet.AdditionalChargesFromAPI

			jsonChargeTypeCode, _ := json.Marshal(t.chargeowner.ChargeTypeCode)
			jsonChargeType, _ := json.Marshal(t.chargeowner.ChargeType)
			dhFilter := fmt.Sprintf(
				energinet.DatahubPricelistFilter,
				jsonChargeTypeCode,
				t.chargeowner.GLN,
				jsonChargeType,
			)
			dhUri := fmt.Sprintf(energinet.DatahubPricelistURI, url.QueryEscape(dhFilter))
			t.log.TRACE.Println("Constructed URI for DatahubPricelist: " + dhUri)

			if err := backoff.Retry(func() error {
				err := client.GetJSON(dhUri, &additionalCharge)
				t.log.TRACE.Printf("%#v", additionalCharge)
				if err != nil {
					t.log.ERROR.Println(err.Error())
				}
				return backoffPermanentError(err)
			}, bo); err != nil {
				once.Do(func() { done <- err })
				t.log.ERROR.Println(err)
				continue
			}

			additionalCharges = energinet.ParseAdditionalChargeRecord(additionalCharge.Records, time.Now())
			additionalChargesTomorrow := energinet.ParseAdditionalChargeRecord(additionalCharge.Records, time.Now().Add(24*time.Hour))

			maps.Copy(additionalCharges, additionalChargesTomorrow)

			t.log.TRACE.Printf("%#v", additionalCharges)

		}

		var res energinet.Prices

		ts := time.Now().Truncate(time.Hour)
		uri := fmt.Sprintf(energinet.ElspotpricesURI,
			ts.Format(energinet.TimeFormat),
			ts.Add(24*time.Hour).Format(energinet.TimeFormat),
			t.region)

		if err := backoff.Retry(func() error {
			return backoffPermanentError(client.GetJSON(uri, &res))
		}, bo()); err != nil {
			once.Do(func() { done <- err })

			t.log.ERROR.Println(err)
			continue
		}

		data := make(api.Rates, 0, len(res.Records))
		for _, r := range res.Records {
			date, _ := time.Parse(energinet.TimeFormatSecond, r.HourUTC)

			var addtionalCharge float64

			addtionalCharge = 0.0 //TODO - Tax

			if t.chargeowner != nil {
				addtionalCharge = additionalCharges[date.Unix()] + addtionalCharge
			} else {
				addtionalCharge = 0 + addtionalCharge
			}

			ar := api.Rate{
				Start: date.Local(),
				End:   date.Add(time.Hour).Local(),
				Price: t.totalPrice(((r.SpotPriceDKK / 1e3) + addtionalCharge)) * 1.25, // TODO - VAT
			}
			data = append(data, ar)
		}

		mergeRates(t.data, data)
		once.Do(func() { close(done) })
	}
}

// Rates implements the api.Tariff interface
func (t *Energinet) Rates() (api.Rates, error) {
	var res api.Rates
	err := t.data.GetFunc(func(val api.Rates) {
		res = slices.Clone(val)
	})
	return res, err
}

// Type implements the api.Tariff interface
func (t *Energinet) Type() api.TariffType {
	return api.TariffTypePriceForecast
}
