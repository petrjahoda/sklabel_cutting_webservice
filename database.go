package main

import (
	"database/sql"
	"time"
)

type Order struct {
	OID            int    `gorm:"primary_key;column:OID"`
	Name           string `gorm:"column:Name"`
	Barcode        string `gorm:"column:Barcode"`
	ProductID      int    `gorm:"column:ProductID"`
	OrderStatusID  int    `gorm:"column:OrderStatusID"`
	CountRequested int    `gorm:"column:CountRequested"`
	WorkplaceID    int    `gorm:"column:WorkplaceID"`
	Cavity         int    `gorm:"column:Cavity"`
}

func (Order) TableName() string {
	return "order"
}

type Device struct {
	OID        int    `gorm:"primary_key;column:OID"`
	IPAddress  string `gorm:"column:IPAddress"`
	Name       string `gorm:"column:Name"`
	DeviceType int    `gorm:"column:DeviceType"`
}

func (Device) TableName() string {
	return "device"
}

type Idle struct {
	OID        int    `gorm:"primary_key;column:OID"`
	Name       string `gorm:"column:Name"`
	Barcode    string `gorm:"column:Barcode"`
	IdleTypeID int    `gorm:"column:IdleTypeID"`
}

func (Idle) TableName() string {
	return "idle"
}

type TerminalInputIdle struct {
	OID                  int       `gorm:"primary_key;column:OID"`
	DTS                  time.Time `gorm:"column:DTS"`
	DTE                  time.Time `gorm:"column:DTE"`
	OrderID              int       `gorm:"column:OrderID"`
	IdleID               int       `gorm:"column:IdleID"`
	UserID               int       `gorm:"column:UserID"`
	Interval             float32   `gorm:"column:Interval"`
	DeviceID             int       `gorm:"column:DeviceID"`
	TerminalInputOrderID int       `gorm:"column:TerminalInputOrderID"`
	Note                 string    `gorm:"column:Note"`
}

func (TerminalInputIdle) TableName() string {
	return "terminal_input_idle"
}

type TerminalInputLogin struct {
	OID      int       `gorm:"primary_key;column:OID"`
	DTS      time.Time `gorm:"column:DTS"`
	DTE      time.Time `gorm:"column:DTE"`
	UserID   int       `gorm:"column:UserID"`
	Interval float32   `gorm:"column:Interval"`
	DeviceID int       `gorm:"column:DeviceID"`
	Note     string    `gorm:"column:Note"`
}

func (TerminalInputLogin) TableName() string {
	return "terminal_input_login"
}

type TerminalInputOrder struct {
	OID             int       `gorm:"primary_key;column:OID"`
	DTS             time.Time `gorm:"column:DTS"`
	DTE             time.Time `gorm:"column:DTE"`
	OrderID         int       `gorm:"column:OrderID"`
	UserID          int       `gorm:"column:UserID"`
	DeviceID        int       `gorm:"column:DeviceID"`
	Interval        float32   `gorm:"column:Interval"`
	Count           int       `gorm:"column:Count"`
	Fail            int       `gorm:"column:Fail"`
	AverageCycle    float32   `gorm:"column:AverageCycle"`
	WorkerCount     int       `gorm:"column:WorkerCount"`
	WorkplaceModeID int       `gorm:"column:WorkplaceModeID"`
	Note            string    `gorm:"column:Note"`
	WorkshiftID     int       `gorm:"column:WorkshiftID"`
}

func (TerminalInputOrder) TableName() string {
	return "terminal_input_order"
}

type User struct {
	OID        int           `gorm:"primary_key;column:OID"`
	Login      string        `gorm:"column:Login"`
	Password   string        `gorm:"column:Password"`
	Name       string        `gorm:"column:Name"`
	FirstName  string        `gorm:"column:FirstName"`
	Rfid       string        `gorm:"column:Rfid"`
	Role       string        `gorm:"column:Role"`
	Barcode    string        `gorm:"column:Barcode"`
	Pin        string        `gorm:"column:Pin"`
	Function   string        `gorm:"column:Function"`
	UserTypeID sql.NullInt32 `gorm:"column:UserTypeID"`
	Email      string        `gorm:"column:Email"`
	Phone      string        `gorm:"column:Phone"`
}

func (User) TableName() string {
	return "user"
}

type Workshift struct {
	OID                 int    `gorm:"primary_key;column:OID"`
	WorkshiftStart      string `gorm:"column:WorkshiftStart"`
	WorkshiftLenght     int    `gorm:"column:WorkshiftLenght"`
	Name                string `gorm:"column:Name"`
	Active              int    `gorm:"column:Active"`
	WorkplaceDivisionID string `gorm:"column:WorkplaceDivisionID"`
}

func (Workshift) TableName() string {
	return "workshift"
}

type Workplace struct {
	OID                 int    `gorm:"primary_key;column:OID"`
	Name                string `gorm:"column:Name"`
	DeviceID            int    `gorm:"column:DeviceID"`
	WorkplaceDivisionID int    `gorm:"column:WorkplaceDivisionID"`
}

func (Workplace) TableName() string {
	return "workplace"
}

type WorkplaceDivision struct {
	OID  int    `gorm:"primary_key;column:OID"`
	Name string `gorm:"column:Name"`
}

func (WorkplaceDivision) TableName() string {
	return "workplace_division"
}
