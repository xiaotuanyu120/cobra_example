## cobra_example

### Overview
Cobra是一个库，提供了一个简单的界面来创建类似于git＆go工具的强大的现代CLI界面。Cobra也是一个应用程序，它可以生成应用程序骨架，以快速开发基于Cobra的应用程序。

cobra提供:
- Easy subcommand-based CLIs: app server, app fetch, etc.
- Fully POSIX-compliant flags (including short & long versions)
- Nested subcommands
- Global, local and cascading flags
- Easy generation of applications & commands with cobra init appname & cobra add cmdname
- Intelligent suggestions (app srver... did you mean app server?)
- Automatic help generation for commands and flags
- Automatic help flag recognition of -h, --help, etc.
- Automatically generated bash autocomplete for your application
- Automatically generated man pages for your application
- Command aliases so you can change things without breaking them
- The flexibility to define your own help, usage, etc.
- Optional tight integration with viper for 12-factor apps

### Concepts
Cobra建立在commands，arguments和flag的结构上。

Commands代表动作，Args代表事物，Flags是这些动作的修饰符。

好的程序应该读起来像句子，用户会根据他们自然的理解就知道如何去使用这个命令。

语法遵循`APPNAME VERB NOUN --ADJECTIVE` 或 `APPNAME COMMAND ARG --FLAG`

