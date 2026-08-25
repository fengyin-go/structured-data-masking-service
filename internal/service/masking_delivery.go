package service

import "datamasking/internal/store"

type MaskingDelivery struct {
	adapter *store.MaskingDeliveryAdapter
}

func NewMaskingDelivery(adapter *store.MaskingDeliveryAdapter) *MaskingDelivery {
	return &MaskingDelivery{adapter: adapter}
}

func (d *MaskingDelivery) Execute(mode string) error {
	var first error
	for attempt := 0; attempt < 2; attempt++ {
		err := d.adapter.Deliver(mode)
		if err == nil {
			return first
		}
		if first == nil {
			first = err
		}
	}
	return first
}
