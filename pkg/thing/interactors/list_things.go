package interactors

import "fmt"

// List fetchs the registered things and return them as an array
func (i *ThingInteractor) List(authorization string) error {
	things, err := i.thingProxy.List(authorization)
	if err != nil {
		return err
	}
	i.logger.Info("Devices obtained")
	fmt.Println(things)
	err = i.msgPublisher.SendThings(things)
	if err != nil {
		return err
	}

	return nil
}
