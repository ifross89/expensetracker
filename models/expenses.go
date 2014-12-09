package models

import (
	"github.com/juju/errors"

	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNegativePence   = errors.New("Pence must be positive")
	ErrInvalidMoneyStr = errors.New("Invalid string passed to PenceFromString")
)

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

type Category int

const (
	CategoryGroceries Category = iota
	CategoryAlcohol
	CategoryDrugs
	CategoryHouseholdItems
	CategoryBills
	CategoryPresents
	CategoryTickets
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
	case CategoryTickets:
		return "Tickets"
	default:
		return "Unknown Category"
	}
}

func (c Category) Validate() error {
	if c == CategoryUnknown {
		return fmt.Errorf("Invalid category: %s.", c)
	}
	return nil
}

func AllCategories() []Category {
	var ret []Category
	for c := CategoryGroceries; c < CategoryUnknown; c++ {
		ret = append(ret, c)
	}
	return ret
}

func StringToCategory(s string) Category {
	// Trim any quotes
	ret, ok := strToCategory[strings.ToLower(strings.Trim(s, `"`))]
	if !ok {
		return CategoryUnknown
	}
	return ret
}

var strToCategory = make(map[string]Category)

type Expense struct {
	Id          int64
	Amount      Pence
	PayerId     int64
	GroupId     int64
	Category    Category
	Description string
	Time        time.Time
	CreatedAt   time.Time
	Assignments []*ExpenseAssignment
}

type ExpenseAssignment struct {
	Id        int64
	UserId    int64
	Amount    Pence
	ExpenseId int64
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

	if e.PayerId <= 0 || e.GroupId <= 0 {
		return errors.New("PayerId and GroupId must be positive")
	}

	return nil
}

func (e *Expense) Assign(userIds []int64) ([]*ExpenseAssignment, error) {
	err := e.validate()
	if err != nil {
		return nil, err
	}

	if e.Id <= 0 {
		return nil, errors.New("Expense must be saved before it is assigned")
	}

	numUsers := int64(len(userIds))
	average := e.Amount / Pence(numUsers)
	remainder := int64(e.Amount) % numUsers
	randIndexes := rand.Perm(int(numUsers))

	var ret []*ExpenseAssignment
	for i := int64(0); i < numUsers; i++ {
		amount := average
		if i < remainder {
			amount += 1
		}
		ret = append(ret, &ExpenseAssignment{
			UserId:    int64(randIndexes[i]),
			Amount:    amount,
			ExpenseId: e.Id,
		})
	}
	return ret, nil
}

func init() {
	for _, c := range AllCategories() {
		strToCategory[strings.ToLower(c.String())] = c
	}
}
