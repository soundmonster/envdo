# `envdo`

Like [`awsudo`](https://github.com/makethunder/awsudo) but for any set of env vars.

## Install

Download binary for your system from the [latest release](https://github.com/soundmonster/envdo/releases); put it somewhere on your path.

## Usage

### Setup

Create a file `~/.envdo.toml`; it should contain data in the following format:

```toml
[profile_name] # profile names should not nest; default profile is called `default`
# Env var names should be upper case with underscores
ENV_VAR_NAME = "value" # values should be TOML strings; no other data type is allowed
```

### Run

```sh
envdo -p <profile_name> <cmd with args>

# run "env" with profile "foo"
envdo -p foo env 

# run "./secret_script.py" and some arguments with implicit profile "default"
envdo ./secret_script.py --yes -R

# run a script under "sh" with profile "foo" 
envdo -p foo sh -c 'echo $FOO'
```

## Difference to the other `envdo`

I realized that another tool called [`envdo`](https://github.com/antonlindstrom/envdo) exists. It follows a different
philosophy: instead of being a per-command switcher like `sudo`, `antonlindstrom`'s `envdo` is more like `su` in that it
switches env var profiles for a session and makes them sticky.

For completeness' sake, there is also [`direnv`](https://direnv.net).
