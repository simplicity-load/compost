// package models

// import (
// 	"encoding/json"
// 	"errors"

// 	"strings"
// )

// type taskStatus int

// const (
// 	Unfinished taskStatus = iota
// 	Finished
// 	OnHold
// )

// var toStr = map[taskStatus]string{
// 	Unfinished: "unfinished",
// 	Finished:   "finished",
// 	OnHold:     "on_hold",
// }
// var fromStr = map[string]taskStatus{
// 	toStr[Unfinished]: Unfinished,
// 	toStr[Finished]:   Finished,
// 	toStr[OnHold]:     OnHold,
// }

// // NEVER AGAIN DO THIS SHIT
// // type TaskStatus interface {
// // 	Status() taskStatus
// // }

// // func (t taskStatus) Status() taskStatus {
// // 	return t
// // }

// func (t taskStatus) String() string {
// 	return toStr[t]
// }

// func (t taskStatus) MarshalJSON() (b []byte, e error) {
// 	// update below if highest len return value is changed
// 	// e.g. unfinished - lol_unfinished
// 	// unfinished     -> 10 chars
// 	// lol_unfinished -> 14 chars
// 	b = make([]byte, 0, 10+2) // +2 - 2x'"'
// 	b = append(b, '"')
// 	b = append(b, []byte(t.String())...)
// 	b = append(b, '"')
// 	return
// }

// func (t *taskStatus) UnmarshalJSON(b []byte) (e error) {
//   var in_str string
//   if e = json.Unmarshal(b, &in_str); e != nil {
//     return e
//   }
//   if e = t.fromString(in_str); e != nil {
//     return e
//   }
// 	return
// }

// func (t *taskStatus) fromString(s string) error {
// 	s = strings.TrimSpace(strings.ToLower(s))
// 	v, ok := fromStr[s]
// 	if !ok {
// 		return errors.New("Non valid taskStatus")
// 	}
//   *t = v
// 	return nil
// }
