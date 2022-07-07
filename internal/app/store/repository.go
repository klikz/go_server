package store

import "premier_api/internal/app/model"

// UserRepository ...
type UserRepository interface {
	Create(*model.User) error
	Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
	ComponentsAll() (interface{}, error)
	GetLast(int) (interface{}, error)
	GetStatus(int) (interface{}, error)
	GetPackingLast() (interface{}, error)
	GetPackingToday() (interface{}, error)
	GetPackingTodayModels() (interface{}, error)
	GetToday(int) (interface{}, error)
	GetTodayModels(int) (interface{}, error)
	GetSectorBalance(int) (interface{}, error)
	SerialInput(int, string) (interface{}, error)
	PackingSerialInput(string, string) (interface{}, error)
	GetPackingTodaySerial() (interface{}, error)
	GetLines() (interface{}, error)
	GetDefectsTypes() (interface{}, error)
	DeleteDefectsTypes(int) (interface{}, error)
	AddDefectsTypes(int, string) (interface{}, error)
	AddDefects(*model.OtkAddDefect, string) (interface{}, error)
	GetByDate(string, string, int) (interface{}, error)
	GetByDateModels(string, string, int) (interface{}, error)
	GetByDateSerial(string, string) (interface{}, error)
	GetRemont() (interface{}, error)
	GetRemontToday() (interface{}, error)
	GetRemontByDate(string, string) (interface{}, error)
	UpdateRemont(string, int) (interface{}, error)
	CheckRole(string, string) (int, error)
	GetInfoBySerial(string) (interface{}, error)
	Galileo(*model.Galileo) (interface{}, error)
}
