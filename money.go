package main

import (
	"strconv"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Money struct {
	pence int
}

func MoneyFrom(str string) Money {
	parts := strings.Split(str, ".")
	pence, _ := strconv.Atoi(parts[1])
	pounds, _ := strconv.Atoi(parts[0])
	return Money{pence: pounds*100 + pence}
}

func (m0 *Money) Add(m1 Money) {
	m0.pence = m0.pence + m1.pence
}

func (m Money) String() string {
	pounds := m.pence / 100
	pence := m.pence % 100
	if pence < 0 {
		pence = -1 * pence
	}
	p := message.NewPrinter(language.BritishEnglish)
	return p.Sprintf("£%d.%02d", pounds, pence)
}

func (m *Money) IsPositive() bool {
	return m.pence >= 0
}

type MoneyMap map[string]Money

func (m *MoneyMap) Add(key string, value Money) {
	total, ok := (*m)[key]
	if !ok {
		total = MoneyFrom("0.00")
	}
	total.Add(value)
	(*m)[key] = total
}
