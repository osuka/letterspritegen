# Auto generation of text images

## Letters

Small helper tool that, given a template (like ./letter-template-01.svg), generates
a set of square transparent images ready to be used as sprite or tiles for word game
applications.

The templates are simply inkscape templates using fonts (keep in mind licenses) and
inkscape effects. Each template has one target object that will be rendered using
inkscape (calling it via command line).

Tested on linux/nixos only.

```sh
./generate-letter -template letter-template-03.svg 
```

