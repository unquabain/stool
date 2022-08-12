# STool

[![Run Tests](https://github.com/unquabain/stool/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/unquabain/stool/actions/workflows/test.yml)

A tool for querying and reformatting structured files.

```
flag provided but not defined: -hel
A tool for querying and reformatting structured files.

Usage ./stool [options] filename [outputfile]

OPTIONS:
  --in -i
    The structured file to read, or - for STDIN (the default).
    This can also be specified with the first positional argument.

  --out -o
    The file to write out, or - for STDOUT (the default).
    This can also be specified with the second positional argument.
  
  --search -s 
    The search path. This is a little like JQuery. It consists mainly of
    dot-delimited path elements. There is an array of results, and each
    segment of the query operates on all the results in parallel. Results
    for which the command does not make sense (e.g. indexing an integer)
    are discarded. Most segment types reduce the number of results or
    keep it the same, brace filters, [*] and flatten() can increase them.
		The template is rendered once for each result. If you want the template
		to operate on all the results as an array, the results() function is
		provided.
      Example: locale.en.errors.not_found.text
			Example: locale.en.errors.keys().results().jsonpretty()
			Example: locale.en.errors[text == "File not found"].code
			Example: locale.en.messages[length() > 5][4]

    Brackets can be used if the element might have conflicting syntax.
      Example: keys["key with spaces and dot."].value

    A few functions have been defined. They use "()" to indicate that they're functions,
      and not members, but none of them take any arguments.

      len(), length(): If the value is an array, map or string, returns the
        length. "1" otherwise.

      json(), js(): Replaces each result with a rendered JSON string.

      jsonpretty(), jspretty(), jpretty(): Same as json, but renders line returns and indents.

			jsoneval(), jeval(): Evaluates each string in the results as embedded JSON

      yaml(), yml(): Replaces each result with a rendered YAML string.

			yamleval(), yeval(): Evaluates each string in the results as embedded YAML

      keys(): Replaces each result that's a map with an array of its keys.

      flatten(), flat(): Expands any result that's a collection into individual results. The 
        template will be rendered for each one individually.

      results(): Takes all the current results and makes them into a single result that's an array.
        The template will be rendered just once with the array of all results as its data.


    The special value "[*]" will resolve to all the values of an array
    or dictionary.
      Example: contacts[*].name

    You can also do simple tests that consist of a path, a comparison,
    and a value.
      Example: contacts[zip_code == "90210"].name
      Example: contacts[phones.length() > 2].name

  --template -t
    A Go Text Template to render the results. The template will be
    repeated for each result, so if you used [*] anywhere in your search
    pattern, the entire template will be repeated in the output. You may
    use Masterminds Sprig functions as well as the "yaml" and "json"
    functions.
      Example: The secret is {{ .client_secret | squote }}

    Currently, the only way to enter a carriage return is with {{ '\x0A' }}.
    So if you're expecting multiple documents, you may want to use
    the --template-file option.
      Example: ./stool -s '[*][client_title == "ekg"].client_id' \
                   -t "ID: {{ . }}{{ "'"\x0A"'" }}" \
                   rockauth_clients.yml

  --template-file -T
    For longer templates, you may specify a file to read instead of
    putting the template on the command line.

  --format -f
    The input file format. If the program cannot guess the file format,
    you may specify it as either "json" or "yaml".
```

The filters `b64dec`, `splitList` and `last` all come from Masterminds Sprig.
