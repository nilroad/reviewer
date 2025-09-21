package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

func loadString(envName string) string {
	validate(envName)

	return viper.GetString(envName)
}

func loadInt(envName string) int {
	validate(envName)

	return viper.GetInt(envName)
}

//nolint:all
func loadInt64(envName string) int64 {
	validate(envName)

	return viper.GetInt64(envName)
}

//nolint:all
func loadUint64(envName string) uint64 {
	validate(envName)

	return viper.GetUint64(envName)
}

func loadBool(envName string) bool {
	validate(envName)

	return viper.GetBool(envName)
}

func loadFloat64(envName string) float64 {
	validate(envName)

	return viper.GetFloat64(envName)
}

func loadStringSlice(envName string) []string {
	validate(envName)

	return viper.GetStringSlice(envName)
}

func validate(envName string) {
	exists := viper.IsSet(envName)
	if !exists {
		panic(fmt.Sprintf("environment variable [%s] does not exist", envName))
	}
}

//nolint:all
func loadDuration(envName string) time.Duration {
	validate(envName)

	return viper.GetDuration(envName)
}
