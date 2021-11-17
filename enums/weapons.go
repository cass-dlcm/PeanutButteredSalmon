package enums

import (
	"fmt"
	"strings"
)

type WeaponSchedule string

const (
	RandommGrizzco WeaponSchedule = "random_gold"
	SingleRandom   WeaponSchedule = "single_random"
	FourRandom     WeaponSchedule = "four_random"
	Set            WeaponSchedule = "set"
)

func GetAllWeapons() []WeaponSchedule {
	return []WeaponSchedule{
		RandommGrizzco,
		SingleRandom,
		FourRandom,
		Set,
	}
}

func stringToWeapon(weaponStr string) (WeaponSchedule, error) {
	switch weaponStr {
	case string(RandommGrizzco), string(SingleRandom), string(FourRandom), string(Set):
		return WeaponSchedule(weaponStr), nil
	}
	return "", fmt.Errorf("weapon not found: %s", weaponStr)
}

func GetWeaponArgs(weaponsStr string) ([]WeaponSchedule, error) {
	weapons := []WeaponSchedule{}
	weaponsStrArr := strings.Split(weaponsStr, " ")
	for i := range weaponsStrArr {
		weaponVal, err := stringToWeapon(weaponsStrArr[i])
		if err != nil {
			return nil, err
		}
		weapons = append(weapons, weaponVal)
	}
	return weapons, nil
}
