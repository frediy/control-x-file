# control-x-file
CLI clipboard tool for cutting and pasting recursive paths or files with original path restoration.

## Usage
```
Usage: cx [<options>] [<path>...]

  Detectes whether [<path>...] is in current dir or clipboard.
  Cuts [<path>...] in current dir recursively to clipboard.
  Pastes [<path>...] recursively from clipboard to current dir.

options:
  -a	paste all clipboard paths into current dir
```

## Installation
```
  # clone
  git clone https://github.com/frediy/control-x-file

  # build
  cd control-x-file
  go build -o "bin/cx" .

  # update path in .bashrc or .zshrc
  echo "export PATH=\$PATH:$(pwd)/bin/" >> ~/.zshrc
```