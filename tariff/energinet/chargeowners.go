package energinet

import (
	"fmt"
	"time"
)

type AdditionalChargeOwner struct {
	GLN            string
	Company        string
	ChargeTypeCode []string
	ChargeType     []string
}

// Values collected from:
// https://github.com/MTrab/energidataservice/blob/master/custom_components/energidataservice/tariffs/energidataservice/chargeowners.py

var AdditionalChargeOwners = map[string]*AdditionalChargeOwner{
	"N1": {GLN: "5790001089030", Company: "N1 A/S - 131", ChargeTypeCode: []string{"CD", "CD R"}, ChargeType: []string{"D03"}},
}

func AdditionalChargeRecordInRange(record AdditionalChargeRecord, inRangeDate time.Time) bool {

	validFrom, err := time.Parse(TimeFormatSecond, record.ValidFrom)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	validTo := time.Now().Add(28 * time.Hour)
	if record.ValidTo != "" {

		validTo, err = time.Parse(TimeFormatSecond, record.ValidTo)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
	}

	if validFrom.Before(inRangeDate) && validTo.After(inRangeDate) {
		return true
	}

	return false
}