举两个例子感受一下
```
hugo server --port=1313
git clone URL --bare
```
扩展资料：
- [cobra.command](https://godoc.org/github.com/spf13/cobra#Command)
- [flag package](https://golang.org/pkg/flag/)

### Installing
``` bash
# clone最新的cobra代码到开发环境的GOPATH中
go get -u github.com/spf13/cobra/cobra

# 在自己的代码里面import cobra
import "github.com/spf13/cobra"
```

### Getting Started
通常基于Cobra的应用程序将遵循以下组织结构：
```
  ▾ appName/
    ▾ cmd/
        add.go
        your.go
        commands.go
        here.go
      main.go
```
在cobra应用程序中，main.go的用途很简单，只有一个目标：初始化cobra
``` go
package main

import (
  "{pathToYourApp}/cmd"
)

func main() {
  cmd.Execute()
}
```

#### Using the Cobra Generator
Cobra提供了自己的程序，可以创建您的应用程序并添加您想要的任何命令。 这是将Cobra整合到您的应用程序中的最简单方法。

参看[这里](https://github.com/spf13/cobra/blob/master/cobra/README.md)来获取详细信息

#### Using the Cobra Library
要手动实现Cobra，您需要创建一个main.go文件和一个rootCmd文件。 您可以选择根据需要提供其他命令。

##### 创建rootCmd
cobra不需要任何特殊的构造函数。只需创建命令即可。

理想情况下，我们需要在`app/cmd/root.go`放入以下内容：
``` go
var rootCmd = &cobra.Command{
  Use:   "hugo",
  Short: "Hugo is a very fast static site generator",
  Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
  Run: func(cmd *cobra.Command, args []string) {
    // Do Stuff Here
  },
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
```
您还可以在init()函数中定义flag和设置配置相关内容。

例如`cmd/root.go`
``` go
import (
  "fmt"
  "os"

  homedir "github.com/mitchellh/go-homedir"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
)

func init() {
  cobra.OnInitialize(initConfig)
  rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
  rootCmd.PersistentFlags().StringVarP(&projectBase, "projectbase", "b", "", "base project directory eg. github.com/spf13/")
  rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "Author name for copyright attribution")
  rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "Name of license for the project (can provide `licensetext` in config)")
  rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
  viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
  viper.BindPFlag("projectbase", rootCmd.PersistentFlags().Lookup("projectbase"))
  viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
  viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
  viper.SetDefault("license", "apache")
}

func initConfig() {
  // Don't forget to read config either from cfgFile or from home directory!
  if cfgFile != "" {
    // Use config file from the flag.
    viper.SetConfigFile(cfgFile)
  } else {
    // Find home directory.
    home, err := homedir.Dir()
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    // Search config in home directory with name ".cobra" (without extension).
    viper.AddConfigPath(home)
    viper.SetConfigName(".cobra")
  }

  if err := viper.ReadInConfig(); err != nil {
    fmt.Println("Can't read config:", err)
    os.Exit(1)
  }
}
```

##### Create your main.go
使用rootCmd，你需要在main function中Excute它。为了清楚起见，应该在root上运行Execute，尽管可以在任何命令上调用它。

`main.go`中只有一个目标，就是初始化cobra
``` go
package main

import (
  "{pathToYourApp}/cmd"
)

func main() {
  cmd.Execute()
}
```

##### Create additional commands
可以定义其他命令，并且通常在`cmd/`目录中为每个命令提供自己的文件。

如果要创建版本命令，可以创建`cmd/version.go`并使用以下内容填充它：
``` go
package cmd

import (
  "fmt"

  "github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Print the version number of Hugo",
  Long:  `All software has versions. This is Hugo's`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
  },
}
```

#### Working with Flags
flag提供修饰符来控制动作命令的操作方式。

##### Assign flags to a command
由于flag是在不同的位置定义和使用的，因此我们需要在外部定义一个具有正确范围的变量来分配要使用的标志。
``` go
var Verbose bool
var Source string
```

分配flag有两种不同的方法。

##### Persistent Flags
flag可以是“persistent”，这意味着该flag可用于它所分配的命令以及该命令下的每个子命令。对于global flag，在根上分配flag作为persistent flag。
``` go
rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
```

##### Local Flags
还可以在本地分配一个flag，该flag仅适用于该特定命令。
``` go
rootCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
```

##### Local Flag on Parent Commands
默认情况下，Cobra仅解析目标命令上的本地标志，忽略父命令上的任何本地标志。通过启用Command.TraverseChildren，Cobra将在执行目标命令之前解析每个命令上的本地标志。
``` go
command := cobra.Command{
  Use: "print [OPTIONS] [COMMANDS]",
  TraverseChildren: true,
}
```

##### Bind Flags with Config
也可以通过viber来绑定flag
``` go
var author string

func init() {
  rootCmd.PersistentFlags().StringVar(&author, "author", "YOUR NAME", "Author name for copyright attribution")
  viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
}
```
在此示例中，persistent flag `author`与viper绑定。请注意，当用户未提供`--author` flag时，变量author将不会设置为config的值。

需要了解更多，参见[viber 文档](	https://github.com/spf13/viper#working-with-flags)

##### Required flags
flag默认是可选的。如果您希望命令在未设置flag时报告错误，请将其flag为必需：
``` go
rootCmd.Flags().StringVarP(&Region, "region", "r", "", "AWS region (required)")
rootCmd.MarkFlagRequired("region")
```

#### Positional and Custom Arguments
可以使用Command的Args字段指定位置参数的验证。

以下是内置的验证器：
- NoArgs - the command will report an error if there are any positional args.
- ArbitraryArgs - the command will accept any args.
- OnlyValidArgs - the command will report an error if there are any positional args that are not in the ValidArgs field of Command.
- MinimumNArgs(int) - the command will report an error if there are not at least N positional args.
- MaximumNArgs(int) - the command will report an error if there are more than N positional args.
- ExactArgs(int) - the command will report an error if there are not exactly N positional args.
- ExactValidArgs(int) = the command will report and error if there are not exactly N positional args OR if there are any positional args that are not in the ValidArgs field of Command
- RangeArgs(min, max) - the command will report an error if the number of args is not between the minimum and maximum number of expected args.

一个自定义验证器的例子：
``` go
var cmd = &cobra.Command{
  Short: "hello",
  Args: func(cmd *cobra.Command, args []string) error {
    if len(args) < 1 {
      return errors.New("requires at least one arg")
    }
    if myapp.IsValidColor(args[0]) {
      return nil
    }
    return fmt.Errorf("invalid color specified: %s", args[0])
  },
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("Hello, World!")
  },
}
```

#### Example
在下面的示例中，我们定义了三个命令。两个位于顶层，一个（cmdTimes）是顶级命令之一的子级。在这种情况下，root不可执行，这意味着需要子命令。这是通过不为'rootCmd'提供'Run'来实现的。

我们只为一个命令定义了一个flag。

有关标志的更多文档，请访问：[https://github.com/spf13/pflag](https://github.com/spf13/pflag)
``` go
package main

