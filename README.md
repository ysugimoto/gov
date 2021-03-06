# gov
Golang versioning tool

## Installation

```
$ go get -u github.com/ysugimoto/gov...
```

You can use `gov` command.

## Setup

`gov` follows [semver](https://semver.org/) as `v(major.minor.patch)` versioning.

### Initialize version

```
$ gov init
```

This command create `.versions` file at current working directory (as project root).

### Show current version

```
$ gov
```

`gov` command will find up `.versions` file and use it. So you can run `gov` on sub directories.

## Bump versions

Once you execute following commands, the `gov` will make new commit and version tag.
Note that versioning should do on `master` branch, `gov` command makes sure you are in `master` branch.

### patch

```
# from v0.0.1 to v0.0.2
$ gov patch
>> v0.0.2
```

### minor

```
# from v0.0.1 to v0.1.0
$ gov minor
>> v0.1.0
```
### major

```
# from v0.0.1 to v1.0.0
$ gov major
>> v1.0.0
```

After that, you can push to remote with new commit and tag:

```
$ git push --follow-tags
```

## Author

Yoshiaki Sugimoto

## License

MIT
