# Documentation Index

Welcome to the Srikandi Sehat GraphQL API documentation! This index will help you navigate to the right documentation for your needs.

## 📖 Documentation Overview

This project includes comprehensive documentation covering architecture, API usage, development, features, deployment, and quick references.

---

## 🎯 Start Here

### New to the Project?
Start with these documents in order:
1. **[README](../README.md)** - Project overview and quick start
2. **[Quick Reference](./QUICKREF.md)** - Common commands and examples
3. **[API Reference](./API.md)** - GraphQL API documentation

### Want to Understand the Architecture?
- **[Architecture Guide](./ARCHITECTURE.md)** - Complete hexagonal architecture explanation

### Ready to Develop?
- **[Developer Guide](./DEVELOPER.md)** - Development workflow and best practices
- **[Quick Reference](./QUICKREF.md)** - Common commands and snippets

### Need to Deploy?
- **[Deployment Guide](./DEPLOYMENT.md)** - Production deployment instructions

### Looking for Specific Features?
- **[Feature Documentation](./FEATURES.md)** - Detailed feature documentation

---

## 📚 Complete Documentation List

### 1. README.md
**Location**: `../README.md`  
**Purpose**: Project overview, quick start, and introduction

**Contents**:
- ✅ Project description
- ✅ Architecture overview
- ✅ Implemented features list
- ✅ Installation instructions
- ✅ Quick start guide
- ✅ Available make commands
- ✅ Project structure
- ✅ Security considerations

**Best for**: First-time users, project overview

---

### 2. ARCHITECTURE.md
**Location**: `./ARCHITECTURE.md`  
**Purpose**: Deep dive into hexagonal architecture implementation

**Contents**:
- ✅ Hexagonal architecture principles
- ✅ Core domain layer
- ✅ Ports layer (interfaces)
- ✅ Service layer (business logic)
- ✅ Adapters (infrastructure & API)
- ✅ Dependency injection
- ✅ Unit of Work pattern
- ✅ DataLoader implementation
- ✅ Benefits and best practices

**Best for**: Understanding the codebase structure, architectural decisions

---

### 3. API.md
**Location**: `./API.md`  
**Purpose**: Complete GraphQL API reference

**Contents**:
- ✅ GraphQL endpoint information
- ✅ Authentication mechanism
- ✅ All queries with examples
- ✅ All mutations with examples
- ✅ Request/response formats
- ✅ Error handling
- ✅ Data types
- ✅ cURL examples
- ✅ Best practices

**Best for**: API consumers, frontend developers, testing

---

### 4. DEVELOPER.md
**Location**: `./DEVELOPER.md`  
**Purpose**: Guide for developers contributing to the project

**Contents**:
- ✅ Getting started
- ✅ Development workflow
- ✅ Commit message conventions
- ✅ Code organization
- ✅ Step-by-step feature addition guide
- ✅ Testing guidelines
- ✅ Database migrations
- ✅ Code style guidelines
- ✅ Common tasks
- ✅ Troubleshooting
- ✅ Best practices

**Best for**: New contributors, adding features, maintaining code

---

### 5. FEATURES.md
**Location**: `./FEATURES.md`  
**Purpose**: Detailed documentation of all implemented features

**Contents**:
- ✅ Authentication system
  - Registration flow
  - Login flow
  - Token refresh
  - Middleware implementation
- ✅ User management
  - Profile retrieval
  - Profile updates
- ✅ Password management
  - Change password
  - Forgot password (OTP)
- ✅ Email services
  - Email templates
  - SMTP/Mailgun configuration
- ✅ Security features
  - Argon2id password hashing
  - PASETO tokens
- ✅ Database management
  - Unit of Work pattern
  - DataLoader (N+1 prevention)
- ✅ Logging and monitoring

**Best for**: Understanding how features work, implementation details

---

### 6. DEPLOYMENT.md
**Location**: `./DEPLOYMENT.md`  
**Purpose**: Production deployment and operations guide

**Contents**:
- ✅ Prerequisites
- ✅ Environment setup
- ✅ Database setup and configuration
- ✅ Application deployment (systemd)
- ✅ Docker deployment
- ✅ Reverse proxy setup (Nginx)
- ✅ SSL/TLS configuration
- ✅ Production checklist
- ✅ Monitoring and logging
- ✅ Scaling strategies
- ✅ Backup and recovery
- ✅ Troubleshooting

**Best for**: DevOps, system administrators, production deployment

---

### 7. QUICKREF.md
**Location**: `./QUICKREF.md`  
**Purpose**: Quick reference for common tasks and commands

**Contents**:
- ✅ Quick start commands
- ✅ Common make commands
- ✅ GraphQL query examples
- ✅ HTTP headers guide
- ✅ cURL examples
- ✅ Environment variables
- ✅ Docker commands
- ✅ Database commands
- ✅ Common errors and solutions
- ✅ Project structure overview
- ✅ Development workflow
- ✅ Security checklist

**Best for**: Quick lookups, copy-paste snippets, daily development

---

## 🎓 Learning Paths

