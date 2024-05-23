package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"super-lender/inits"
	"super-lender/schemas"
	"super-lender/utils"
	"sync"
	"time"
)

// function to synchronize CC vintages with pesaflow
func SyncSpCCVintagesWithPesaflow() {
	db := inits.CurrentDB
	var ccVintages []schemas.GetSpCcVintagesResultSchema
	selectQuery := utils.FindCCVintages(db, "select")
	err := selectQuery.Scan(&ccVintages).Error
	if err != nil {
		log.Println("Error querying database:", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(ccVintages))

	client := &http.Client{
		Timeout: 10 * time.Second, // Example timeout, adjust as needed
	}

	for _, ccVintage := range ccVintages {
		go func(ccVintage schemas.GetSpCcVintagesResultSchema) {
			defer wg.Done()

			// Construct the URL
			params := url.Values{
				"customerName":        {ccVintage.CustomerName},
				"phoneNumber":         {ccVintage.PhoneNumber},
				"nationalId":          {ccVintage.NationalId},
				"loanId":              {strconv.Itoa(ccVintage.LoanId)},
				"loanApplicationDate": {ccVintage.LoanApplicationDate},
				"loanDefaultedDate":   {ccVintage.LoanDefaultedDate},
				"loanBal":             {strconv.FormatFloat(ccVintage.LoanBal, 'f', -1, 64)},
				"agentEmail":          {ccVintage.AgentEmail},
				"branch":              {ccVintage.Branch},
				"loanStatus":          {ccVintage.LoanStatus},
			}
			url := fmt.Sprintf("https://call.pesaflow.co.ke/api/customer?%s", params.Encode())

			// Call the Pesaflow API
			response, err := client.Get(url)
			if err != nil {
				log.Printf("Error calling Pesaflow API: %s\n", err)
				return
			}
			defer response.Body.Close()

			// Check for response status
			if response.StatusCode != http.StatusOK {
				log.Printf("Error from Pesaflow API: %s\n", response.Status)
				return
			}

			// Log successful synchronization
			log.Printf("Data synchronized successfully for customer: %s\n", ccVintage.CustomerName)
		}(ccVintage)
	}

	wg.Wait()

	log.Printf("Number of ccVintage records processed: %d\n", len(ccVintages))
}
