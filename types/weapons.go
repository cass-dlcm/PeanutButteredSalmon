package types

import (
	"fmt"
	"strings"
)

// WeaponSchedule consists of an indicator as to which weapons are in the shift.
type WeaponSchedule string

// The four possibilities for weapon set types.
const (
	RandommGrizzco WeaponSchedule = "random_gold"
	SingleRandom   WeaponSchedule = "single_random"
	FourRandom     WeaponSchedule = "four_random"
	Set            WeaponSchedule = "set"
)

// GetAllWeapons returns a slice containing every WeaponSchedule constant.
func GetAllWeapons() []WeaponSchedule {
	return []WeaponSchedule{
		RandommGrizzco,
		SingleRandom,
		FourRandom,
		Set,
	}
}

// GetWeaponArgs turns a string of space seperated WeaponSchedule strings into a slice of WeaponSchedule.
func GetWeaponArgs(weaponsStr string) ([]WeaponSchedule, error) {
	weapons := []WeaponSchedule{}
	weaponsStrArr := strings.Split(weaponsStr, " ")
	for i := range weaponsStrArr {
		switch weaponsStrArr[i] {
		case string(RandommGrizzco), string(SingleRandom), string(FourRandom), string(Set):
			weapons = append(weapons, WeaponSchedule(weaponsStrArr[i]))
		default:
			return nil, fmt.Errorf("weapon not found: %s", weaponsStrArr[i])
		}
	}
	return weapons, nil
}

// IsElementExists finds whether the given WeaponSchedule is in the WeaponSchedule slice.
func (w *WeaponSchedule) IsElementExists(arr []WeaponSchedule) bool {
	for _, v := range arr {
		if v == *w {
			return true
		}
	}
	return false
}
