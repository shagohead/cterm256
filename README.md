cterm256
========

8/16 => 256 color table generator & collection of 256-table'based color schemes & configurations for CLI/TUI apps.

### Goals of this project
- Use same color coding across ~~all~~ most of CLI/TUI apps
- Ability to change colors of ~~all~~ most apps output with one command
- Generate 256 colors table based on 8/16 colors of any custom theme

#### Use same color coding across ~~all~~ most of CLI/TUI apps
I think it’s nice to see semantically colored text in one color coding space across diffirent terminal apps output.

#### Ability to change colors of ~~all~~ most apps output with one command
If applications will use indexed color table, you can just change colors in that table. All that apps will change thier colorschemes with that one action. It is especially usefull when you're need to switch from dark to light theme and vice-versa.

#### Generate 256 colors table based on 8/16 colors of any custom theme
Most (if not all) of color schemes for terminals defines only 8 or 16 colors. With that assumption we can try to generate another colors in table based on that 8/16.

## Reason of this project

Years ago i was used base16 terminal themes. But for me, 16 is not enough, i need variations of these 16 base colors. Mainly due to the need for different backgrounds: added/removed/changed blocks of code in diffs and UI elements of TUI apps like Vim/Neovim, Tig and so on.

I decided to use 8-bit ANSI color table. Most of color schemes at internet defines only standard & high intensity colors which stands at 0-7 and 8-15 numbers in that table. Another 240 colors are: 6x6x6 cube and 24-length grayscale. How to generate these colorspace from theme defined 8/16 colors?

## Algorithm

Grayscale can be maded by generating gradient of lightness based on black color (which number is 0). For 6x6x6 cube i've found these correlations:

- Top-left color of first side is the black (#0) color.
- Top row (5 colors, except top-left black) of first side is the gradient of blue (#12).
- Left column (5 colors, except top-left black) of first side is the gradient red (#9).
- Top-left color of all 5 remaining sides is the gradient of green (#10).
- Top row of remaining sides is the gradient made by mixing of top-left green and blue from the top row of the first side.
- Left column of remaining sides is the gradient made by mixing of top-left green and red from the left column of the first side.
- All remain colors made by mixing intersecting top row and left column and top-left cell colors.

Because some of these colors will be used as backgrounds, they need to have same lightness as the background. With that, cube generated from standard colors will be slightly diffirent to standard cube.

Lightness variations generated in HSLuv (developer oriented CIELUV) colors space, which produces [accurate results](https://www.hsluv.org/comparison/), which [especially important for backgrounds](https://www.kuon.ch/post/2020-03-08-hsluv/).

## TODO
Try github.com/lucasb-eyer/go-colorful, which can prevent creating wrong colors for RGB: https://github.com/lucasb-eyer/go-colorful?tab=readme-ov-file#q-labluvhcl-seem-broken-your-library-sucks

Autofix black & white colors for light-background themes by reversing them.

Add images & video previews as examples of how it works and feels.
