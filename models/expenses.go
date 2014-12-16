package models

import (
	"github.com/juju/errors"

	"database/sql/driver"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrNegativePence is returned when a Pence with a negative values is
	// given to a function where it must not be negative
	ErrNegativePence = errors.New("Pence must not be negative")

	// ErrInvalidMoneyStr indicates that a string used to represent money
	// could not be parsed and converted to Pence
	ErrInvalidMoneyStr = errors.New("Invalid string passed to PenceFromString")

	// ErrStructNotSaved is returned when an operation is performed on a struct
	// that must first be saved to the database
	ErrStructNotSaved = errors.New("Invalid operation: struct must be saved first")
)

// Pence is an amount of money used in Payments & Expenses. There are 100 Pence
// in a Pound (Sterling)
type Pence int64

// String formats pence into a human readable string
func (p Pence) String() string {
	var negativeString = ""
	if p < 0 {
		p = -p
		negativeString = "-"
	}

	return fmt.Sprintf("%s£%01d.%02d", negativeString, p/100, p%100)
}

func (p *Pence) Scan(src interface{}) error {
	if src == nil {
		return errors.New("nil value cannot be assigned to Pence value")
	}

	n, ok := src.(int64)
	if !ok {
		return errors.New("cannot convert non-int64 to Pence")
	}

	fmt.Println("Pence.Scan: got int64 of ", n)
	*p = Pence(n)
	return nil
}

func (p Pence) Value() (driver.Value, error) {
	fmt.Println("Pence.Value() called for Pence=", p)
	return int64(p), nil
}

// Validate ensures that the pence is positive.
func (p Pence) Validate() error {
	if p <= 0 {
		return ErrNegativePence
	}
	return nil
}

// PenceFromString parses a string and returns the amount of pence.
func PenceFromString(s string) (Pence, error) {
	//Quick check to see if there are any non-numerical digits
	if strings.ToUpper(s) != strings.ToLower(s) {
		return Pence(0), ErrInvalidMoneyStr
	}

	//Ensure there is, at most, 1 decimal point.
	nDec := strings.Count(s, ".")
	if nDec > 1 {
		return Pence(0), ErrInvalidMoneyStr
	} else if nDec == 0 {
		return poundsStrToPence(s)
	}

	pIndex := strings.Index(s, ".")

	if pIndex == -1 {
		//No decimal place
		return penceStrToPence(s)
	}

	if pIndex == 0 {
		//decimal point at beginning
		return penceStrToPence(strings.TrimLeft(s, "."))
	}

	strs := strings.Split(s, ".")
	pence, err := penceStrToPence(strs[1])
	if err != nil {
		return Pence(0), err
	}
	pounds, err := poundsStrToPence(strs[0])
	if err != nil {
		return Pence(0), err
	}
	return pence + pounds, nil
}

// penceStrToPence converts a string containing an amount of pence to the value
// in pence.
func penceStrToPence(s string) (Pence, error) {
	s = strings.TrimRight(s, "0")
	if len(s) > 2 {
		return Pence(0), ErrInvalidMoneyStr
	}
	if len(s) == 0 {
		return Pence(0), nil
	}

	if len(s) == 1 {
		s += "0"
	}
	//Now everything should have correct number of Zeros
	s = strings.TrimLeft(s, ".")
	ret, err := strconv.Atoi(s)
	if err != nil {
		return Pence(0), err
	}

	return Pence(ret), nil
}

// poundsStrToPence returns the number of pence in the pounds passed as an
// argument
func poundsStrToPence(s string) (Pence, error) {
	if s == "" {
		return Pence(0), nil
	}

	ret, err := strconv.Atoi(s)
	if err != nil {
		return Pence(0), err
	}

	return Pence(ret * 100), nil
}

// Category is used to group expenses
type Category int

