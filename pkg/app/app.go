package app

import (
	"bytes"
	"fmt"
	"os"
	"reflect"

	"github.com/WeiXinao/daily_fresh/pkg/app/configurator"
	"github.com/WeiXinao/daily_fresh/pkg/app/configurator/subscriber"
	cliflag "github.com/WeiXinao/daily_fresh/pkg/common/cli/flag"
	"github.com/WeiXinao/daily_fresh/pkg/common/cli/globalflag"
	"github.com/WeiXinao/daily_fresh/pkg/common/term"
	"github.com/WeiXinao/daily_fresh/pkg/common/version"
	"github.com/WeiXinao/daily_fresh/pkg/common/version/verflag"
	"github.com/WeiXinao/daily_fresh/pkg/errors"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"

	"github.com/WeiXinao/daily_fresh/pkg/log"
)

var (
	progressMessage = color.GreenString("==>")
	//nolint: deadcode,unused,varcheck
	usageTemplate = fmt.Sprintf(`%s{{if .Runnable}}
  %s{{end}}{{if .HasAvailableSubCommands}}
  %s{{end}}{{if gt (len .Aliases) 0}}

%s
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

%s
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

%s{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  %s {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

%s
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

%s
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

%s{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "%s --help" for more information about a command.{{end}}
`,
		color.CyanString("Usage:"),
		color.GreenString("{{.UseLine}}"),
		color.GreenString("{{.CommandPath}} [command]"),
		color.CyanString("Aliases:"),
		color.CyanString("Examples:"),
		color.CyanString("Available Commands:"),
		color.GreenString("{{rpad .Name .NamePadding }}"),
		color.CyanString("Flags:"),
		color.CyanString("Global Flags:"),
		color.CyanString("Additional help topics:"),
		color.GreenString("{{.CommandPath}} [command]"),
	)
)

// App is the main structure of a cli application.
// It is recommended that an app be created with the app.NewApp() function.
type App[T CliOptions] struct {
	basename           string
	name               string
	description        string
	options            T
	runFunc            RunFunc
	stopFunc           StopFunc
	silence            bool
	noVersion          bool
	noConfig           bool
	subscriberInitFunc func(T) (subscriber.Subscriber, error)
	cfgr               configurator.Configurator[T]
	commands           []*Command
	args               cobra.PositionalArgs
	cmd                *cobra.Command
}

// Option defines optional parameters for initializing the application
// structure.
type Option[T CliOptions] func(*App[T])

// WithOptions to open the application's function to read from the command line
// or read parameters from the configuration file.
func WithOptions[T CliOptions](opt T) Option[T] {
	return func(a *App[T]) {
		a.options = opt
	}
}

// RunFunc defines the application's startup callback function.
type RunFunc func(basename string) error

// StopFunc defines the application's shutdown callback function.
type StopFunc func() error

// WithRunFunc is used to set the application startup callback function option.
func WithRunFunc[T CliOptions](run RunFunc) Option[T] {
	return func(a *App[T]) {
		a.runFunc = run
	}
}

func WithStopFunc[T CliOptions](stop StopFunc) Option[T] {
	return func(a *App[T]) {
		a.stopFunc = stop
	}
}

// WithDescription is used to set the description of the application.
func WithDescription[T CliOptions](desc string) Option[T] {
	return func(a *App[T]) {
		a.description = desc
	}
}

// WithSilence sets the application to silent mode, in which the program startup
// information, configuration information, and version information are not
// printed in the console.
func WithSilence[T CliOptions]() Option[T] {
	return func(a *App[T]) {
		a.silence = true
	}
}

// WithNoVersion set the application does not provide version flag.
func WithNoVersion[T CliOptions]() Option[T] {
	return func(a *App[T]) {
		a.noVersion = true
	}
}

// WithNoConfig set the application does not provide config flag.
func WithNoConfig[T CliOptions]() Option[T] {
	return func(a *App[T]) {
		a.noConfig = true
	}
}

// WithValidArgs set the validation function to valid non-flag arguments.
func WithValidArgs[T CliOptions](args cobra.PositionalArgs) Option[T] {
	return func(a *App[T]) {
		a.args = args
	}
}

// WithDefaultValidArgs set default validation function to valid non-flag arguments.
func WithDefaultValidArgs[T CliOptions]() Option[T] {
	return func(a *App[T]) {
		a.args = func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}

			return nil
		}
	}
}

func WithSubscribeInitFunc[T CliOptions](f func(cfg T) (subscriber.Subscriber, error)) Option[T] {
	return func(a *App[T]) {
		a.subscriberInitFunc = f
	}
}

// NewApp creates a new application instance based on the given application name,
// binary name, and other options.
func NewApp[T CliOptions](name string, basename string, opts ...Option[T]) *App[T] {
	a := &App[T]{
		name:     name,
		basename: basename,
	}

	for _, o := range opts {
		o(a)
	}

	a.buildCommand()

	return a
}

