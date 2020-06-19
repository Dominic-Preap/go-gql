// https://github.com/vektah/dataloaden/issues/35

package dataloader

import (
	"context"
	"net/http"
	"time"

	"github.com/my/app/model"
	"github.com/my/app/service"
)

type ctxKeyType struct{ name string }

var ctxKey = ctxKeyType{"userCtx"}

// Loaders .
type Loaders struct {
	UserByID      *UserLoader
	TodosByUserID *TodosLoader
}

// Middleware .
func Middleware(s *service.Service, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ldrs := Loaders{}

		// set this to zero what happens without dataloading
		wait := 250 * time.Microsecond

		// simple 1:1 loader, fetch an address by its primary key
		ldrs.UserByID = &UserLoader{
			wait:     wait,
			maxBatch: 100,
			fetch: func(keys []int) ([]*model.User, []error) {
				users, err := s.User.FindAll(&service.UserFilter{IDs: keys})

				if err != nil {
					return nil, []error{err}
				}

				u := make(map[int]*model.User, len(users))
				for _, user := range users {
					u[user.ID] = user
				}

				result := make([]*model.User, len(keys))
				for i, id := range keys {
					result[i] = u[id]
				}

				return result, nil
			},
		}

		// 1:M loader
		ldrs.TodosByUserID = &TodosLoader{
			wait:     wait,
			maxBatch: 100,
			fetch: func(keys []int) ([][]*model.Todo, []error) {
				todos, err := s.Todo.FindAll(&service.TodoFilter{UserIDs: keys})
				if err != nil {
					return nil, []error{err}
				}

				t := make(map[int][]*model.Todo, len(keys))
				for _, todo := range todos {
					t[todo.UserID] = append(t[todo.UserID], todo)
				}

				result := make([][]*model.Todo, len(keys))
				for i, key := range keys {
					result[i] = t[key]
				}

				return result, nil
			},
		}

		// // 1:M loader
		// ldrs.ordersByCustomer = &OrderSliceLoader{
		// 	wait:     wait,
		// 	maxBatch: 100,
		// 	fetch: func(keys []int) ([][]*Order, []error) {
		// 		var keySql []string
		// 		for _, key := range keys {
		// 			keySql = append(keySql, strconv.Itoa(key))
		// 		}

		// 		fmt.Printf("SELECT * FROM orders WHERE customer_id IN (%s)\n", strings.Join(keySql, ","))
		// 		time.Sleep(5 * time.Millisecond)

		// 		orders := make([][]*Order, len(keys))
		// 		errors := make([]error, len(keys))
		// 		for i, key := range keys {
		// 			id := 10 + rand.Int()%3
		// 			orders[i] = []*Order{
		// 				{ID: id, Amount: rand.Float64(), Date: time.Now().Add(-time.Duration(key) * time.Hour)},
		// 				{ID: id + 1, Amount: rand.Float64(), Date: time.Now().Add(-time.Duration(key) * time.Hour)},
		// 			}

		// 			// if you had another customer loader you would prime its cache here
		// 			// by calling `ldrs.ordersByID.Prime(id, orders[i])`
		// 		}

		// 		return orders, errors
		// 	},
		// }

		// // M:M loader
		// ldrs.itemsByOrder = &ItemSliceLoader{
		// 	wait:     wait,
		// 	maxBatch: 100,
		// 	fetch: func(keys []int) ([][]*Item, []error) {
		// 		var keySql []string
		// 		for _, key := range keys {
		// 			keySql = append(keySql, strconv.Itoa(key))
		// 		}

		// 		fmt.Printf("SELECT * FROM items JOIN item_order WHERE item_order.order_id IN (%s)\n", strings.Join(keySql, ","))
		// 		time.Sleep(5 * time.Millisecond)

		// 		items := make([][]*Item, len(keys))
		// 		errors := make([]error, len(keys))
		// 		for i := range keys {
		// 			items[i] = []*Item{
		// 				{Name: "item " + strconv.Itoa(rand.Int()%20+20)},
		// 				{Name: "item " + strconv.Itoa(rand.Int()%20+20)},
		// 			}
		// 		}

		// 		return items, errors
		// 	},
		// }

		dlCtx := context.WithValue(r.Context(), ctxKey, ldrs)
		next.ServeHTTP(w, r.WithContext(dlCtx))
	})
}

// For ..
func For(ctx context.Context) Loaders {
	return ctx.Value(ctxKey).(Loaders)
}
