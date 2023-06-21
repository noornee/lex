# LEX <br> <a href="https://github.com/cmd777/lex/actions/workflows/build_all_os.yml">![](https://img.shields.io/github/actions/workflow/status/cmd777/lex/build_all_os.yml?logo=github&style=flat-square)</a> <a href="https://github.com/cmd777/lex/blob/main/LICENSE">![](https://img.shields.io/github/license/cmd777/lex?logo=opensourceinitiative&style=flat-square)</a>

## Table of Contents

[**ℹ️ Information about the project**](#ℹ️-information-about-the-project)
- [What is LEX? / Issues with Reddit](#what-is-lex--issues-with-reddit)
- [How does LEX Fix these Issues?](#how-does-lex-fix-these-issues)
- [How do I use LEX?](#how-do-i-use-lex)
  
[**⚙️ Installation Instructions**](#%EF%B8%8F-installation-instructions)
- [Building Instructions](#building-instructions)
- [Downloading Binaries](#downloading-binaries)

[**🔬 Compatibility**](#-compatibility)

[**🚩 Bugs, Issues, and other Important Information**](#-bugs-issues-and-other-important-information)

[**📜 Legal Disclaimer**](#-legal-disclaimer)

[**🧰 Technologies that were used to create LEX**](#-technologies-that-were-used-to-create-lex)

---

# ℹ️ Information about the project

## What is LEX? / Issues with Reddit
**LEX** (**L**azer**EX**) is a Lightweight, Open Source Frontend for Reddit written in Go.<br>
The name "Lazer" is an intentional misspelling of the word "Laser"

## Why?
Reddit has a lot of potential, but is purposely degraded.<br>
Why? It was likely done to push their mobile app, which is kind of saddening.

The issues with Reddit currently are some of the following:
- NSFW is heavily blocked (even if the post is not actually NSFW)
  - [On mobile, Reddit won't even allow you to view ANY subreddit. Be it NSFW or Not.](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/fullexperience.png)
  <br><sub>You can even test this, by setting your user agent to something like `Mozilla/5.0 (Android 13; Mobile; rv:68.0) Gecko/68.0 Firefox/112.0`</sub>
- The UI is extremely slow, and unresponsive.
  - To the point that it can take 20 seconds just to go back from a comment section.
- Too much bloat, and therefore the data usage is very high.

These issues can be observed on all browsers (for example: Chrome, Firefox, Brave, etc.)

These are some of the issues that LEX aims to fix.

## How does LEX Fix these Issues?

- NSFW content is allowed, but only if you allow it.
  - In the settings page (located at `http://localhost:9090/config` there is an option to enable NSFW subreddits, and posts.)
- Optional Javascript and Configurations.
  - A lot of configurations, so you can use LEX as you want to!
  <br>Configurations include:
    - Enabling HLS Videos
    - Enabling Infinite "Scroll"
    - Enabling Gallery Navigation via Arrow Keys
    - Allowing NSFW Subreddits and Posts
    - Allowing Images from Unknown Sources
    - Using Advanced Math to Re-size the Page
    - Setting Preferred Image Resolutions
- No Ads.

## How do I use LEX?

Unfortunately, unlike other open source reddit frontends, **LEX does not provide a demo website**, and there are three ways to build/use LEX.

1. [Downloading a pre-built binary](#downloading-binaries)
2. [Building LEX yourself](#building-instructions)
3. [Using go install](#go-install)

# ⚙️ Installation Instructions

## Building Instructions

The only requirement to build the project is [Go](https://go.dev/dl)

If Go is installed, run the following commands in your terminal
```shell
git clone https://github.com/cmd777/lex.git &&
cd lex/cmd/lex &&
go get -u &&
go build
```
That's everything! All you need to do next, is just navigate to `http://localhost:9090/r/{your_favorite_subreddit}`

## Downloading Binaries

LEX provides downloadable, pre-built binaries for all major operating systems, including Windows, Linux, and MacOS.

To install LEX via these automatically built binaries, go to the [releases](https://github.com/cmd777/lex/releases/latest) tab, and download the appropriate zip for your OS + ARCH, and extract it.

- For 32 bit Windows machines, download `lex-windows.zip`, then launch lex-i386-windows.exe
- For 64 bit Windows machines, download `lex-windows.zip`, then launch lex-amd64-windows.exe

<br>

- For 32 bit Linux machines, download `lex-linux.zip`, then launch lex-i386-linux
- For 64 bit Linux machines, download `lex-linux.zip`, then launch lex-amd64-linux

<br>

- For Intel MacOS machines, download `lex-osx.zip`, then launch lex-amd64-osx
- For M1 MacOS machines, download `lex-osx.zip`, then launch lex-arm64-osx

That's everything! All you need to do next, is just navigate to `http://localhost:9090/r/{your_favorite_subreddit}`

<sub>* Pre-built binaries are published every time there is a push to the main branch.</sub>

## Go Install

LEX can also be installed via the `go install` command.

If you have at least `go 1.16 or later`, you can run the following command:
```shell
go install github.com/cmd777/lex/cmd/lex@latest
```

which will install the latest version of LEX.

Afterwards, simply run `lex` in your terminal, and that's everything!

> **Note** you may need to add `export PATH=$PATH:$(go env GOPATH)/bin` to your `.profile` file to be able to use the command `lex`!

If you wish to uninstall LEX, simply run the following command:
```shell
rm $(go env GOPATH)/bin/lex
```

# 🔬 Compatibility

<figcaption>No JavaScript</figcaption>

| ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/chrome.svg) Chrome | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/edge.svg) Edge | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/safari.svg) Safari | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/firefox.svg) Firefox | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/opera.svg) Opera | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/iexplorer.svg) IE |
|:------:|:----:|:------:|:-------:|:-----:|:----------------:|
| 26>    | 19>  | 6.1>   | 49>     | 12.1> | 11<sup>(?)</sup> |

<details>
<summary>INFScroll Setting Enabled</summary>

| ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/chrome.svg) Chrome | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/edge.svg) Edge | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/safari.svg) Safari | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/firefox.svg) Firefox | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/opera.svg) Opera | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/iexplorer.svg) IE |
|:------:|:----:|:------:|:-------:|:-----:|:--:|
| 42>    | -    | 10.1>  | -       | 29>   | X  |

</details>

<details>
<summary>GalleryNav Setting Enabled</summary>

| ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/chrome.svg) Chrome | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/edge.svg) Edge | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/safari.svg) Safari | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/firefox.svg) Firefox | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/opera.svg) Opera | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/iexplorer.svg) IE |
|:------:|:----:|:------:|:-------:|:-----:|:--:|
| 45>    | -    | 9>     | -       | 32>   | X  |