import (
  "fmt"
  "strings"

  "github.com/spf13/cobra"
)

func main() {
  var echoTimes int

  var cmdPrint = &cobra.Command{
    Use:   "print [string to print]",
    Short: "Print anything to the screen",
    Long: `print is for printing anything back to the screen.
For many years people have printed back to the screen.`,
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
      fmt.Println("Print: " + strings.Join(args, " "))
    },
  }

  var cmdEcho = &cobra.Command{
    Use:   "echo [string to echo]",
    Short: "Echo anything to the screen",
    Long: `echo is for echoing anything back.
Echo works a lot like print, except it has a child command.`,
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
      fmt.Println("Print: " + strings.Join(args, " "))
    },
  }

  var cmdTimes = &cobra.Command{
    Use:   "times [# times] [string to echo]",
    Short: "Echo anything to the screen more times",
    Long: `echo things multiple times back to the user by providing
a count and a string.`,
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
      for i := 0; i < echoTimes; i++ {
        fmt.Println("Echo: " + strings.Join(args, " "))
      }
    },
  }

  cmdTimes.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")

  var rootCmd = &cobra.Command{Use: "app"}
  rootCmd.AddCommand(cmdPrint, cmdEcho)
  cmdEcho.AddCommand(cmdTimes)
  rootCmd.Execute()
}
```
想了解一个完整的大型的应用程序，请查看[Hugo](http://gohugo.io/).

#### Help Command
当您有子命令时，Cobra会自动为您的应用程序添加一个帮助命令。当用户运行“app help”时会调用此方法。此外，帮助还将支持所有其他命令作为输入。比如说，你有一个名为'create'的命令，没有任何额外的配置;当'app help create'被调用时，Cobra会工作。每个命令都会自动添加' - help'标志。

帮助就像任何其他命令一样。它周围没有特殊的逻辑或行为。事实上，如果你愿意，你可以提供自己的。

##### Defining your own help
您可以提供自己的帮助命令或自己的模板，以使用以下函数使用的默认命令：
``` go
cmd.SetHelpCommand(cmd *Command)
cmd.SetHelpFunc(f func(*Command, []string))
cmd.SetHelpTemplate(s string)
```
后两个命令也支持任何子命令。

#### Usage Message
当用户提供了一个无效的flag或者无效的命令时，cobra会返回改应用的usage。

##### Example
您可以从上面的帮助中认识到这一点。那是因为默认帮助将用法嵌入其输出中。
``` bash
$ cobra --invalid
Error: unknown flag: --invalid
Usage:
  cobra [command]

Available Commands:
  add         Add a command to a Cobra Application
  help        Help about any command
  init        Initialize a Cobra Application

Flags:
  -a, --author string    author name for copyright attribution (default "YOUR NAME")
      --config string    config file (default is $HOME/.cobra.yaml)
  -h, --help             help for cobra
  -l, --license string   name of license for the project
      --viper            use Viper for configuration (default true)

Use "cobra [command] --help" for more information about a command.
```

##### Defining your own usage
您可以提供自己的usage功能或模板供Cobra使用。与帮助一样，函数和模板可以通过公共方法覆盖：

#### Version Flag
如果在root命令上设置了Version字段，Cobra会添加顶级'--version'标志。使用'--version'标志运行应用程序将使用版本模板将版本打印到stdout。可以使用`cmd.SetVersionTemplate(s string)`函数自定义模板。

#### PreRun and PostRun Hooks
可以在命令的主运行功能之前或之后运行功能。`PersistentPreRun`和`PreRun`函数将在Run之前执行。PersistentPostRun和PostRun将在Run之后执行。`Persistent * Run`函数将由子项继承，如果它们不声明它们自己的。这些功能按以下顺序运行：
- `PersistentPreRun`
- `PreRun`
- `Run`
- `PostRun`
- `PersistentPostRun`

下面是使用所有这些功能的两个命令的示例。执行子命令时，它将运行root命令的`PersistentPreRun`，但不运行root命令的`PersistentPostRun`：
``` go
package main

