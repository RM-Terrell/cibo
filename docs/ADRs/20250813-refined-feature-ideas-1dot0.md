# Refined feature ideas for 1.0

- Save daily price data to parquet
- Fair value price calculation pipeline, save to parquet with daily price data
- When ignoring a bad price value, make sure its logged somehow and bubbled up into the TUI or WUI or in a log file
- Make the full application run as a Docker container (not just dev container), and setup compilation pipelines for MacOS and Linux to run as a bin.

# 2.0

- WebUI
- Robust messaging in the case of a companies earnings being all negative when doing graham lynch calculations
