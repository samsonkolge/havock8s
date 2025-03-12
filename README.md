# Havock8s Documentation

This directory contains the source files for the Havock8s documentation site, which is built using [Jekyll](https://jekyllrb.com/) and hosted on [GitHub Pages](https://pages.github.com/).

## Local Development

To run the documentation site locally:

1. Install Ruby and Bundler:
   ```bash
   # For Ubuntu/Debian
   sudo apt-get install ruby-full build-essential zlib1g-dev
   
   # For macOS (using Homebrew)
   brew install ruby
   
   # Install Bundler
   gem install bundler
   ```

2. Install dependencies:
   ```bash
   cd docs
   bundle install
   ```

3. Run the Jekyll server:
   ```bash
   bundle exec jekyll serve
   ```

4. Open your browser and navigate to `http://localhost:4000/havock8s/`

## Documentation Structure

- `_config.yml`: Jekyll configuration file
- `_layouts/`: Custom layout templates
- `assets/`: Static assets (CSS, images, etc.)
- `*.md`: Markdown content files

## Adding Content

### Adding a New Page

1. Create a new Markdown file in the `docs` directory
2. Add the following front matter at the top of the file:
   ```yaml
   ---
   layout: default
   title: Your Page Title
   ---
   ```
3. Add your content using Markdown

### Adding Images

1. Place image files in the `assets/images/` directory
2. Reference them in your Markdown files:
   ```markdown
   ![Alt text](assets/images/your-image.png)
   ```

## Styling

The documentation site uses a customized version of the Jekyll Cayman theme. Custom styles are defined in `assets/css/style.css`.

## Deployment

The documentation is automatically deployed to GitHub Pages when changes are pushed to the main branch.

## Contributing

Please see the main [CONTRIBUTING.md](../CONTRIBUTING.md) file for general contribution guidelines.

For documentation-specific contributions:

1. Keep the documentation up-to-date with the code
2. Use clear, concise language
3. Include examples where appropriate
4. Test any code examples to ensure they work
5. Preview your changes locally before submitting a pull request 