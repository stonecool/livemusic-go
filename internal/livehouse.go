package internal

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/database"
)

type Livehouse struct {
	ID   int    `json:"ID"`
	Name string `json:"name"`
}

func (h *Livehouse) init(house *database.Livehouse) {
	h.ID = house.ID
	h.Name = house.Name
}

func (h *Livehouse) Add() error {
	data := map[string]interface{}{
		"name": h.Name,
	}

	if house, err := database.AddLiveHouse(data); err != nil {
		return err
	} else {
		h.init(house)
		return nil
	}
}

func (h *Livehouse) Get() error {
	if house, err := database.GetLiveHouse(h.ID); err != nil {
		return err
	} else {
		h.init(house)
		return nil
	}
}

func (h *Livehouse) GetAll() ([]*Livehouse, error) {
	if houses, err := database.GetLiveHouses(); err != nil {
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

	return database.EditLiveHouse(h.ID, data)
}

func (h *Livehouse) Delete() error {
	exist, err := database.ExistLivehouse(h.ID)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf("not Exist")
	}

	return database.DeleteLiveHouse(h.ID)
}

func (h *Livehouse) SetId(id int) {
	h.ID = id
}

func (h *Livehouse) Exist() (bool, error) {
	return database.ExistLivehouse(h.ID)
}
