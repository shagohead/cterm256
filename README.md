cterm256
========

Generate full ANSI color table based on primary 8 or 16 colors.

It is like base16, but with two benefits:

- it has background colors too (based on foreground ones mixed with primary background)
- color scheme applied only on terminal emulator, not for each your application separately.

The latter requries configuring the applications to use ANSI indexed colors.
Accordingly, the application must have this capability.
Luckily, most of them have it.

With that approach you will have one color system for most of your terminal apps.
And when you will change color scheme, changes will be applied for all of these apps at once.

## Installation

Color table generator is an Go CLI app. Which can be installed like this:

```go
go install -v github.com/shagohead/cterm256/cmd/cterm256@latest
```

Currently generator supports only [kitty](https://sw.kovidgoyal.net/kitty/) color themes, but it can be extended by implementation of [`cterm.FileType`](https://pkg.go.dev/github.com/shagohead/cterm256/internal/cterm#ColorScheme) interface and adding that implementation into `internal/filetypes/filetypes.go` flag values.

Configurations which are uses generated color scheme located are in `./configs` directory.

## Reason of this project

Years ago i was used base16 terminal themes. But for me, 16 is not enough, i need variations of these 16 base colors. Mainly due to the need for different backgrounds: added/removed/changed blocks of code in diffs and UI elements of TUI apps like Vim/Neovim, Tig and so on.

I decided to use 8-bit ANSI color table. Most of color schemes at internet defines only standard & high intensity colors which stands at 0-7 and 8-15 numbers in that table. Another 240 colors are: 6x6x6 cube and 24-length grayscale. How to generate these colorspace from theme defined 8/16 colors?

## Generated color table cheatsheat

There two types of colors that can be distinguished: foreground and background.

Background colors are based created from main red, green and blue color hue and saturation, but with lightness of main background color. Red, green and blue are used for creation of 5 colors each with lightness from background to foreground. These gradients are:

- Red: 52, 88, 124, 160, 196.
- Green: 22, 28, 34, 40, 46.
- Blue: 17, 18, 19, 20, 21.

Grayscale 232-252 corresponds to transition from background color to foreground or white or «bright white» color, which one will be more contrast to background.

Lightness variations generated in HSLuv (developer oriented CIELUV) colors space, which produces [accurate results](https://www.hsluv.org/comparison/), which [especially important for backgrounds](https://www.kuon.ch/post/2020-03-08-hsluv/).

In light color schemes, if white is lighter than black they will be swaped. Because switching between light and dark themes should not change semantics of colors, and white color should be high contrast to background.

## TODO
Fix colors out of range.

Add images & video previews as examples of how it works and feels.
