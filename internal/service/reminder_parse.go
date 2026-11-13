package service

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const maxSymbols = 500

var (
	ErrBadFormat = errors.New("bad format")  // не "ДД.ММ ЧЧ:ММ ...", или такой даты нет. Правильно писать например вот так: 02.01.2028 15:00 Дот это круто!
	ErrInPast    = errors.New("in the past") // время уже прошло
	ErrEmptyText = errors.New("empty text")  // дата и время есть, а текста нет
	ErrTooLong   = errors.New("text too long")
)

func parseReminder(args string, now time.Time, loc *time.Location,
) (remindAt time.Time, text string, err error) {

	stringNoSpaces := strings.TrimSpace(args)
	if stringNoSpaces == "" {
		return time.Time{}, "", ErrEmptyText
	}
	if len(stringNoSpaces) > maxSymbols {
		return time.Time{}, "", ErrTooLong
	}
	splitMessage := strings.SplitN(stringNoSpaces, " ", 3)
	if len(splitMessage) < 2 {
		return time.Time{}, "", ErrBadFormat
	}
	datePart := splitMessage[0]
	timePart := splitMessage[1]
	if len(splitMessage) < 3 || strings.TrimSpace(splitMessage[2]) == "" {
		return time.Time{}, "", ErrEmptyText
	}

	text = splitMessage[2]
	dateTime := fmt.Sprintf("%s %s", datePart, timePart)

	remindAt, err = time.ParseInLocation(
		"02.01.2006 15:04",
		dateTime,
		loc,
	)

	if err != nil {
		return time.Time{}, "", ErrBadFormat
	}
	if !remindAt.After(now.In(loc)) {
		return time.Time{}, "", ErrInPast
	}
	return remindAt, text, nil
}
