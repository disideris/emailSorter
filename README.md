# emailSorter

There are two implementations of the customer importer.
customerimporter and concurrentcustomerimporter.

customerimporter is a more straight forward implementation leveraging byte indexes for optimizing performance. 

concurrentcustomerimporter is a concurrent implementation. This importer creates a worker pool, reads each line of file, extracts the domain (leveraging the byte indexes as above), pushes domain strings to a work channel and does some concurrent optimazations by grouping domains before sending them to actual workers. The worker pool also takes a function as parameter (for demonstrating reasons) in case we wanted to perform any extra processing.
