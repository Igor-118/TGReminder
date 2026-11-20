package service

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

const maxSymbols = 500
const reminderParts = 3
const dateLayout = "02.01.2006 15:04"

var (
	errBadFormat = errors.New("bad format")  // не "ДД.ММ ЧЧ:ММ ...", или такой даты нет. Правильно писать например вот так: 02.01.2028 15:00 Вот это круто!
	errInPast    = errors.New("in the past") // время уже прошло
	errEmptyText = errors.New("empty text")  // дата и время есть, а текста нет
	errTooLong   = errors.New("text too long")
)

func parseReminder(args string, now time.Time, loc *time.Location,
) (remindAt time.Time, text string, err error) {

	trimmedString := strings.TrimSpace(args)

	if trimmedString == "" {
		return time.Time{}, "", errEmptyText
	}

	splitMessage := strings.SplitN(trimmedString, " ", reminderParts)

	if len(splitMessage) < 2 {
		return time.Time{}, "", errBadFormat
	}

	if len(splitMessage) < reminderParts || strings.TrimSpace(splitMessage[2]) == "" {
		return time.Time{}, "", errEmptyText
	}

	datePart := splitMessage[0]
	timePart := splitMessage[1]

	text = strings.TrimSpace(splitMessage[2])

	if utf8.RuneCountInString(text) > maxSymbols {
		return time.Time{}, "", errTooLong
	}

	dateTime := fmt.Sprintf("%s %s", datePart, timePart)

	remindAt, err = time.ParseInLocation(
		dateLayout,
		dateTime,
		loc,
	)

	switch {
	case err != nil:
		return time.Time{}, "", errBadFormat
	case !remindAt.After(now.In(loc)):
		return time.Time{}, "", errInPast
	}

	return remindAt, text, nil
}
