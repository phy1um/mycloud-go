# MyCloud GO

_Human scale file management and sharing_

Very WIP. Created by Tom Marks

My Links: [website](https://coding.tommarks.xyz) | [youtube](https://www.youtube.com/c/TomMarksTalksCode) | [twitch](https://twitch.tv/phylum919)

## What

* Single-binary distribution
* Upload files over SCP
* Manage files over a text interface over SSH
* Generate keys to share files with HTTP links
* Login solved by using existing public key authentication with SSH
* Tag files for easier searching, state is stored in a sqlite database


## Why

I wanted a solution to manage files and sometimes share them with people or download them easily.
Simple text interfaces are cool, so this was an experiment rather than building a webapp.
Doing all file management over SSH meant that some problems were solved out the gate.

## Status

- [x] Simple file management over SSH (tag files, make access keys)
- [x] Download files using access keys over HTTP
- [x] Pretty thorough logging using `zerolog` (except for a few log.Println statements left)
- [ ] Make it not look ugly (experimenting with [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss))
- [ ] Paginated file view (currently limited to 10 files :D)
- [ ] More options for HTTP access keys
