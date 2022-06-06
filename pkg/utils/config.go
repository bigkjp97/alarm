package utils

// import (
// 	"alarm/pkg/utils/server"
// 	"fmt"
// 	"flag"
// 	yaml "gopkg.in/yaml.v2"
// )

// type Config struct {
// 	DBConfig server.DBServer `yaml:"db_config,omitempty`
// }

// func (c Config) String() string {
// 	b, err := yaml.Marshal(c)
// 	if err != nil {
// 		return fmt.Sprintf("<error creating config string: %s>", err)
// 	}

// 	fmt.Printf(string(b))
// 	return string(b)
// }



// func DefaultUnmarshal(dst DynamicCloneable, args []string, fs *flag.FlagSet) error {
// 	return Unmarshal(dst,
// 		// First populate the config with defaults including flags from the command line
// 		Defaults(fs),
// 		// Next populate the config from the config file, we do this to populate the `common`
// 		// section of the config file by taking advantage of the code in ConfigFileLoader which will load
// 		// and process the config file.
// 		ConfigFileLoader(args, "config.file"),
// 		// Apply any dynamic logic to set other defaults in the config. This function is called after parsing the
// 		// config files so that values from a common, or shared, section can be used in
// 		// the dynamic evaluation
// 		dst.ApplyDynamicConfig(),
// 		// Load configs from the config file a second time, this will supersede anything set by the common
// 		// config with values specified in the config file.
// 		ConfigFileLoader(args, "config.file"),
// 		// Load the flags again, this will supersede anything set from config file with flags from the command line.
// 		Flags(args, fs),
// 	)
// }