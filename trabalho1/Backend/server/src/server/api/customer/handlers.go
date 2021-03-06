package customer

import (
	"database/sql"
	"encoding/json"
	"github.com/elithrar/simple-scrypt"
	"math/rand"
	"net/http"
	"regexp"
	"server/authentication"
	"server/globals/sqlite"
	"server/helpers"
	"server/models"
	"strconv"
	"time"
	"unicode/utf8"
)

var visaRegex, _ = regexp.Compile("^4[0-9]{12}(?:[0-9]{3})?$")
var masterCardRegex, _ = regexp.Compile("^(?:5[1-5][0-9]{2}|222[1-9]|22[3-9][0-9]|2[3-6][0-9]{2}|27[01][0-9]|2720)[0-9]{12}$")

func registerCustomer(w http.ResponseWriter, r *http.Request) {
	claims := helpers.GetAuth(r)
	if claims != nil {
		http.Error(w, "Cannot register an authenticated customer", 400)
		return
	}

	// Decoding Input
	decoder := json.NewDecoder(r.Body)
	newCustomer := models.NewCustomer()
	err := decoder.Decode(newCustomer)
	if err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	// Check if all fields are set and correct
	error := ""
	if !(utf8.RuneCountInString(newCustomer.Name) >= 3) {
		error += "name>=3 "
	}
	if !(utf8.RuneCountInString(newCustomer.Username) >= 3) {
		error += "username>=3 "
	}
	if !(utf8.RuneCountInString(newCustomer.Password) >= 6) {
		error += "password>=6 "
	}
	if !(newCustomer.CreditCard.Month > 0 && newCustomer.CreditCard.Month <= 12) {
		error += "0<creditcard.month<=12 "
	}
	year := time.Now().Year()
	if !(newCustomer.CreditCard.Year >= year) {
		error += "creditcard.year>=" + strconv.Itoa(year) + " "
	}
	if newCustomer.CreditCard.IdCreditCardType == -1 {
		error += "missingCreditCardType "
	}

	if error != "" {
		http.Error(w, error, 400)
		return
	}

	creditCardType, err := sqlite.DB.GetCreditCardType(newCustomer.CreditCard.IdCreditCardType)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid creditCardType", 400)
		return
	} else if err != nil {
		http.Error(w, "Couldn't get creditCardType", 500)
		return
	}
	switch creditCardType.Description {
	case "Visa":
		if !visaRegex.MatchString(newCustomer.CreditCard.Code) {
			http.Error(w, "Invalid Visa card code", 400)
			return
		}
	case "Mastercard":
		if !masterCardRegex.MatchString(newCustomer.CreditCard.Code) {
			http.Error(w, "Invalid Mastercard card code", 400)
			return
		}
	default:
		http.Error(w, "missing creditCard Validation", 500)
		return
	}

	// Check if username is already registered
	_, err = sqlite.DB.GetCustomer(newCustomer.Username)
	if err != sql.ErrNoRows {
		http.Error(w, "A customer already exists with that username", 400)
		return
	}

	// Hash the password
	hashedPassword, err := scrypt.GenerateFromPassword([]byte(newCustomer.Password), scrypt.DefaultParams)
	if err != nil {
		http.Error(w, "Failed to hash and salt password", 500)
		return
	}
	newCustomer.Password = string(hashedPassword)

	// Assign a 4 digit random PIN
	newCustomer.PIN = rand.Intn(9000) + 1000

	// Insert the new customer
	err = sqlite.DB.InsertCustomer(newCustomer)
	if err != nil {
		http.Error(w, "Failed inserting new Customer", 400)
		return
	}

	// Create a token to be used in later authenticated requests
	// Should come in the HTML Header, Authentication: Bearer <token>
	tokenString, err := authentication.CreateToken(newCustomer)
	if err != nil {
		http.Error(w, "Failed to create token", 500)
		return
	}

	// Replies with JSON with token and PIN
	response := []byte(`{"token":"` + tokenString + `","PIN":` + strconv.Itoa(newCustomer.PIN) + "}")
	w.Write(response)
}

func loginCustomer(w http.ResponseWriter, r *http.Request) {
	claims := helpers.GetAuth(r)
	if claims != nil {
		http.Error(w, "Already authenticated", 400)
		return
	}

	// Decoding Input
	decoder := json.NewDecoder(r.Body)
	postCustomer := models.NewCustomer()
	err := decoder.Decode(postCustomer)
	if err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	// Check if all fields are set and correct
	error := ""
	if !(utf8.RuneCountInString(postCustomer.Username) >= 3) {
		error += "username>=3 "
	}
	if !(utf8.RuneCountInString(postCustomer.Password) >= 6) {
		error += "password>=6 "
	}
	if error != "" {
		http.Error(w, error, 400)
		return
	}

	// Check if username is already registered
	dbCustomer, err := sqlite.DB.GetCustomer(postCustomer.Username)
	if err != nil {
		http.Error(w, "Failed querying customer by username", 500)
		return
	}
	err = scrypt.CompareHashAndPassword([]byte(dbCustomer.Password), []byte(postCustomer.Password))
	if err != nil {
		http.Error(w, "Bad credentials", 401)
		return
	}

	// Create a token to be used in later authenticated requests
	// Should come in the HTML Header, Authentication: Bearer <token>
	tokenString, err := authentication.CreateToken(dbCustomer)
	if err != nil {
		http.Error(w, "Failed to create token", 500)
		return
	}

	// Replies with JSON with token and PIN
	response := []byte(`{"token":"` + tokenString + `","PIN":` + strconv.Itoa(dbCustomer.PIN) + "}")
	w.Write(response)
}

func getCustomerVouchers(w http.ResponseWriter, r *http.Request) {
	claims := helpers.GetAuth(r)

	vouchers, err := sqlite.DB.GetVouchers(claims.IdCustomer, -1)
	if err != nil {
		http.Error(w, "Failed retrieving vouchers", 500)
		return
	}
	vouchersSlice, err := json.Marshal(vouchers)
	if err != nil {
		http.Error(w, "Failed converting vouchers to JSON", 500)
		return
	}

	w.Write(vouchersSlice)
}

func getCustomerOrders(w http.ResponseWriter, r *http.Request) {
	claims := helpers.GetAuth(r)

	orders, err := sqlite.DB.GetOrders(claims.IdCustomer)
	if err != nil {
		http.Error(w, "Failed retrieving orders", 500)
		return
	}
	ordersSlice, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, "Failed converting orders to JSON", 500)
		return
	}

	w.Write(ordersSlice)
}
