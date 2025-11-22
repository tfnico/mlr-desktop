Goal: Provide a desktop tool for transforming files using Miller (mlr)

Background: The open source Go library/command mlr provides a powerful way to operate on large files (e.g. CSV, TSV, JSON, etc.).

Its command line interface has a lot of options and can be difficult to learn. Therefore we'll create a desktop application to make it easier to use.

The application should offer three basic components:

* input: a way to add text or a text file to transform, as well format options for the input file. It should also provide a free form text field with options.
* transformation: a way to select the transformation to apply. This should allow adding a chain of mlr verbs to be combined with the "then" keyword. It should be possible to provide options for each verb in a text field.
* output: a preview of the transformed text, as well as format options for the output. It should also provide a free form text field with options.

Implementation details:

The application should be written in Go and use the Wails framework for the desktop application.

Secondary goals:

* Save/load command: It should be possible to save/load the transformation instructions.
* Export as Go program: It should be possible to export the prepared transformation as a main.go file which uses the mlr stream package to perform the transformation.

