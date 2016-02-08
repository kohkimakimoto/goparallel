# Goparallel

Execute commands in parallel.

## Installation

Goparallel is provided as a single binary. You can download it and drop it in your $PATH.

[Download latest version](https://github.com/kohkimakimoto/goparallel/releases/latest)

## Usage

Pass commands list that are executed in parallel to `goparallel` by the stdin.

Example1) From a file using `cat` command.

```
cat commands.ltsv
cmd:echo one
cmd:echo two

cat commands.ltsv | goparallel
one
two
```

Example2) Use echo.

```
echo -e "cmd:echo one\ncmd:echo two" | goparallel
one
two
```

Example3) You can use a argument instead of stdin.

```
$ goparallel "cmd:echo one
dquote> cmd:echo two"
one
two
```

The commands list is defined in a LTSV format at default.
Above examples use `cmd` key that defines command to execute.

You can use `prefix` key to output with a prefix.

```
echo -e "cmd:echo one\tprefix:[aaa]\ncmd:echo two\tprefix:[bbb]" | goparallel
[aaa] one
[bbb] two
```

Goparallel supports other formats. use `-f` option.

YAML format.
```
$ cat commands.yml
- {cmd: "echo one", prefix: "[aaa]"}
- {cmd: "echo two", prefix: "[bbb]"}

$ cat commands.yml  | goparallel -f=yaml
[aaa] one
[bbb] two
```

JSON format.
```
$ cat commands.json
[
  {"cmd":"echo one", "prefix":"[aaa]"},
  {"cmd":"echo two", "prefix":"[bbb]"},
]

$ cat commands.json | goparallel -f=json
[aaa] one
[bbb] two
```

## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

The MIT License (MIT)
