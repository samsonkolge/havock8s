# Contributing to StatefulChaos

Thank you for your interest in contributing to StatefulChaos! This document provides guidelines and instructions for contributing to this project.

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