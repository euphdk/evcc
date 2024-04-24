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

func ParseAdditionalChargeRecord(records []AdditionalChargeRecord, date time.Time) map[int64]float64 {

	additionalCharges := make(map[int64]float64)


	for _, record := range records {
		if AdditionalChargeRecordInRange(record, date) {
			baseTime := time.Date(
				date.Year(),
				date.Month(),
				date.Day(),
				0, 0, 0, 0, date.Location(),
			)
			additionalCharges[baseTime.Unix()] = record.Price1 + additionalCharges[baseTime.Unix()]
			additionalCharges[baseTime.Add(1 * time.Hour).Unix()] = record.Price2 + additionalCharges[baseTime.Add(1 * time.Hour).Unix()]
			additionalCharges[baseTime.Add(2 * time.Hour).Unix()] = record.Price3 + additionalCharges[baseTime.Add(2 * time.Hour).Unix()]
			additionalCharges[baseTime.Add(3 * time.Hour).Unix()] = record.Price4 + additionalCharges[baseTime.Add(3 * time.Hour).Unix()]
			additionalCharges[baseTime.Add(4 * time.Hour).Unix()] = record.Price5 + additionalCharges[baseTime.Add(4 * time.Hour).Unix()]
			additionalCharges[baseTime.Add(5 * time.Hour).Unix()] = record.Price6 + additionalCharges[baseTime.Add(5 * time.Hour).Unix()]
			additionalCharges[baseTime.Add(6 * time.Hour).Unix()] = record.Price7 + additionalCharges[baseTime.Add(6 * time.Hour).Unix()]
			additionalCharges[baseTime.Add(7 * time.Hour).Unix()] = record.Price8 + additionalCharges[baseTime.Add(7 * time.Hour).Unix()]
			additionalCharges[baseTime.Add(8 * time.Hour).Unix()] = record.Price9 + additionalCharges[baseTime.Add(8 * time.Hour).Unix()]
			additionalCharges[baseTime.Add(9 * time.Hour).Unix()] = record.Price10 + additionalCharges[baseTime.Add(9 * time.Hour).Unix()]
			additionalCharges[baseTime.Add(10 * time.Hour).Unix()] = record.Price11 + additionalCharges[baseTime.Add(10 * time.Hour).Unix()]
			additionalCharges[baseTime.Add(11 * time.Hour).Unix()] = record.Price12 + additionalCharges[baseTime.Add(11 * time.Hour).Unix()]
			additionalCharges[baseTime.Add(12 * time.Hour).Unix()] = record.Price13 + additionalCharges[baseTime.Add(12 * time.Hour).Unix()]
			additionalCharges[baseTime.Add(13 * time.Hour).Unix()] = record.Price14 + additionalCharges[baseTime.Add(13 * time.Hour).Unix()]
			additionalCharges[baseTime.Add(14 * time.Hour).Unix()] = record.Price15 + additionalCharges[baseTime.Add(14 * time.Hour).Unix()]
			additionalCharges[baseTime.Add(15 * time.Hour).Unix()] = record.Price16 + additionalCharges[baseTime.Add(15 * time.Hour).Unix()]
			additionalCharges[baseTime.Add(16 * time.Hour).Unix()] = record.Price17 + additionalCharges[baseTime.Add(16 * time.Hour).Unix()]
			additionalCharges[baseTime.Add(17 * time.Hour).Unix()] = record.Price18 + additionalCharges[baseTime.Add(17 * time.Hour).Unix()]
			additionalCharges[baseTime.Add(18 * time.Hour).Unix()] = record.Price19 + additionalCharges[baseTime.Add(18 * time.Hour).Unix()]
			additionalCharges[baseTime.Add(19 * time.Hour).Unix()] = record.Price20 + additionalCharges[baseTime.Add(19 * time.Hour).Unix()]
			additionalCharges[baseTime.Add(20 * time.Hour).Unix()] = record.Price21 + additionalCharges[baseTime.Add(20 * time.Hour).Unix()]
			additionalCharges[baseTime.Add(21 * time.Hour).Unix()] = record.Price22 + additionalCharges[baseTime.Add(21 * time.Hour).Unix()]
			additionalCharges[baseTime.Add(22 * time.Hour).Unix()] = record.Price23 + additionalCharges[baseTime.Add(22 * time.Hour).Unix()]
			additionalCharges[baseTime.Add(23 * time.Hour).Unix()] = record.Price24 + additionalCharges[baseTime.Add(23 * time.Hour).Unix()]
		}
	}

	return additionalCharges

}