import (
  "fmt"

  "github.com/spf13/cobra"
)

func main() {

  var rootCmd = &cobra.Command{
    Use:   "root [sub]",
    Short: "My root command",
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
      fmt.Printf("Inside rootCmd PersistentPreRun with args: %v\n", args)
    },
    PreRun: func(cmd *cobra.Command, args []string) {
      fmt.Printf("Inside rootCmd PreRun with args: %v\n", args)
    },
    Run: func(cmd *cobra.Command, args []string) {
      fmt.Printf("Inside rootCmd Run with args: %v\n", args)
    },
    PostRun: func(cmd *cobra.Command, args []string) {
      fmt.Printf("Inside rootCmd PostRun with args: %v\n", args)
    },
    PersistentPostRun: func(cmd *cobra.Command, args []string) {
      fmt.Printf("Inside rootCmd PersistentPostRun with args: %v\n", args)
    },
  }

  var subCmd = &cobra.Command{
    Use:   "sub [no options!]",
    Short: "My subcommand",
    PreRun: func(cmd *cobra.Command, args []string) {
      fmt.Printf("Inside subCmd PreRun with args: %v\n", args)
    },
    Run: func(cmd *cobra.Command, args []string) {
      fmt.Printf("Inside subCmd Run with args: %v\n", args)
    },
    PostRun: func(cmd *cobra.Command, args []string) {
      fmt.Printf("Inside subCmd PostRun with args: %v\n", args)
    },
    PersistentPostRun: func(cmd *cobra.Command, args []string) {
      fmt.Printf("Inside subCmd PersistentPostRun with args: %v\n", args)
    },
  }

  rootCmd.AddCommand(subCmd)

  rootCmd.SetArgs([]string{""})
  rootCmd.Execute()
  fmt.Println()
  rootCmd.SetArgs([]string{"sub", "arg1", "arg2"})
  rootCmd.Execute()
}
```
output是
```
Inside rootCmd PersistentPreRun with args: []
Inside rootCmd PreRun with args: []
Inside rootCmd Run with args: []
Inside rootCmd PostRun with args: []
Inside rootCmd PersistentPostRun with args: []

Inside rootCmd PersistentPreRun with args: [arg1 arg2]
Inside subCmd PreRun with args: [arg1 arg2]
Inside subCmd Run with args: [arg1 arg2]
Inside subCmd PostRun with args: [arg1 arg2]
Inside subCmd PersistentPostRun with args: [arg1 arg2]
```
#### Suggestions when "unknown command" happens
当“未知命令”错误发生时，Cobra将打印自动建议。当错字发生时，这允许Cobra的行为类似于git命令。例如:
``` go
$ hugo srever
Error: unknown command "srever" for "hugo"

Did you mean this?
        server

Run 'hugo --help' for usage.
```
根据注册的每个子命令自动提出建议，使用了[Levenshtein distance](http://en.wikipedia.org/wiki/Levenshtein_distance)的实现。每个匹配最小距离为2（忽略大小写）的注册命令将显示为建议。

如果您需要在命令中禁用建议或调整字符串距离，请使用：
``` go
command.DisableSuggestions = true
```
or
``` go
command.SuggestionsMinimumDistance = 1
```

您还可以使用SuggestFor属性显式设置要为其指定命令的名称。 这允许建议字符串距离不接近的字符串，但在您的命令集和一些您不想要别名的字符串中有意义。 例：
``` bash
$ kubectl remove
Error: unknown command "remove" for "kubectl"

Did you mean this?
        delete

Run 'kubectl help' for usage.
```