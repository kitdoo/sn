package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-ozzo/ozzo-validation/v4"

	"github.com/kitdoo/sn/internal/app"

	"github.com/spf13/cobra"
)

var cmdServer = &cobra.Command{
	Use:          "sn",
	Short:        "Run the sn service",
	Run:          cmdServerRun,
	SilenceUsage: true,
}

func init() {
	RootCmd.AddCommand(cmdServer)
	RootCmd.Run = cmdServerRun
}

func cmdServerRun(_ *cobra.Command, _ []string) {
	configFile := RootCmd.PersistentFlags().Lookup("config").Value.String()
	if configFile == "" {
		configFile, _ = os.LookupEnv("CONFIG_FILE")
	}

	serv := app.New(configFile)
	err := serv.Start()
	if err != nil {
		println(PrettyError(err))
		return
	}
	serv.Wait().Shutdown()
}
func PrettyError(err error) string {
	var verrs validation.Errors
	if errors.As(err, &verrs) {
		verrs = ValidationErrorsToFlatMap(verrs)
		var strs = make([]string, 0, len(verrs))
		for f, e := range verrs {
			strs = append(strs, fmt.Sprintf("  - %s: %v", f, e))
		}
		return "Configuration in not valid:\n" +
			strings.Join(strs, "\n")
	}
	return err.Error()
}

func ValidationErrorsToFlatMap(in validation.Errors) validation.Errors {
	out := make(validation.Errors)
	for key, value := range in {
		if valueAsMap, ok := value.(validation.Errors); ok {
			sub := ValidationErrorsToFlatMap(valueAsMap)
			for subKey, subValue := range sub {
				out[key+"."+subKey] = subValue
			}
			continue
		}
		out[key] = value
	}
	return out
}