### Path 1: API Consumer
If you're building a client application:
1. [README](../README.md) - Understand the project
2. [API Reference](./API.md) - Learn the API
3. [Quick Reference](./QUICKREF.md) - Get examples

### Path 2: Backend Developer
If you're developing features:
1. [README](../README.md) - Setup and introduction
2. [Architecture Guide](./ARCHITECTURE.md) - Understand the structure
3. [Developer Guide](./DEVELOPER.md) - Learn the workflow
4. [Quick Reference](./QUICKREF.md) - Daily reference

### Path 3: DevOps/Deployment
If you're deploying to production:
1. [README](../README.md) - Understand the project
2. [Deployment Guide](./DEPLOYMENT.md) - Deploy step-by-step
3. [Quick Reference](./QUICKREF.md) - Common commands

### Path 4: Architect/Tech Lead
If you're evaluating or reviewing:
1. [README](../README.md) - Project overview
2. [Architecture Guide](./ARCHITECTURE.md) - Architecture deep dive
3. [Features](./FEATURES.md) - Implementation details
4. [Deployment Guide](./DEPLOYMENT.md) - Production setup

---

## 🔍 Quick Find

### How do I...

**...start the project?**
- [README - Installation & Setup](../README.md#installation--setup)
- [Quick Reference - Quick Start](./QUICKREF.md#quick-start)

**...use the API?**
- [API Reference](./API.md)
- [Quick Reference - GraphQL Queries](./QUICKREF.md#graphql-queries--mutations)

**...add a new feature?**
- [Developer Guide - Adding New Features](./DEVELOPER.md#adding-new-features)

**...understand the architecture?**
- [Architecture Guide](./ARCHITECTURE.md)

**...deploy to production?**
- [Deployment Guide](./DEPLOYMENT.md)

**...run tests?**
- [Developer Guide - Testing](./DEVELOPER.md#testing)
- [Quick Reference - Testing](./QUICKREF.md#common-commands)

**...create a database migration?**
- [Developer Guide - Database Migrations](./DEVELOPER.md#database-migrations)
- [Quick Reference - Database](./QUICKREF.md#database)

**...configure environment variables?**
- [README - Configure Environment](../README.md#configure-environment-variables)
- [Quick Reference - Environment Variables](./QUICKREF.md#environment-variables)

**...troubleshoot issues?**
- [Developer Guide - Troubleshooting](./DEVELOPER.md#troubleshooting)
- [Deployment Guide - Troubleshooting](./DEPLOYMENT.md#troubleshooting)
- [Quick Reference - Common Errors](./QUICKREF.md#common-errors)

---

## 📋 Document Status

| Document | Status | Last Updated | Completeness |
|----------|--------|--------------|--------------|
| README.md | ✅ Complete | 2025-11-10 | 100% |
| ARCHITECTURE.md | ✅ Complete | 2025-11-10 | 100% |
| API.md | ✅ Complete | 2025-11-10 | 100% |
| DEVELOPER.md | ✅ Complete | 2025-11-10 | 100% |
| FEATURES.md | ✅ Complete | 2025-11-10 | 100% |
| DEPLOYMENT.md | ✅ Complete | 2025-11-10 | 100% |
| QUICKREF.md | ✅ Complete | 2025-11-10 | 100% |

---

## 🤝 Contributing to Documentation

Found an error or want to improve the documentation?

1. **Fork the repository**
2. **Make your changes**
3. **Test any code examples**
4. **Submit a pull request**

### Documentation Guidelines

- Use clear, concise language
- Include code examples where appropriate
- Keep formatting consistent
- Update the index if adding new docs
- Test all code snippets
- Include screenshots for complex UI steps

---

## 📞 Getting Help

If you can't find what you're looking for:

1. **Search the documentation** - Use Ctrl+F in your browser
2. **Check the Quick Reference** - Common tasks and commands
3. **Review the FAQ** - In Developer Guide
4. **Check existing issues** - On GitHub
5. **Ask the team** - Create a new issue

---

## 📈 Documentation Roadmap

### Planned Additions

- [ ] FAQ section
- [ ] Video tutorials
- [ ] Architecture decision records (ADRs)
- [ ] Performance tuning guide
- [ ] Testing strategy guide
- [ ] Security audit checklist
- [ ] Migration guides for major versions
- [ ] API versioning strategy

---

## 📝 Feedback

We value your feedback on this documentation!

- **What's working well?**
- **What's confusing?**
- **What's missing?**

Please open an issue on GitHub with the tag `documentation` to share your thoughts.

---

## 📜 License

This documentation is part of the Srikandi Sehat GraphQL API project and is licensed under the MIT License.

---

**Last Updated**: 2025-11-10  
**Version**: 1.0.0  
**Maintained by**: Srikandi Sehat Development Team

---

## Quick Links

- [Main README](../README.md)
- [Architecture Guide](./ARCHITECTURE.md)
- [API Reference](./API.md)
- [Developer Guide](./DEVELOPER.md)
- [Feature Documentation](./FEATURES.md)
- [Deployment Guide](./DEPLOYMENT.md)
- [Quick Reference](./QUICKREF.md)

---

*Happy coding! 🚀*
