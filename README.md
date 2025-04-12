# Pet Daisy

Pet Daisy started as a simple WebSocket server to act as a backend for my JavaFX final. After the final I decided to port it to the web for fun, and it picked up some traction at my school. Since then it has become my stomping ground for developing performant and scalable backend systems.

This README acts as a quick journal for me to look back on and see just how far I've come. I'm writing the first version pretty late at night, so please forgive any grammatical errors. thx.


# What I've learned so far

This project has taught me a lot of important software engineering concepts that I was simply never exposed to while learning programming.

### 1. GoLang

This project is my first time using Go and it has been excellent for diving into the features that make Go special such as how it handles concurrency.

### 2. Project Structuring

This is by far the largest project I have worked on and so it forced me to understand how to structure a project. 

In the beginning, everything was at the root of the project and there was no encapsulation or separation of concerns. Things went wherever felt most convenient. Naturally, this got messy fast and I had to refactor the whole thing into separate packages. Now, I'm even learning about and implementing **dependency injection** to clean up the project organization even further. This concept really excites me so far. It makes the system so much easier to visualize in my head and I'm excited to work in this codebase after I'm done refactoring it to use dependency injection.

### 3. Technologies

This project also introduced me to two very useful technologies, and gave me a little bit of hands-on experience.

#### 1. Docker & Compose

While not used in *this* repo, a rewrite I was working on earlier used Docker Compose to bundle the server's binary, MySQL, and Redis into one neat container. I am definitely not proficient in using Docker yet, but the exposure feels invaluable and I'm no longer intimidated by it.

#### 2. Redis

My current implementation of the leaderboard is clunky. I query the SQLite database with `SELECT user_id, display_name, pets FROM users ORDER BY pets DESC LIMIT 10`. This won't scale well and I needed to find a solution. That's when I learned about Redis `ZSets`. I got to play around with these a little and I'm excited to properly implement them into this repo ASAP.

### 4. Testing

Prior to this project, I would only test code by manually playing around with it. Now I'm getting some hands-on experience with writing tests.

## What I've Built

- A Realtime WebSocket server
- Redis & MySQL implementations (separate repo)
- Fully testable DB layer
- Clear requirements documentation
- Dataflow diagrams & ERDs

## Skills I've Developed
- Stronger understanding of concurrency
- Thinking in systems
- Writing extensible code
- Project Organization
- Debugging & Testing

