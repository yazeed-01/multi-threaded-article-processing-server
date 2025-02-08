# Multi-Threaded Article Processing Server

## Overview

This project implements a multi-threaded server that processes user articles from Medium. It stores the processed articles in Elasticsearch for full-text search and sub-text search and tracks the processing status in an SQL database.

The server receives a user ID, fetches the top 20 articles for that user, and processes each article concurrently using multiple threads. The system ensures efficient queue management and provides API endpoints for tracking the processing status.

## Features

- **Fetch Articles from Medium**: The server fetches the top 20 articles for a given user from Medium.
- **Multi-threaded Processing**: Articles are processed concurrently, with each article stored in Elasticsearch for full-text search.
- **Tracking**: The system tracks the status of each article in an SQL database.
- **Queue Management**: Uses an array blocking queue for efficient thread management and a database-based queue for persistent queuing.
- **API Endpoints**:
    - Submit User ID for processing
    - Fetch processed articles for a user
    - Check processing status
 
