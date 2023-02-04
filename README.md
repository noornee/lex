# LEX
[![build](https://github.com/cmd777/lex/actions/workflows/build_all_os.yml/badge.svg)](https://github.com/cmd777/lex/actions/workflows/build_all_os.yml)

# Table of Contents
**🇺🇸 English**

[**ℹ️ Information about the project**](#%E2%84%B9%EF%B8%8F-information-about-the-project)
- [What is Lex?](#what-is-lex)
  - [Why create Lex?](#why-create-lex)
- [How do I use Lex?](#how-do-i-use-lex)
  - [Is Lex Free to Use? Do I need to pay to use Lex?](#is-lex-free-to-use-do-i-need-to-pay-to-use-lex)
  
[**⚙️ Installation Instructions**](#%EF%B8%8F-installation-instructions)
- [Building Instructions](#building-instructions)
- [Downloading Binaries](#downloading-binaries)

[**🚩 Bugs, Issues, and other Important Information**](#-bugs-issues-and-other-important-information)

[**🔨 Technologies that were used to create Lex**](#-technologies-that-were-used-to-create-lex)

---

**[🇭🇺 Magyar README](https://github.com/cmd777/lex/blob/main/README_hu.md)**

---

# ℹ️ Information about the project

## What is Lex?
LEX (LazerEX) is a lightweight, open source frontend for reddit written in Go.

## Why create Lex?

Reddit can be very slow, and is very bloated, that's no secret. It was really annoying to either wait a super long time for a single post to load, or the UI freezing up for no apparent reason, these were just some of the issues that were common for me in every browser I tried (Chrome, Firefox, Brave, etc.) so, I decided to try to re-create it. I did try to re-create it previously, and named it Lazer, but that attempt failed miserably, and I decided to give it one last shot. And so LazerEX was created, which is a far superior version than the previous attempt.

Lex can use up to 60% less bandwidth, while keeping the same image/video quality

As for the UI, I went for the newer reddit redesign, but with a little tweaking.

## How do I use Lex?

Unfortunately, unlike other open source reddit frontends, **Lex does not provide a demo website**, and there is two ways to build/use Lex.

1. [Downloading a Binary and running it](#downloading-binaries)
2. [Building Lex yourself](#building-instructions)

## Is Lex Free to Use? Do I need to pay to use Lex?

**Lex is, and always will be open source, and free to use either commerically, and/or privately, for as long as the origin of the software is not misrepresented.**

For more information, refer to the [LICENSE](https://github.com/cmd777/lex/blob/main/LICENSE).

# ⚙️ Installation Instructions

## Building Instructions

The only requirement to build the project is [Go](https://go.dev/dl)

If Go is installed, run the following commands in your terminal
```cmd
git clone https://github.com/cmd777/lex.git
cd lex
cd src
go get -u
go build
```
and, you're pretty much done, all that's left to do is launch the built binary, and navigate to `localhost:9090/r/{subreddit}`

## Downloading Binaries

To install LEX via automatically built binaries, go to the [releases](https://github.com/cmd777/lex/releases/latest) tab, and download the appropriate zip for your OS + ARCH, and extract it.

- For 32 bit Windows machines, download `lex-windows.zip`, then launch lex-i386-windows.exe
- For 64 bit Windows machines, download `lex-windows.zip`, then launch lex-amd64-windows.exe

<br>

- For 32 bit Linux machines, download `lex-linux.zip`, then launch lex-i386-linux
- For 64 bit Linux machines, download `lex-linux.zip`, then launch lex-amd64-linux

<br>

- For Intel MacOS machines, download `lex-osx.zip`, then launch lex-amd64-osx
- For M1 MacOS machines, download `lex-osx.zip`, then launch lex-arm64-osx

After launching, navigate to `localhost:9090/r/{subreddit}`, and you're done.

# 🚩 Bugs, Issues, and other Important Information

Lex is still in very early stages of development, a lot of features are missing, and it might be prone to bugs, such as text overflowing, gallery buttons acting weird, and more.

I try to fix most of the critical bugs before pushing any changes, but if you find a bug, feel free to [create an issue](https://github.com/cmd777/lex/issues) for it.

I also work on Lex in my free time as a hobby, so development may be slow, thank you for your patience!

# 🔨 Technologies that were used to create Lex

- [Go](https://go.dev/) ➡️ Powering just About Everything
- [Humanize (go-humanize)](https://github.com/dustin/go-humanize) ➡️ Formatting to Human Friendly Units 
- [Gzip (gin-contrib)](https://github.com/gin-contrib/gzip) ➡️ Gin Middleware to Enable Gzip Support
- [Gin](https://github.com/gin-gonic/gin) ➡️ HTTP Web Framework
- [Bluemonday](https://github.com/microcosm-cc/bluemonday) ➡️ HTML Sanitizer
- [Blackfriday](https://github.com/russross/blackfriday/tree/v2) ➡️ Markdown Processor
- [SASS](https://sass-lang.com) ➡️ CSS Extension

<sub>Also used to create LEX:</sub>

- [Josefin Sans](https://fonts.google.com/specimen/Josefin+Sans) ➡️ Navbar, index font
- [Open Sans](https://fonts.google.com/specimen/Open+Sans) ➡️ Subreddit font
- [SVGRepo](https://www.svgrepo.com) ➡️ SVGs