# LEX
[![build](https://github.com/cmd777/lex/actions/workflows/build_all_os.yml/badge.svg)](https://github.com/cmd777/lex/actions/workflows/build_all_os.yml)

# Tartalomjegyzék

**🇭🇺 Magyar**

[**ℹ️ Információk a projektről**](#%E2%84%B9%EF%B8%8F-inform%C3%A1ci%C3%B3k-a-projektr%C5%91l)
- [Mi az a Lex?](#mi-az-a-lex)
  - [Miért hoztad létre a Lexet?](#mi%C3%A9rt-hoztad-l%C3%A9tre-a-lexet)
- [Hogyan használom a Lexet?](#hogyan-kell-haszn%C3%A1lni-a-lexet)
  - [Ingyenesen használható a Lex? Fizetnem kell a Lex használatáért?](#a-lex-ingyenesen-haszn%C3%A1lhat%C3%B3-fizetnem-kell-a-lex-haszn%C3%A1lat%C3%A1%C3%A9rt)
  
[**⚙️ Telepítési útmutató**](#%EF%B8%8F-telep%C3%ADt%C3%A9si-%C3%BAtmutat%C3%B3)
- [Építési útmutató](#%EF%B8%8F-telep%C3%ADt%C3%A9si-%C3%BAtmutat%C3%B3)
- [Binárisok letöltése](#bin%C3%A1ris-f%C3%A1jlok-let%C3%B6lt%C3%A9se)

[**🚩 Hibák, problémák és egyéb fontos információk**](#-hib%C3%A1k-probl%C3%A9m%C3%A1k-%C3%A9s-egy%C3%A9b-fontos-inform%C3%A1ci%C3%B3k)

[**🔨 A Lex létrehozásához használt technológiák**](#-a-lex-l%C3%A9trehoz%C3%A1s%C3%A1hoz-haszn%C3%A1lt-technol%C3%B3gi%C3%A1k)

---

# ℹ️ Információk a projektről

## Mi az a Lex?
A LEX (LazerEX) egy könnyű, nyílt forráskódú kezelőfelület a reddithez, amit Go-ban készítettem.

## Miért hoztad létre a Lexet?

A Reddit nagyon lassú tud lenni, és rengeteg internetet használ, ez nem titok. Nagyon bosszantó volt, hogy vagy nagyon sokáig vártam, amíg egy poszt betöltődik, vagy a felhasználói felület nyilvánvaló ok nélkül lefagyott. Ez csak néhány olyan probléma, amit én néhány böngészőben (Chrome, Firefox, Brave stb.) tapasztaltam, ezért úgy döntöttem, hogy megpróbálom újra csinálni. Korábban megpróbáltam újra csinálni a redditet, és Lazernek neveztem el, de az a próbálkozás nem sikerült, és úgy döntöttem, adok neki egy utolsó esélyt. Így jött létre a LazerEX, ami egy sokkal jobb verzió, mint az előző próbálkozás.

A Lex akár 60%-kal kevesebb internetet tud használni, miközben megőrzi ugyanazt a kép/videó minőséget

Ami a felhasználói felületet illeti, az újabb reddit újratervezést választottam, de egy kis finomítással.

## Hogyan kell használni a Lexet?

Sajnos, ellentétben más nyílt forráskódú reddit frontendekkel, a **Lex nem biztosít bemutató webhelyet**, és kétféleképpen lehet létrehozni/használni.

1. [Bináris letöltése és futtatása](#bin%C3%A1ris-f%C3%A1jlok-let%C3%B6lt%C3%A9se)
2. [A Lex elkészítése magadnak](#%EF%B8%8F-telep%C3%ADt%C3%A9si-%C3%BAtmutat%C3%B3)

## A Lex ingyenesen használható? Fizetnem kell a Lex használatáért?

**A Lex nyílt forráskódú, és mindig is az lesz, és szabadon használható kereskedelmi és/vagy magáncélú indokokért, mindaddig, amíg a szoftver eredetét félre nem ábrázolják.**

További információért tekintse meg a [LICENSE](https://github.com/cmd777/lex/blob/main/LICENSE) dokumentumot.

# ⚙️ Telepítési útmutató

## Építési útmutató

A projekt felépítésének egyetlen feltétele a [Go](https://go.dev/dl)

Ha a Go telepítve van, futtassa a következő parancsokat a terminálon
```cmd
git clone https://github.com/cmd777/lex.git
cd lex
cd src
go get -u
go build
```
és nagyjából kész is van, nincs más dolgunk, mint elindítani az épített bináris fájlt, és navigálni a `localhost:9090/r/{subreddit}` címre.

## Bináris fájlok letöltése

Ha a LEX-et automatikusan felépített bináris fájlokon keresztül szeretné telepíteni, lépjen a [releases](https://github.com/cmd777/lex/releases/latest) lapra, töltse le az operációs rendszer + processzor architektúra megfelelő zip-fájlt, és csomagolja ki.

- 32 bites Windows rendszerű gépek esetén töltse le a "lex-windows.zip"-et, majd indítsa el a lex-i386-windows.exe fájlt.
- 64 bites Windows rendszerű gépek esetén töltse le a "lex-windows.zip"-et, majd indítsa el a lex-amd64-windows.exe fájlt.

<br>

- 32 bites Linux gépek esetén töltse le a "lex-linux.zip"-et, majd indítsa el a lex-i386-linux fájlt
- 64 bites Linux gépek esetén töltse le a "lex-linux.zip"-et, majd indítsa el a lex-amd64-linux fájlt

<br>

- Intel MacOS gépek esetén töltse le a "lex-osx.zip"-et, majd indítsa el a lex-amd64-osx fájlt
- M1 MacOS gépek esetén töltse le a "lex-osx.zip"-et, majd indítsa el a lex-arm64-osx fájlt

Az indítás után navigáljon a `localhost:9090/r/{subreddit}` címre, és kész.

# 🚩 Hibák, problémák és egyéb fontos információk

A Lex még nagyon korai fejlesztési szakaszban van, sok funkció hiányzik, és hajlamos lehet a hibákra, mint például a szöveg túlcsordulása, a galéria gombjai furcsán működnek, és így tovább.

Igyekszem kijavítani a legtöbb kritikus hibát, mielőtt bármilyen változtatást végrehajtanék, de ha hibát talál, nyugodtan [hozzon létre egy problémát a githubon](https://github.com/cmd777/lex/issues).

Szabadidőmben hobbiból dolgozom a Lexen, így lehet, hogy lassú a fejlődés, köszönöm a türelmet!

# 🔨 A Lex létrehozásához használt technológiák

- [Go](https://go.dev/) ➡️ Körülbelül minden futtatásához
- [Humanize (go-humanize)](https://github.com/dustin/go-humanize) ➡️ Formázás emberbarát egységekre
- [Gzip (gin-contrib)](https://github.com/gin-contrib/gzip) ➡️ Gin Köztes szoftver a Gzip támogatás engedélyezéséhez
- [Gin](https://github.com/gin-gonic/gin) ➡️ HTTP Webes Keretrendszer
- [Bluemonday](https://github.com/microcosm-cc/bluemonday) ➡️ HTML-fertőtlenítő
- [Blackfriday](https://github.com/russross/blackfriday/tree/v2) ➡️ Markdown processzor
- [SASS](https://sass-lang.com) ➡️ CSS Kiterjesztés

<sub>Ezenkívül ezeket használtam a LEX létrehozásához:</sub>

- [Josefin Sans](https://fonts.google.com/specimen/Josefin+Sans) ➡️ Navigációs sáv, index betűtípus
- [Open Sans](https://fonts.google.com/specimen/Open+Sans) ➡️ Subreddit betűtípus
- [SVGRepo](https://www.svgrepo.com) ➡️ SVGs
- [Yqnn's SVG Path Editor](https://github.com/Yqnn/svg-path-editor) ➡️ Szinte az összes SVG szerkesztéséhez

⚠️ <sub>A fordítás nagy része a google fordító volt (A fordítási munka felgyorsítása érdekében), szóval lehet hogy nem a legjobb fordítás lett.</sub>