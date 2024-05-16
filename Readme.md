# URLBook
__URLBook__ is a sample project of URL-shortener platform with focus on backend side. 

## System design
### Functional requirement
- There should be an API (without OAuth action) that user can submit a valid URL for shortening.
- Shorted URL can have a system generated phrase or user customized value.
- The limitation for system generated phrase is up to 7 character
- The limitation for user customized phrase is up to 16 character
- The Shorted URL will not expire
- Shorted URL must redirect to the destination with status of 302
- With having shorted URL, can monitor the number of clicks and the device types which used for browse the link.

### Non-Functional requirement
- System should be highly available (99.9% uptime)
- 500k url submitted to the system per day
- 200M redirection request per month ( 200M / (30Day * 24h * 3600s) = 80tps )

### High-level Design Diagram
<div style='width: auto; max-width: 1000px; margin: 10px auto;'>
    <img src='docs/design.png' alt='high-level-system-design' />
</div>

### How project caching works?
As non-functional requirement noted, we will receive 500k url submission per day which for cache this amount of url
in about 5 years without expiration and average size of 200 characters, we might need 200GB-460GB of RAM (centralized or decentralized) 
to store and read the data from cache.

So inefficient, right! Let's think that the probability of redirection occurrence for a single shortened url decreased by the time passed
from its creation date, e.g. a short url which created at 15 minutes ago has 50% probability of occupying redirection traffic and a short url 
which created at last year has lower probability, like 0.5%.

Here are our strategy: 
 1. Cache a newly generated short url for plenty of time (like 3 hours)
 2. Set new cache for not cached short url on redirection
 3. Prepare a job to update caches in a reasonable sequence for last 6 months (which requires 18GB of RAM usage at most)

## Project architecture
As considered, _Hexagonal Architecture_ selected for project structure to make the project flexible and isolated its parts.

Here are some links about the architecture:
 - https://medium.com/@pthtantai97/hexagonal-architecture-with-golang-part-1-7f82a364b29 
 - https://medium.com/@pthtantai97/hexagonal-architecture-with-golang-part-2-681ee2a0d780
 - https://herbertograca.com/2017/11/16/explicit-architecture-01-ddd-hexagonal-onion-clean-cqrs-how-i-put-it-all-together/

## Todo
- [x] Create project system design
- [x] Setup required services with docker-compose
  - [x] Database (mysql)
  - [x] Cache (memcached)
- [x] Only submit a url and get a system generated short url
- [x] Redirect the system generated short-url to its original url with 302 HTTP status code
- [x] Submit a url with custom name for shortening
- [x] Bring caching mechanism
- [ ] Add configurable cache update job
- [ ] Add some tracking mechanism on urls
  - [ ] Number of clicks with date
  - [ ] The devices used to visit the link
  - [ ] etc.

