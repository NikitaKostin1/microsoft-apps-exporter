# Microsoft Apps Exporter

## The project is still being developed - the tests and deployments are to be introduced yet.

## Overview

The Microsoft Apps Exporter project was initiated to address the need for a dynamic, scalable solution to export SharePoint list data to a production database without requiring manual intervention for each new list. This solution ensures efficient data handling, automation, and adaptability to evolving business needs.

## Problem Statement

Managing data from multiple SharePoint lists in a production environment is often challenging, requiring frequent manual updates and code modifications whenever a new list is introduced. The key challenges include:

- Enabling automatic integration of new SharePoint lists without modifying the code.
- Ensuring a structured and reliable approach to map SharePoint list columns to database tables.
- Ensures reliability in a production setting with structured logging and error handling.
- Incorporating structured logging, error handling, and monitoring.
- Minimizing manual interventions while supporting a growing number of lists or other Microsoft products.

## Solution Approach

### Config-driven SharePoint List Integration

- A configuration file (`config.yaml`) where new SharePoint lists can be defined.
- Each list entry contains its **Site ID, List ID, Database Table Name, and Column Mapping**.
- Supports incremental updates to minimize redundant data transfers.

### Automated Data Export

- By webhook request fetches data from SharePoint via the **Graph API**.
- The data is transformed and stored in the corresponding production database table.

### Scalable Architecture

- Developed using **Go** for efficiency, concurrency, and reliability.
- Uses Viper for configuration management, ensuring flexibility.
- **Dockerized deployment** for consistent and scalable execution.
- **CI/CD pipeline** integration for automated testing and deployment.

### Robust Error Handling & Logging

- Implements structured logging for enhanced observability and debugging.
- Graceful handling of API failures with retry mechanisms.
- Tracks export status for auditing and monitoring purposes.

## Current Development Stage

This project is in its early development phase, with the following milestones completed:

- ✅ Core Logic Implementation: Fetching, transforming, and storing SharePoint list data.
- ✅ Unit Tests: Validation of key functionalities and edge cases.

### Next Steps

- Implementing **integration tests**.
- Finalizing **database migration strategies**.
- Developing a **CI/CD pipeline**.
- Deploying the system for real-world usage.

## Installation & Setup (WIP)

## Future Enhancements

- Implementing **real-time monitoring and alerts**.
- Optimizing **batch processing** for high-volume data.
- Enhancing **authentication and access control mechanisms**.

This project is actively evolving, and feedback is highly valued to ensure it meets the business requirements effectively!
