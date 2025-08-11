# Bubbletea

Bubbletea was chosen in order to provide a slightly more interactive and UI like terminal user interface (TUI). Bubbletea specifically was chosen for its wide range of support for complex visualizations in the terminal and is probably overkill for this application, but room to grow was important.

Official docs on Bubbletea can be found [here](https://github.com/charmbracelet/bubbletea)

The `tui_model.go` code follows a Model-View-Update structure as described in the libraries docs, and the function by those names are intended to carry out those tasks. This project was the original authors first venture in Bubbletea, and their first venture into a TUI, but care was taken to stick close to best practices for such systems.
