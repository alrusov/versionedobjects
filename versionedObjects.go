package versionedObjects

import (
	"fmt"
	"reflect"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------------//

type (

	// Служебные поля объекта
	SysData struct {
		Ver           string     `json:"ver" db:"o.ver" dbSimple:"ver" default:"" comment:"Object version" readonly:"true"`
		ChangeGUID    string     `json:"-" db:"o.change_guid" dbSimple:"change_guid" default:"00000000-0000-0000-0000-000000000000" comment:"Change GUID" readonly:"true"`
		CreatedByID   uint64     `json:"-" db:"o.created_by" dbSimple:"created_by" default:"0" comment:"The user ID who created the object" readonly:"true"`
		CreatedBy     Person     `json:"createdBy,omitempty" db:"-" dbSimple:"-" comment:"The user who created the object" readonly:"true"`
		CreatedAt     *time.Time `json:"createdAt,omitempty" db:"o.created_at" dbSimple:"created_at" default:"1970-01-01T00:00:00" comment:"Date of creation of object" readonly:"true"`
		UpdatedByID   uint64     `json:"-" db:"o.updated_by" dbSimple:"updated_by" default:"0" comment:"The user ID who updated the object" readonly:"true"`
		UpdatedBy     Person     `json:"updatedBy,omitempty" db:"-" dbSimple:"-" comment:"The user who updated the object" readonly:"true"`
		LifeStart     time.Time  `json:"lifeStart" db:"o.life_start" dbSimple:"life_start" comment:"Object version life start" readonly:"true"`
		LifeEnd       time.Time  `json:"lifeEnd" db:"o.life_end" dbSimple:"life_end" comment:"Object version life end" readonly:"true"`
		ActionComment string     `json:"actionComment,omitempty" db:"o.action_comment" dbSimple:"action_comment" default:"" comment:"Comment for version change"`
		Changes       Changes    `json:"changes,omitempty" db:"-" dbSimple:"-" comment:"Changes from the previous version" readonly:"true" ref:"objChanges"`
	}

	// Персона создавшая или изменившая объект
	Person struct {
		ID   uint64 `json:"id,omitempty" db:"-" dbSimple:"-" default:"0" comment:"ID of the user" readonly:"true"`
		GUID string `json:"guid,omitempty" db:"-" dbSimple:"-" default:"" comment:"GUID of the user" readonly:"true"`
		Name string `json:"name,omitempty" db:"-" dbSimple:"-" default:"" comment:"Name of the user" readonly:"true"`
	}

	// Набор изменений
	Changes []*Change

	// Изменение
	Change struct {
		FieldName string `json:"fieldName" readonly:"true"`
		OldVal    any    `json:"oldVal" readonly:"true"`
		NewVal    any    `json:"newVal" readonly:"true"`
	}

	StdHeader struct {
		ID          uint64   `json:"id,omitempty" db:"o.id" dbSimple:"id" comment:"ID" role:"primary" readonly:"true"`
		Sys         *SysData `json:"sys,omitempty" db:"" dbSimple:"" comment:"System data" ref:"sysData"`
		GUID        string   `json:"guid,omitempty" db:"o.guid" dbSimple:"guid" comment:"GUID" role:"key" readonly:"true"`
		Name        string   `json:"name,omitempty" db:"o.name" dbSimple:"name" comment:"Name"`
		Description string   `json:"description,omitempty" db:"o.description" dbSimple:"description" default:"" comment:"Description"`
		Flags       Flags    `json:"flags,omitempty" db:"o.flags" dbSimple:"flags" default:"0" comment:"Flags"`
	}
)

//----------------------------------------------------------------------------------------------------------------------------//

// Стандартные поля

const (
	FieldID    = "id"
	FieldGUID  = "guid"
	FieldName  = "name"
	FieldFlags = "flags"

	FieldStdHeader = "StdHeader"
)

const (
	DbFieldID    = "o." + FieldID
	DbFieldGUID  = "o." + FieldGUID
	DbFieldName  = "o." + FieldName
	DbFieldFlags = "o." + FieldFlags
)

//----------------------------------------------------------------------------------------------------------------------------//

// Флаги

type Flags uint64

const (
	FlagSystem      = Flags(0x80000000) // Системный
	FlagTest        = Flags(0x40000000) // Тестовый объект
	FlagWithProtect = Flags(0x20000000) // При сохранении в базу защитить от изменения определенные флаги (набор зависит от типа объекта), сам этот флаг в базу не сохранять, реализовано в триггере
	FlagCustom1     = Flags(0x10000000) // Пользовательский флаг 1
	FlagCustom2     = Flags(0x08000000) // Пользовательский флаг 2
	FlagCustom3     = Flags(0x04000000) // Пользовательский флаг 2
)

//----------------------------------------------------------------------------------------------------------------------------//

var stdHeaderType = reflect.TypeOf(StdHeader{})

func StdHeaderForObject(obj any) (sh *StdHeader, err error) {
	objV := reflect.ValueOf(obj)
	return StdHeaderForObjectV(objV)
}

func StdHeaderForObjectV(objV reflect.Value) (sh *StdHeader, err error) {
	if objV.Kind() != reflect.Pointer {
		err = fmt.Errorf("is not a pointer")
		return
	}

	objV = objV.Elem()

	if objV.Kind() != reflect.Struct {
		err = fmt.Errorf("is not a struct")
		return
	}

	shV := objV.FieldByName(FieldStdHeader)
	if !shV.IsValid() || shV.Kind() != reflect.Struct || shV.Type() != stdHeaderType {
		// поля нет или оно не то
		err = fmt.Errorf("object has no valid %s", FieldStdHeader)
		return
	}

	sh = shV.Addr().Interface().(*StdHeader) // тип проверили выше, уже не упадет
	return
}

//----------------------------------------------------------------------------------------------------------------------------//
