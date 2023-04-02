package cmd

import (
	"api/utils"
	"encoding/json"
	"os/exec"
	"sync"
)

var s sync.Mutex

func CallJSON[T any](command string, args ...string) (T, error) {
	s.Lock()
	defer s.Unlock()
	var result T
	r, err := exec.Command(command, args...).Output()
	if err != nil {
		return getZero[T](), err
	}
	err = json.Unmarshal(r, &result)
	if err != nil {
		return getZero[T](), err
	}
	return result, nil
}

func CallString(command string, args ...string) (string, error) {
	s.Lock()
	defer s.Unlock()

	r, err := exec.Command(command, args...).Output()
	if err != nil {
		return "", err
	}
	return string(r), nil
}

func CallArrayJSON[T any](command string, args ...string) ([]T, error) {
	s.Lock()
	defer s.Unlock()
	var result []T
	r, err := exec.Command(command, args...).Output()
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return getZeroArray[T](), err
	}
	err = json.Unmarshal(r, &result)
	if err != nil {
		utils.WrapErrorLog(err.Error())
		return getZeroArray[T](), err
	}
	return result, nil
}

func getZero[T any]() T {
	var result T
	return result
}

func getZeroArray[T any]() []T {
	var result []T
	return result
}