func (a *App[T]) buildCommand() {
	cmd := cobra.Command{
		Use:   FormatBaseName(a.basename),
		Short: a.name,
		Long:  a.description,
		// stop printing usage when the command errors
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          a.args,
	}
	// cmd.SetUsageTemplate(usageTemplate)
	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.Flags().SortFlags = true
	cliflag.InitFlags(cmd.Flags())

	if len(a.commands) > 0 {
		for _, command := range a.commands {
			cmd.AddCommand(command.cobraCommand())
		}
		cmd.SetHelpCommand(helpCommand(a.name))
	}
	if a.runFunc != nil {
		cmd.RunE = a.runCommand
	}

	var namedFlagSets cliflag.NamedFlagSets
	if !reflect.ValueOf(a.options).IsZero() {
		namedFlagSets = a.options.Flags()
		fs := cmd.Flags()
		for _, f := range namedFlagSets.FlagSets {
			fs.AddFlagSet(f)
		}

		usageFmt := "Usage:\n  %s\n"
		cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
		cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
			cliflag.PrintSections(cmd.OutOrStdout(), namedFlagSets, cols)
		})
		cmd.SetUsageFunc(func(cmd *cobra.Command) error {
			fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
			cliflag.PrintSections(cmd.OutOrStderr(), namedFlagSets, cols)

			return nil
		})
	}

	if !a.noVersion {
		verflag.AddFlags(namedFlagSets.FlagSet("global"))
	}

	if !a.noConfig {
		addConfigFlag(a.basename, namedFlagSets.FlagSet("global"))
	}

	globalflag.AddGlobalFlags(namedFlagSets.FlagSet("global"), cmd.Name())

	a.cmd = &cmd
}

func (a *App[T]) setUpConfigurator() error {
	var err error
	if a.subscriberInitFunc == nil {
		return nil
	}
	subscriber, err := a.subscriberInitFunc(a.options)
	if err != nil {
		return err
	}

	a.cfgr, err = configurator.NewConfigCenter[T](configurator.Config{
		Type: "yaml",
		Log:  true,
	}, subscriber)
	if err != nil {
		return err
	}
	return nil
}

func (a *App[T]) refreshConfig() error {
	cfg, err := a.cfgr.GetConfigString()
	if err != nil {
		return err
	}

	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBufferString(cfg))
	if err := viper.Unmarshal(a.options); err != nil {
		return err
	}

	return nil
}

func (a *App[T]) addConfigListener(f func (key string, raw string, data T)) {
	a.cfgr.AddListener(f)
}

// Run is used to launch the application.
func (a *App[T]) Run() {
	if err := a.cmd.Execute(); err != nil {
		fmt.Printf("%v %v\n", color.RedString("Error:"), err)
		os.Exit(1)
	}
}

// Command returns cobra command instance inside the application.
func (a *App[T]) Command() *cobra.Command {
	return a.cmd
}

func (a *App[T]) runCommand(cmd *cobra.Command, args []string) error {
	printWorkingDir()
	cliflag.PrintFlags(cmd.Flags())
	if !a.noVersion {
		// display application version information
		verflag.PrintAndExitIfRequested()
	}

	if !a.noConfig {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}

		if err := viper.Unmarshal(a.options); err != nil {
			return err
		}
	}

	err := a.setUpConfigurator()
	if err != nil {
		return nil
	}
	
	err = a.refreshConfig()
	if err != nil {
		return err
	}


	if !a.silence {
		log.Infof("%v Starting %s ...", progressMessage, a.name)
		if !a.noVersion {
			log.Infof("%v Version: `%s`", progressMessage, version.Get().ToJSON())
		}
		if !a.noConfig {
			log.Infof("%v Config file used: `%s`", progressMessage, viper.ConfigFileUsed())
		}
	}
	if !reflect.ValueOf(a.options).IsZero() {
		if err := a.applyOptionRules(); err != nil {
			return err
		}
	}

	eg := errgroup.Group{}

	a.addConfigListener(func(key, raw string, data T) {
		eg.Go(func() error {
			if a.stopFunc != nil {
				a.stopFunc()
			}

			err := a.refreshConfig(); 
			if err != nil {
				return err
			}

			if a.runFunc != nil {
				return a.runFunc(a.basename)
			}
			return nil
		})
	})

	eg.Go(func() error {
		// run application
		if a.runFunc != nil {
			return a.runFunc(a.basename)
		}
		return nil
	})

	return eg.Wait()
}

func (a *App[T]) applyOptionRules() error {
	if completeableOptions, ok := any(a.options).(CompleteableOptions); ok {
		if err := completeableOptions.Complete(); err != nil {
			return err
		}
	}

	if errs := a.options.Validate(); len(errs) != 0 {
		return errors.NewAggregate(errs)
	}

	if printableOptions, ok := any(a.options).(PrintableOptions); ok && !a.silence {
		log.Infof("%v Config: `%s`", progressMessage, printableOptions.String())
	}

	return nil
}

func printWorkingDir() {
	wd, _ := os.Getwd()
	log.Infof("%v WorkingDir: %s", progressMessage, wd)
}
