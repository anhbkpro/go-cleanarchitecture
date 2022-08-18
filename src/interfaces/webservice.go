package interfaces

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/anhbkpro/go-cleanarchitecture/src/usercases"
)

type OrderInteractor interface {
	Items(userId, orderId int) ([]usercases.Item, error)
	Add(userId, orderId, itemId int) error
}

type WebServiceHandler struct {
	OrderInteractor OrderInteractor
}

func (handler WebServiceHandler) ShowOrder(rw http.ResponseWriter, r *http.Request) {
	userId, _ := strconv.Atoi(r.FormValue("userId"))
	orderId, _ := strconv.Atoi(r.FormValue("orderId"))
	items, _ := handler.OrderInteractor.Items(userId, orderId)
	for _, item := range items {
		io.WriteString(rw, fmt.Sprintf("item id: %d\n", item.Id))
		io.WriteString(rw, fmt.Sprintf("item name: %v\n", item.Name))
		io.WriteString(rw, fmt.Sprintf("item value: %f\n", item.Value))
	}
}
