package main

import (
	"fmt"
	"net/http"
	"os"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
)

func chargeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	customerParams := &stripe.CustomerParams{Email: r.Form.Get("stripeEmail")}
	customerParams.SetSource(r.Form.Get("stripeToken"))

	var customerID string

	if true {
		// Save the mapping of customerParams Email and customerID
		// in your database and reuse for future calls
		newCustomer, err := customer.New(customerParams)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		customerID = newCustomer.ID
	}

	chargeParams := &stripe.ChargeParams{
		Amount:   500,
		Currency: "usd",
		Desc:     "Demo charge on credit card",
		Customer: customerID,
	}

	if _, err := charge.New(chargeParams); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Charged the card successfully")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,
		`<html>
		<head>
			<title>Stripe Demo - Checkout</title>
		</head>
		<body>
			<form action="/charge" method="post" class="payment">
				<article>
				  <label class="amount">
					<span>Amount: $5.00</span>
				  </label>
				</article>
			
				<script
					src="https://checkout.stripe.com/checkout.js" 
					class="stripe-button"
					data-key="%s"
					data-amount="500"
					data-description="Stripe-Demo Payment"
					data-image="https://stripe.com/img/documentation/checkout/marketplace.png"
					data-locale="auto">
  				</script>
			</form>
		</body>
		</html>`, publishableKey)
}

var publishableKey string

func init() {
	publishableKey = os.Getenv("PUBLISHABLE_KEY")
	stripe.Key = os.Getenv("SECRET_KEY")
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/charge", chargeHandler)

	http.ListenAndServe(":8080", nil)
}
