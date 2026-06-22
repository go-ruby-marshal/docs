# go-ruby-marshal/docs

Documentation site for [`go-ruby-marshal/marshal`](https://github.com/go-ruby-marshal/marshal) —
a pure-Go (`CGO_ENABLED=0`) implementation of Ruby's **Marshal** binary
serialization format (version 4.8), byte-for-byte identical to MRI Ruby.

Built with [MkDocs Material](https://squidfunk.github.io/mkdocs-material/) and
versioned with [mike](https://github.com/jimporter/mike). Served at
**<https://go-ruby-marshal.github.io/docs/>**.

## Local development

```sh
python -m venv .venv && . .venv/bin/activate
pip install -r requirements.txt
mkdocs serve            # live preview at http://127.0.0.1:8000/
mkdocs build --strict   # build to ./site, failing on warnings
```

The deploy is handled by `.github/workflows/docs.yml`: on every push to `main`,
mike publishes the versioned site to the `gh-pages` branch, which GitHub Pages
serves.

## License

BSD-3-Clause.
