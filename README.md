# STool

[![Run Tests](https://github.com/unquabain/stool/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/unquabain/stool/actions/workflows/test.yml)

A tool for querying and reformatting structured files.

```
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
    dot-delimited path elements.
      Example: locale.en.errors.0.text

    Brackets can be used if the element might have conflicting syntax.
      Example: keys["key with spaces and dot."].value

    A few functions have been defined. They are "length()", "yaml()",
    "json()" and "jspretty()".
      Example: local.en.errors.yaml()

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

## Example

Given a file like:

```json
{
  "auths": {
    "my-private.site": {
      "auth": "dXNlcjpteXNlY3JldHBhc3N3b3JkCg=="
    }
  }
}
```

This command will extract the password:

```bash
stool -i config.json \
  -s 'auths["my-private.site"].auth' \
  -t '{{ . | b64dec | splitList ":" | last }}'
```

The filters `b64dec`, `splitList` and `last` all come from Masterminds Sprig.
