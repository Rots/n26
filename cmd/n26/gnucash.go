package main

import (
	"bytes"
	"fmt"
	"github.com/Rots/n26"
	"os"
	"strconv"
)

type gncWriter struct{}

func NewGncWriter() gncWriter {
	return gncWriter{}
}

func (w gncWriter) WriteTransactions(transactions *n26.Transactions) error {
	data := [][]string{}
	for _, transaction := range *transactions {
		amount := strconv.FormatFloat(transaction.Amount, 'f', -1, 64)
		var location string
		if transaction.MerchantCity != "" {
			location = transaction.MerchantCity
			if transaction.MerchantCountry != 0 {
				location += ", "
			}
		}
		if transaction.MerchantCountry != 0 {
			location += "Country Code: " + fmt.Sprint(transaction.MerchantCountry)
		}
		//Join fields with separator ';' ignoring empty values
		var description bytes.Buffer
		for _, s := range []string{transaction.PartnerName,
			transaction.MerchantName,
			transaction.PartnerIban,
			transaction.PartnerBic,
			location} {
			if s != "" {
				_, err := description.WriteString(s)
				if err != nil {
					return err
				}
				_, err = description.WriteRune(';')
				if err != nil {
					return err
				}
			}
		}
		//Trim the trailing ;
		if description.Len() > 1 {
			description.Truncate(description.Len() - 1)
		}
		data = append(data,
			[]string{
				transaction.VisibleTS.String(),
				description.String(),
				amount,
				transaction.CurrencyCode,
				transaction.Type,
			},
		)
	}
	writer, err := NewCsvWriter(os.Stdout)
	if err != nil {
		return err
	}
	return writer.WriteData([]string{"Time", "Description", "Amount", "Currency", "Type"}, data)
}