const (
	// CategoryGroceries indicates an expense spend on groceries
	// e.g. Boursin
	CategoryGroceries Category = iota
	// CategoryAlcohol indicates an expense spend on alcohol e.g. K Cider
	CategoryAlcohol
	// CategoryDrugs indicates an expense spend on medicine e.g cough
	// medicine
	CategoryDrugs
	// CategoryHouseholdItems indicates an expense spend on household items
	// e.g. An iron
	CategoryHouseholdItems
	// CategoryBills indicates an expense spend to on bill e.g. Netflix
	CategoryBills
	// CategoryPresents indicates an expense spent on presents e.g. singing
	// lessons
	CategoryPresents
	// CategoryTickets indicates an expense spent on tickets e.g. Balls Deep
	CategoryTickets
	// CategoryMisc is used for expenses that do not fit into any of the
	// other categories
	CategoryMisc
	// CategoryUnknown is not a category that should be used and indicates
	// a logic error.
	CategoryUnknown //Should always be last
)

func (c Category) String() string {
	switch c {
	case CategoryGroceries:
		return "Groceries"
	case CategoryAlcohol:
		return "Alcohol"
	case CategoryDrugs:
		return "Drugs"
	case CategoryHouseholdItems:
		return "Household Items"
	case CategoryBills:
		return "Bills"
	case CategoryPresents:
		return "Presents"
	case CategoryMisc:
		return "Misc"
	case CategoryTickets:
		return "Tickets"
	default:
		return "Unknown Category"
	}
}

// Validate ensures the Category is valid
func (c Category) Validate() error {
	if c == CategoryUnknown {
		return fmt.Errorf("Invalid category: %s.", c)
	}
	return nil
}

// AllCategories returns a slice containing all the valid categories
func AllCategories() []Category {
	var ret []Category
	for c := CategoryGroceries; c < CategoryUnknown; c++ {
		ret = append(ret, c)
	}
	return ret
}

// StringToCategory converts a string to the appropriate category.
// This performs a lower case compare and strips any surrounding double quote
// characters.
func StringToCategory(s string) Category {
	// Trim any quotes
	ret, ok := strToCategory[strings.ToLower(strings.Trim(s, `"`))]
	if !ok {
		return CategoryUnknown
	}
	return ret
}

var strToCategory = make(map[string]Category)

// Expense represents an expense made that is to be shared with the group
type Expense struct {
	ID          int64                `db:"id"`
	Amount      Pence                `db:"amount"`
	PayerID     int64                `db:"payer_id"`
	GroupID     int64                `db:"group_id"`
	Category    Category             `db:"category"`
	Description string               `db:"description"`
	CreatedAt   time.Time            `db:"created_at"`
	Assignments []*ExpenseAssignment `db:"-"`
}

// ExpenseAssignment represents the amount of money assigned to each user when
// an expense is made.
type ExpenseAssignment struct {
	ID        int64 `db:"id"`
	UserID    int64 `db:"user_id"`
	Amount    Pence `db:"amount"`
	ExpenseID int64 `db:"expense_id"`
}

func (e Expense) validate() error {
	err := e.Amount.Validate()
	if err != nil {
		return err
	}

	err = e.Category.Validate()
	if err != nil {
		return err
	}

	if e.PayerID <= 0 || e.GroupID <= 0 {
		return errors.New("PayerId and GroupId must be positive")
	}

	return nil
}

// Assign assigns an expense to the users given. If the amount is not equally
// divisable by the number of users the expense applies to, then the
// remaining amount is assigned at random.
func (e *Expense) Assign(userIds []int64) ([]*ExpenseAssignment, error) {
	err := e.validate()
	if err != nil {
		return nil, err
	}

	if e.ID <= 0 {
		return nil, errors.Trace(ErrStructNotSaved)
	}

	numUsers := int64(len(userIds))
	average := e.Amount / Pence(numUsers)
	remainder := int64(e.Amount) % numUsers
	randIndexes := rand.Perm(int(numUsers))

	var ret []*ExpenseAssignment
	for i := int64(0); i < numUsers; i++ {
		amount := average
		if i < remainder {
			amount++
		}
		ret = append(ret, &ExpenseAssignment{
			UserID:    int64(randIndexes[i]),
			Amount:    amount,
			ExpenseID: e.ID,
		})
	}
	return ret, nil
}

func init() {
	for _, c := range AllCategories() {
		strToCategory[strings.ToLower(c.String())] = c
	}
}
