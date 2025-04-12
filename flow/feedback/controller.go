package feedback

import (
	"fmt"
	"log"

	"github.com/ulngollm/teleflow"
	tele "gopkg.in/telebot.v4"
)

const FlowName string = "default"

type FlowController struct {
	manager *teleflow.FlowManager
}

type eventData map[string]string

func New(manager *teleflow.FlowManager) *FlowController {
	return &FlowController{manager: manager}
}

func (r *FlowController) Init(c tele.Context) error {
	id := c.Sender().ID
	flow := teleflow.NewSimpleFlow(id, StateUndefined, FlowName)
	if err := r.manager.InitFlow(flow); err != nil {
		return err
	}
	teleflow.SaveToCtx(c, flow)
	if err := c.Send("Выберите категорию вопроса"); err != nil {
		return fmt.Errorf("send: %v", err)
	}
	return r.checkoutState(flow, StateAskedCategory)
}

func (r *FlowController) AskProduct(c tele.Context) error {
	flow := teleflow.GetCurrentFlow(c)
	d := make(eventData)
	d[StateAskedCategory] = c.Message().Text
	flow.SetData(d)

	if err := c.Send("назовите продукт"); err != nil {
		return fmt.Errorf("send: %v", err)
	}
	return r.checkoutState(flow, StateAskedProduct)
}

func (r *FlowController) AskDetails(c tele.Context) error {
	flow := teleflow.GetCurrentFlow(c)
	d := flow.Data().(eventData)
	d[StateAskedProduct] = c.Message().Text
	flow.SetData(d)

	if err := c.Send("опишите, как воспроизвести ошибку"); err != nil {
		return fmt.Errorf("send: %v", err)
	}
	return r.checkoutState(flow, StateAskedDetails)
}

func (r *FlowController) AskScreenshot(c tele.Context) error {
	flow := teleflow.GetCurrentFlow(c)
	d := flow.Data().(eventData)
	d[StateAskedDetails] = c.Message().Text
	flow.SetData(d)

	if err := c.Send("приложите скриншот"); err != nil {
		return fmt.Errorf("send: %v", err)
	}
	return r.checkoutState(flow, StateAskedScreenshot)
}

func (r *FlowController) Thank(c tele.Context) error {
	flow := teleflow.GetCurrentFlow(c)
	d := flow.Data().(eventData)
	d[StateAskedScreenshot] = c.Message().Text
	flow.SetData(d)
	log.Printf("%v", d)

	if err := c.Send("спасибо за обратную связь! мы передали ваше сообщение в поддержку"); err != nil {
		return fmt.Errorf("send: %v", err)
	}
	return r.checkoutState(flow, StateComplete)
}

func (r *FlowController) checkoutState(flow teleflow.Flow, state string) error {
	flow.SetState(state)
	return nil
}
