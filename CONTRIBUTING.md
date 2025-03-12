# Contributing to havock8s

We love your input! We want to make contributing to havock8s as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

## Development Process

We use GitHub to host code, to track issues and feature requests, as well as accept pull requests.

### Pull Requests

1. Fork the repository to your own GitHub account
2. Clone the project to your machine
3. Create a branch locally with a succinct but descriptive name
4. Commit changes to the branch
5. Follow any formatting and testing guidelines specific to this repo
6. Push changes to your fork
7. Open a PR in our repository and follow the PR template so that we can efficiently review the changes

### Development Workflow

Here's how you can set up your development environment:

```bash
# Clone your fork of the repo
git clone https://github.com/YOUR_USERNAME/havock8s.git

# Navigate to the newly cloned directory
cd havock8s

# Add the original repository as a remote called "upstream"
git remote add upstream https://github.com/havock8s/havock8s.git

# Get the latest changes from upstream
git pull upstream main

# Create a new topic branch
git checkout -b my-feature-branch

# Make your changes, add tests, and make sure the tests still pass
make test

# Commit your changes
git commit -m "Descriptive commit message"

# Push your changes to your fork
git push origin my-feature-branch
```

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct. Please be respectful and considerate of others.

## Getting Started

### Prerequisites

- Go 1.19 or higher
- Kubernetes cluster (for testing)
- kubectl configured to communicate with your cluster
- Docker (for building container images)

### Setting Up Your Development Environment

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/statefulchaos.git
   cd statefulchaos
   ```
3. Add the original repository as an upstream remote:
   ```bash
   git remote add upstream https://github.com/statefulchaos/statefulchaos.git
   ```
4. Install dependencies:
   ```bash
   go mod download
   ```

## Development Workflow

1. Create a new branch for your feature or bugfix:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes, following our coding conventions

3. Run tests to ensure your changes don't break existing functionality:
   ```bash
   make test
   ```

4. Build and run the operator locally:
   ```bash
   make run
   ```

5. Commit your changes with a descriptive commit message:
   ```bash
   git commit -m "Add feature: your feature description"
   ```

6. Push your branch to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

7. Create a Pull Request from your fork to the main repository

## Adding New Chaos Types

StatefulChaos is designed to be extensible. To add a new chaos type:

1. Define the chaos type in the CRD (`api/v1alpha1/statefulchaosexperiment_types.go`)
2. Create a new injector in `pkg/chaos/your_chaos_type.go`
3. Register your injector in `pkg/chaos/injector.go`
4. Add controller logic in `controllers/statefulchaosexperiment_controller.go`
5. Create example YAML files in the `examples/` directory
6. Add tests for your new chaos type

See the [Developer Guide](docs/developer-guide.md) for more detailed instructions.

## Pull Request Process

1. Ensure your code passes all tests and linting
2. Update documentation, including the README if necessary
3. Add or update tests for your changes
4. Make sure your commits are clean and focused
5. Your PR will be reviewed by maintainers, who may request changes
6. Once approved, your PR will be merged

## Coding Conventions

- Follow standard Go coding conventions
- Use meaningful variable and function names
- Write clear comments, especially for complex logic
- Add appropriate error handling
- Include unit tests for new functionality

## Testing

- Write unit tests for all new functionality
- Ensure existing tests pass with your changes
- For chaos injectors, include tests that verify both injection and cleanup work correctly
- Test with different Kubernetes versions if possible

## Documentation

- Update documentation for any changed functionality
- Document new features thoroughly
- Include examples of how to use new features
- Keep the README up to date

## Community

- Join our [Slack channel](#) for discussions
- Participate in issues and pull requests
- Help answer questions from other users
- Share your use cases and feedback

Thank you for contributing to StatefulChaos! 