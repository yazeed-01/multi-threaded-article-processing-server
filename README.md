# multi-threaded-article-processing-server
 A multi-threaded server that processes user articles from Medium, stores them in Elasticsearch for full-text search, and tracks processing states in an SQL database. The system uses efficient thread management through an array blocking queue, with a database-based queue option for persistent processing.
