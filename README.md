[![Go Report Card](https://goreportcard.com/badge/github.com/frediy/control-x-file)](https://goreportcard.com/report/github.com/frediy/control-x-file)

# control-x-file
CLI clipboard tool for cutting and pasting paths and subdirs with exact restoration.

No need to remember paths that no longer exist when calling `git restore`. Avoid accidentally combining paths incorrectly  with `cp -r` due to a rogue "/".


**Use case: Restore hard to remember files fast**
```bash
# cut
~/my-repo $ git checkout branch-with-good-stuff
~/my-repo $ cx src/deep/into/the/module/space/SomeObscureClassWithLongName.code
~/my-repo $ cx src/deep/into/another/part/of/the/module/space/SomeEvenMoreObscureClassWithLongName.code

~/my-repo $ git checkout feature-branch
~/my-repo $ cx -a

~/my-repo $ ls src/deep/into/the/module/space/SomeObscureClassWithLongName.code
src/deep/into/the/module/space/SomeObscureClassWithLongName.code

~/my-repo $ ls src/deep/into/another/part/of/the/module/space/SomeEvenMoreObscureClassWithLongName.code
src/deep/into/another/part/of/the/module/space/SomeEvenMoreObscureClassWithLongName.code
```


**Use case: Quick reset**
```bash
~/my-repo $ cat some/module_with_many_files/file1.code
def clean_awesome_code():
  return yeah
~/my-repo $ cat some/module_with_many_files/file2.code
def best_code():
  return so_good

# cut
~/my-repo $ cx some/module_with_many_files
~/my-repo $ git checkout .

# *make changes*..

~/my-repo $ cat some/module_with_many_files/file1.code
def dirty_sad_code():
  return no
~/my-repo $ cat some/module_with_many_files/file2.code
def mediocre_code():
  return so_medium

# restore
~/my-repo $ cx some/module_with_many_files

~/my-repo $ cat some/module_with_many_files/file1.code
def clean_awesome_code():
  return yeah
~/my-repo $ cat some/module_with_many_files/file2.code
def best_code():
  return so_good
```

**Use case: Compose new repo**
```bash
~/repo1 $ cx src/utilities
~/repo2 $ cx src/utilities/other
~/repo3 $ cx Dockerfile
~/repo4 $ cx bash_scripts/deployment

~/new-repo $ cx -a
~/new-repo $ tree .
.
└── src
|   └── utilities
|       ├── subdir
|       │   └── repo1a.sh
|       ├── repo1b.sh
|       └── other
|           ├── repo2a.sh
|           └── repo2b.sh
├── bash_scripts
│   └── deployment
│       └── repo4.sh
└── Dockerfile
```

## Usage
```
Usage: cx [<options>] [<path>...]

  Detectes whether [<path>...] is in current dir or clipboard.
  Cuts [<path>...] in current dir recursively to clipboard.
  Pastes [<path>...] recursively from clipboard to current dir.

options:
  -a, --all    paste all clipboard paths into current dir
  -h, --help   show help
  -k, --keep   keep paths in workdir after cut and in clipboard after paste
  -l, --list   list all paths in clipboard
```

## Installation
```
  # clone
  git clone https://github.com/frediy/control-x-file

  # build
  cd control-x-file
  go mod tidy
  go build -o "bin/cx" .

  # update path in .bashrc or .zshrc
  echo "export PATH=\$PATH:$(pwd)/bin/" >> ~/.zshrc
```