package internal

import "github.com/stonecool/livemusic-go/internal/model"

type Livehouse struct {
	ID   int    `json:"ID"`
	Name string `json:"name"`
}

func (h *Livehouse) init(house *model.Livehouse) {
	h.ID = house.ID
	h.Name = house.Name
}

func (h *Livehouse) Add() error {
	data := map[string]interface{}{
		"name": h.Name,
	}

	if house, err := model.AddLiveHouse(data); err != nil {
		return err
	} else {
		h.init(house)
		return nil
	}
}

func (h *Livehouse) Get() error {
	if house, err := model.GetLiveHouse(h.ID); err != nil {
		return err
	} else {
		h.init(house)
		return nil
	}
}

func (h *Livehouse) GetAll() ([]*Livehouse, error) {
	if houses, err := model.GetLiveHouses(); err != nil {
		return nil, err
	} else {
		var s []*Livehouse

		for _, msg := range houses {
			house := &Livehouse{}
			house.init(msg)
			s = append(s, house)
		}

		return s, nil
	}
}

func (h *Livehouse) Edit() error {
	data := map[string]interface{}{
		"name": h.Name,
	}

	return model.EditLiveHouse(h.ID, data)
}

func (h *Livehouse) Delete() error {
	house, err := model.GetLiveHouse(h.ID)
	if err != nil {
		return err
	}

	return model.DeleteLiveHouse(house)
}
