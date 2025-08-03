# Changelog

## [v2.2.0] - 03 August 2025

    * Fix refresh crashing some applications because wParam was not NULL
    * Prefix with dash '-' in INI section to delete variable. Leaving it empty issues a warning.
    * Some code refactoring

## [v2.1.1] - 04 June 2019

    * Fixed panic in call to Windows function
    * Work around appveyor build hangs

## [v2.0.0] - 19 July 2017

    * Simplify options
    * Allow stdin
    * Refresh environment via WM_SETTINGCHANGED

## [v1.1.1] - 12 June 2015

    * Fix handling of sections that appear twice

## [v1.1.0] - 11 June 2015

    * Internal rework
    * Added more tests
    * Added -checkpaths option

## [v1.0.0] - 08 June 2015

    * First version
