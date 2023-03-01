# LEX
[![build](https://github.com/cmd777/lex/actions/workflows/build_all_os.yml/badge.svg)](https://github.com/cmd777/lex/actions/workflows/build_all_os.yml)

# Tartalomjegyzék

**🇭🇺 Magyar**

[**ℹ️ Információk a projektről**](#-információk-a-projektről)
- [Mi az a Lex?](#mi-az-a-lex)
  - [Miért hoztad létre a Lexet?](#miért-hoztad-létre-a-lexet)
  
[**⚙️ Telepítési útmutató**](#-telepítési-útmutató)
- [Építési útmutató](#építési-útmutató)
- [Binárisok letöltése](#bináris-fájlok-letöltése)

[**🔬 Böngésző kompatibilitás**](#-böngésző-kompatibilitás)

[**🚩 Hibák, problémák és egyéb fontos információk**](#-hibák-problémák-és-egyéb-fontos-információk)

[**📜 Jogi nyilatkozat**](#-jogi-nyilatkozat)

[**🧰 A Lex létrehozásához használt technológiák**](#-a-lex-létrehozásához-használt-technológiák)

[**📝 Egyéb Információk**](#-egyéb-információk)
- [Szükségem van-e reddit fiókra a Lex használatához?](#szükségem-van-e-reddit-fiókra-a-lex-használatához)
- [Upvoteolhatom / Kommentelhetek posztokra?](#upvoteolhatom--kommentelhetek-posztokra)

---

# ℹ️ Információk a projektről

## Mi az a Lex?
A LEX (LazerEX) egy könnyű, nyílt forráskódú kezelőfelület a reddithez, amit Go-ban készítettem.

## Miért hoztad létre a Lexet?

A Reddit nagyon lassú tud lenni, és rengeteg internetet használ, ez nem titok. Nagyon bosszantó volt, hogy vagy nagyon sokáig vártam, amíg egy poszt betöltődik, vagy a felhasználói felület nyilvánvaló ok nélkül lefagyott. Ez csak néhány olyan probléma, amit én néhány böngészőben (Chrome, Firefox, Brave stb.) tapasztaltam. A Lex akár 60%-kal kevesebb internetet használhat fel, miközben nagyon hasonló kép-/videóminőséget biztosít, és az interaktívvá válás ideje átlagosan körülbelül 800 ms-1,2s (az összes kép, szkript, stíluslap betöltésével együtt, megjegyzés: a videók nincsenek előre betöltve.)

Ami a felhasználói felületet illeti, az újabb reddit újratervezést választottam, de egy kis finomítással.

## Hogyan kell használni a Lexet?

Sajnos, ellentétben más nyílt forráskódú reddit frontendekkel, a **Lex nem biztosít bemutató webhelyet**, és kétféleképpen lehet létrehozni/használni.

1. [Bináris letöltése és futtatása](#bináris-fájlok-letöltése)
2. [A Lex elkészítése magadnak](#építési-útmutató)

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

# 🔬 Böngésző kompatibilitás
|Chrome  | Edge | Safari | Firefox  | Opera  | IE   |
|:-----: | :--: | :----: | :------: | :---:  | :--: |
| > 41   | > 18 | > 10   |  > 48    | > 28   | ❌  |

<sub>A táblázat a böngészők ajánlott verzióit jeleníti meg, bár a régebbi verziók *lehet*, hogy valamelyest kompatibilisek.</sub>

<sub>A [Can I Use](https://caniuse.com) adatok alapján</sub>

# 🚩 Hibák, problémák és egyéb fontos információk

A Lex még nagyon korai fejlesztési szakaszban van, sok funkció hiányzik, és hajlamos lehet a hibákra, mint például a szöveg túlcsordulása, a galéria gombjai furcsán működnek, és így tovább.

Igyekszem kijavítani a legtöbb kritikus hibát, mielőtt bármilyen változtatást végrehajtanék, de ha hibát talál, nyugodtan [hozzon létre egy problémát a githubon](https://github.com/cmd777/lex/issues).

A tervezett új funkciókra vagy a javítandó dolgokra vonatkozóan tekintse meg a [TODO Listát](https://github.com/cmd777/lex/blob/main/TODO.md)

<sub>**A lista nem teljes. Az új funkciók vagy a javítandó hibák felfedezésekor, vagy új ötletek alapján kerülnek hozzáadásra**</sub>

<sub>**A lista nem prioritás szerint van rendezve.**</sub>

Szabadidőmben hobbiból dolgozom a Lexen, így lehet, hogy lassú a fejlődés, köszönöm a türelmet!

# 📜 Jogi nyilatkozat

A LEX nem összefüggő, nem támogatott, és nem jóváhagyott a Reddit által.

A LEX-en megjelenő összes tartalom a Redditről származik. LEX nem őriz semmilyen tartalmat.

Bármely bejegyzéshez kapcsolódó probléma esetén, például szerzői jogi sértés, védjegyjogi sértés vagy a Reddit közösségi szabályainak megsértése esetén, a jelentéseket a Redditnek kell címezni.

# 🧰 A Lex létrehozásához használt technológiák

- [Go](https://go.dev) ➡️ Programozási Nyelv
- [Humanize (go-humanize)](https://github.com/dustin/go-humanize) ➡️ Idő, számok stb. formázása emberbarát egységekre
- [Fiber](https://github.com/gofiber/fiber) ➡️ HTTP Webes Keretrendszer
- [Go-JSON](https://github.com/goccy/go-json) ➡️ Gyors JSON dekóder
- [Bluemonday](https://github.com/microcosm-cc/bluemonday) ➡️ HTML-fertőtlenítő
- [Blackfriday](https://github.com/russross/blackfriday/tree/v2) ➡️ Markdown processzor
- [SASS](https://sass-lang.com) ➡️ CSS Kiterjesztés

<sub>Ezenkívül ezeket használtam a LEX létrehozásához:</sub>

- [Josefin Sans](https://fonts.google.com/specimen/Josefin+Sans) ➡️ Navigációs sáv, index betűtípus
- [Open Sans](https://fonts.google.com/specimen/Open+Sans) ➡️ Subreddit betűtípus
- [SVGRepo](https://www.svgrepo.com) ➡️ SVGs
- [Yqnn's SVG Path Editor](https://github.com/Yqnn/svg-path-editor) ➡️ Szinte az összes SVG szerkesztéséhez

# 📝 Egyéb információk

## Szükségem van-e reddit fiókra a Lex használatához?

Nem, nem szükséges reddit fiók a Lex használatához.

## Upvoteolhatom / Kommentelhetek posztokra?

Sajnos, jelenleg nem tervezem az upvoteolás és a kommentálás fejlesztését.

⚠️ <sub>A fordítás nagy része a google fordító volt (A fordítási munka felgyorsítása érdekében), szóval lehet hogy nem a legjobb fordítás lett.</sub>