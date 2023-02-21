<div align="center">
<h1 align="center">disgo-i18n</h1>
<p><a href="https://godoc.org/github.com/srvcontrol/disgo-i18n"><img src="https://godoc.org/github.com/srvcontrol/disgo-i18n?status.svg" alt="GoDoc"></a>
</div>

<div align="center">
<img src="https://img.shields.io/github/actions/workflow/status/srvcontrol/disgo-i18n/build.yml?style=for-the-badge" alt="Build Status">
<a href="https://goreportcard.com/report/github.com/srvcontrol/disgo-i18n"><img src="https://goreportcard.com/badge/github.com/srvcontrol/disgo-i18n?style=for-the-badge" alt="Report Card"></a> 
<a href="https://codecov.io/gh/srvcontrol/disgo-i18n"><img src="https://img.shields.io/codecov/c/github/srvcontrol/disgo-i18n?style=for-the-badge" alt="codecov"></a></p>
</div>

<p align="center">A simple and lightweight Go package that helps you translate Go programs into <a href="https://discord.com/developers/docs/reference#locales"> languages supported by Discord</a>.</p>

## Features

- Adapted from [kaysoro/discordgo-i18n](https://github.com/kaysoro/discordgo-i18n) for easy use with [disgo](https://github.com/disgoorg/disgo).
- Supports multiple strings per key to make your bot "more alive".
- Supports variables in strings using [text/template](http://golang.org/pkg/text/template/) syntax.
- Supports JSON language files.


## Installing

This assumes you already have a working Go environment, if not please see
[this page](https://golang.org/doc/install) first.

`go get` *will always pull the latest tagged release from the master branch.*

```sh
go get github.com/srvcontrol/disgo-i18n
```

## Usage

Import the package into your project.

```go
import i18n "github.com/srvcontrol/disgo-i18n"
```

Load a locale bundle.

```go
err := i18n.LoadBundle(discord.LocaleFrench, "path/to/your/file.json")
```

The bundle format must respect the schema below; note [text/template](http://golang.org/pkg/text/template/) syntax is used to inject variables.  
For a given key, value can be a:
- string
- string array to randomize translations
- deep structure to group translations as needed. 

If any other type is provided, it is mapped to string automatically.

```json
{
    "hello_world": "Hello world!",
    "hello_anyone": "Hello {{ .anyone }}!",
    "image": "https://media2.giphy.com/media/Ju7l5y9osyymQ/giphy.gif",
    "bye": ["See you", "Bye!"],
    "command": {
        "scream": {
            "name": "scream",
            "description": "Screams something",
            "dog": "Waf waf! 🐶",
            "cat": "Miaw! 🐱"
        }
    }
}
```

By default, the locale fallback used when a key does not have any translations is `discord.LocaleEnglishGB`. To change it, use the following method.

```go
i18n.SetDefault(discord.LocaleChineseCN)
```

To get translations use the below thread-safe method; if any translation cannot be found or an error occurred even with the fallback, key is returned.

```go
helloWorld := i18n.Get(discord.LocaleEnglishGB, "hello_world")
fmt.Println(helloWorld)
// Prints "Hello world!"

helloAnyone := i18n.Get(discord.LocaleEnglishGB, "hello_anyone")
fmt.Println(helloAnyone)
// Prints "Hello {{ .anyone }}!"

helloAnyone = i18n.Get(discord.LocaleEnglishGB, "hello_anyone", i18n.Vars{"anyone": "Nick"})
fmt.Println(helloAnyone)
// Prints "Hello Nick!"

bye := i18n.Get(discord.LocaleEnglishGB, "bye")
fmt.Println(bye)
// Prints randomly "See you" or "Bye!"

keyDoesNotExist := i18n.Get(discord.LocaleEnglishGB, "key_does_not_exist")
fmt.Println(keyDoesNotExist)
// Prints "key_does_not_exist"

dog := i18n.Get(discord.LocaleEnglishGB, "command.scream.dog")
fmt.Println(dog)
// Prints "Waf waf! 🐶"
```

To get localizations for a command name, description, options or other fields, use the below thread-safe method. It retrieves a `map[discord.Locale]string` based on the loaded bundles.

```go
screamCommand := discord.SlashCommandCreate{
    Name:                     "scream",
    Description:              "Here a default description for my command",
    NameLocalizations:        i18n.GetLocalizations("command.scream.name"),
    DescriptionLocalizations: i18n.GetLocalizations("command.scream.description"),
}
```

Here an example of how it can work with interactions.

```go
func HelloWorld(event *discord.EventsApplicationCommandInteractionCreate) {

    embed := discord.NewEmbedBuilder().
		SetTitle(i18n.Get(event.Locale(), "hello_world")).
		SetDescription(i18n.Get(event.Locale(), "hello_anyone",
			i18n.Vars{"anyone": event.Member().Nick})
		).SetImage(i18n.Get(event.Locale(), "image")).Build()
	
    err := event.CreateMessage(discord.NewMessageCreateBuilder().SetEmbeds(embed).Build())

    // ...
}
```




## License

disgo-i18n is available under the same [MIT license](LICENSE) as the original project.
