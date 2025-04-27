# Simple Go MapReduce Prototype

This project is a basic implementation of the MapReduce concept in Go, focusing on demonstrating concurrency using goroutines and leveraging Go generics for flexible data handling.

It's designed to process data in parallel using user-defined map and reduce functions. This prototype builds upon the principles of a worker pool for managing concurrent tasks.

**In this example, the prototype is used to analyze OpenSSH logs to identify suspicious IP addresses based on the number of failed login attempts.**

**Key Features:**

* **Concurrent Processing:** Utilizes Go goroutines to simulate parallel execution.
* **Generic Types:** Supports custom input, intermediate, and output data types.
* **User-Defined Functions:** Allows users to provide their own map and reduce logic.

**Learn More:**

For a more detailed explanation and walkthrough, check out my article on Medium: [Link](https://medium.com/@darshmamtora2122/map-reduce-prototype-log-analysis-6b529dfc488f)
