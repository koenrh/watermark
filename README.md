# watermark

A simple command-line tool to quickly add a text-based watermarks to all pages
of a PDF document.

## Usage

Add a text-based watermark by specifying one or more `-text` options, one per text
line. Optionally, specify `-date` to include the current date.

`watermark -text "Framer B.V." -text "Confidential" -date very_secret_report.pdf`