</details>

<details>
<summary>HLS Setting Enabled</summary>

| ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/chrome.svg) Chrome | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/edge.svg) Edge | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/safari.svg) Safari | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/firefox.svg) Firefox | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/opera.svg) Opera | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/iexplorer.svg) IE |
|:------:|:----:|:------:|:-------:|:-----:|:--:|
| 41>    | ?    | 9>     | -       | ?     | X  |

</details>

<details>
<summary>Advanced Math Enabled</summary>

| ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/chrome.svg) Chrome | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/edge.svg) Edge | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/safari.svg) Safari | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/firefox.svg) Firefox | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/opera.svg) Opera | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/iexplorer.svg) IE |
|:------:|:----:|:------:|:-------:|:-----:|:--:|
| 79>    | 79>  | 13.1>  | 75>     | 66>   | X  |

</details>

<details>
<summary>Legend</summary>

`>` -> Any equal, or newer version is supported

`?` -> Unsure of compatability

`X` -> Not compatible

`-` -> No newer version required

</details>

<sub>The table displays the recommended versions for browsers, though older versions *may* be compatible to an extent.</sub>

<sub>Based on [Can I Use](https://caniuse.com) data</sub>

---

| ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/windows.svg) Windows | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/linux.svg) Linux | ![](https://raw.githubusercontent.com/cmd777/lex/main/docs/images/macos.svg) macOS |
| :-----: | :---: | :---: |
| Windows 7 and higher | Kernel version 2.6.32 or later | macOS High Sierra 10.13 or newer

<sub>The table displays the **required** minimum versions for operating systems.</sub>

<sub>**LEX does not require administrator, or firewall privileges.**</sub>

<sub>Based on the [Go Documentation](https://github.com/golang/go/wiki/MinimumRequirements)</sub>

# 🚩 Bugs, Issues, and other Important Information

LEX is still in early stages of development, a lot of features are missing, and it might be prone to bugs.

I try to fix most of the critical bugs before pushing any changes, but if you find a bug, or have any questions, feel free to [create an issue](https://github.com/cmd777/lex/issues) for it.

For the features that are planned to be added, or the things that need to be fixed, you can once again, take a look at the [issues tab](https://github.com/cmd777/lex/issues) 

I work on LEX in my free time as a hobby, so development may be slow, thank you for your patience!

# 📜 Legal Disclaimer

LEX is not affiliated with, sponsored, or endorsed by Reddit.

All content that is displayed on LEX has been sourced from Reddit. LEX does not host any of the content

In case of any issues with a post, such as copyright infringement, trademark infringement, or violation of Reddit's community rules, the reports should be directed to Reddit.

# 🧰 Technologies that were used to create LEX

<details>
  <summary><a href="https://go.dev">Go</a> ➡️ Programming Language</summary>
  https://github.com/golang/go/blob/master/LICENSE
</details>

<details>
  <summary><a href="https://github.com/dustin/go-humanize">Humanize (go-humanize)</a> ➡️ Formatting time, numbers, etc.. to Human Friendly Units</summary>
  https://github.com/dustin/go-humanize/blob/master/LICENSE
</details>

<details>
  <summary><a href="https://github.com/gofiber/fiber">Fiber</a> ➡️ HTTP Web Framework</summary>
  https://github.com/gofiber/fiber/blob/master/LICENSE
</details>

<details>
  <summary><a href="https://github.com/gofiber/utils">(Fiber) Utils</a> ➡️ Common functions with better performance</summary>
  https://github.com/gofiber/utils/blob/master/LICENSE
</details>

<details>
  <summary><a href="https://github.com/goccy/go-json">Go-JSON</a> ➡️ Fast JSON Decoder</summary>
  https://github.com/goccy/go-json/blob/master/LICENSE
</details>

<details>
  <summary><a href="https://github.com/microcosm-cc/bluemonday">Bluemonday</a> ➡️ HTML Sanitizer</summary>
  https://github.com/microcosm-cc/bluemonday/blob/main/LICENSE.md
</details>

<details>
  <summary><a href="https://github.com/russross/blackfriday/tree/v2">Blackfriday</a>  ➡️ Markdown Processor</summary>
  https://github.com/russross/blackfriday/blob/master/LICENSE.txt
</details>

<details>
  <summary><a href="https://github.com/sass/sass">SASS</a> ➡️ CSS Extension</summary>
  https://github.com/sass/sass/blob/main/LICENSE
</details>

<sub>Also used to create LEX:</sub>

<details>
  <summary><a href="https://github.com/googlefonts/josefinsans">Josefin Sans</a> ➡️ Navbar, index font</summary>
  https://github.com/googlefonts/josefinsans/blob/master/OFL.txt
</details>

<details>
  <summary><a href="https://github.com/googlefonts/opensans">Open Sans</a> ➡️ Subreddit font</summary>
  https://github.com/googlefonts/opensans/blob/main/OFL.txt
</details>

<details>
  <summary><a href="https://www.svgrepo.com">SVGRepo</a> ➡️ SVGs</summary>
  https://www.svgrepo.com/page/licensing
</details>

<details>
  <summary><a href="https://github.com/Yqnn/svg-path-editor">Yqnn's SVG Path Editor</a> ➡️ Was used to edit almost all SVGs</summary>
  https://github.com/Yqnn/svg-path-editor/blob/master/LICENSE
</details>
