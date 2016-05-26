package pongor

// import (
// 	"gopkg.in/flosch/pongo2.v3"
// 	"time"
// )

// func init() {
// 	pongo2.RegisterFilter("DateFormat", filterDate)
// }

// func filterDate(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
// 	t, isTime := in.Interface().(time.Time)
// 	if !isTime {
// 		return nil, &pongo2.Error{
// 			Sender:   "filter:date",
// 			ErrorMsg: "Filter input argument must be of type 'time.Time'.",
// 		}
// 	}
// 	return pongo2.AsValue(t.Format(param.String())), nil
// }
