
# The North Star

After years as a software engineer, I’ve realized every developer needs a guiding philosophy. For me, a good application is one that:

1.  Works reliably in the present.
2.  Can evolve from simple to complex as requirements grow.

### Working Reliably in the Present
Reliability has two parts: first, making the application work, and second, making it work well. A TDD mindset is essential here.

TDD isn’t just about unit tests; it’s a mindset where you proactively anticipate failure points. This involves building an observer/supervisor system to provide feedback at every step, catching issues quickly. Then, you refactor to reduce technical debt, making the code readable and loosely coupled for future enhancements.

This mindset applies to various levels:
- **Daily Coding**: the first thing you need to learns is writing test cases (unit, integrate, e2e test) to have a feedback system, then "Make the test work quickly, committing whatever sins necessary in the process." as Kent Beck said. Finally reduce technical debt to open the ability for extending by refactoring the code|test, make it readable, loosely coupled.
- **Application Development**: the first thing we need to do is implement Fault Injection & Observer System (Log, Trace, Metric) to have a feedback system, then we need to implement resilient techniques like Rate Limit, Retry, Timeout, Circuit Breaker, Bulk Head. Finally, we need to separate those implemented from the business logic code, make it replaceable.
- **Deployment System**: the first thing is to introduce chaos testing tools like Chaos Monkey to verify resilience.
- **Agile Management**: the first thing we need is the Acceptance Criteria.
- **Team Management**: the first thing is to establish a team charter to set shared values and goals.

### Having the ability to evolve
For me, the best software architecture is the one that simple enough to satisfy the current requirement and can evolve into a complex system if needed. We need not start a project with fancy stuff like microservice, event-driven architecture, CQRS or eventsourcing... on the simple problem. If possible, we should start with a loosely coupled DDD Modular Monolith (based on context) to make sure our solution is working, when the time comes we can break it into services, choose the suitable communication method between services, by then our architecture will naturally evolve to what it needs to be, don't force it. Hexagon|Onion|Clean|Micro-Kernel architecture is quite fit for this problem.  

In conclusion, we need The North Star (Observer/Feedback system) to guide us and be prepared for the next challenges (change requirements).
  
# DOWN THE ROAD
With those in mind, I choose TDD & Clean Architecture & some design patterns to set up a simple shortened URL service to see how good it can be. The requirement is simple: we input an URL and get a short link that can redirect to the original URL. Ex: https://vnexpress.net => http://localhost/xjfepj

## Project Structure
I started with Modular Monolith in which each module should be defined by using the DDD approach (Module = Bounded Context).

Each module will have the same folder structure which is implemented in Clean Architecture (including config, container) so that we can easily break them into services later.

![](https://i.imgur.com/59leWsj.png)
  
## Loose coupled and highly cohesive

In this application, I purposely not using any libraries/framework at the beginning so we can have a better control on the project structure. Only after the whole application structure is laid out, I will consider replacing some components of the application with libraries.  

Some benefits of using Clean Architecture:

> 1. Independent of Frameworks. The architecture does not depend on the existence of some library of feature software. This allows you to use such frameworks as tools, rather than having to cram your system into their limited constraints.
> 2. Testable. The business rules can be tested without the UI, Database, Web Server, or any other external element.
> 3. Independent of UI. The UI can change easily, without changing the rest of the system. A Web UI could be replaced with a console UI, for example, without changing the business rules.
> 4. Independent of Database. You can swap out Oracle or SQL Server, for Mongo, BigTable, CouchDB, or something else. Your business rules are not bound to the database.
> 5. Independent of any external agency. In fact, your business rules simply don’t know anything at all about the outside world.

Thanks to the benefit of Clean Architecture the way a service communicates (REST, gRPC, MQ...) without touching the business rules. When we want to add or change the way to communicate to the outside we just need to add a new handler/client in the 'interface' package. This is the interface between our application and outside data service, Ex: another gRPC service. All the data conversion and transformation happened here, so the business logic code doesn’t need to be aware of the detailed implementation (whether it is gRPC or REST) of outside services. For example:

![](https://i.imgur.com/XUhsWTx.png)
  
![](https://i.imgur.com/l0tyT60.png)

## The Glue
To glue all loosely coupled components together as well as make them easier to be changed without touching existed code. I added the configuration and container to the project  

**Configuration** act as a blueprint in which we define which database we'll use and used for which use-case, we can also config which logger (zap, logrus) we want to use. And then **Container**, act as a factory, take in the blueprint and produce concrete instances (use case, logger, repository...)  

For now, configuration and container are implemented as simple as possible. In the future, the configuration could be read from a configuration server, support dynamic reloading of application configurations from a configuration server, a must-have feature in microservice. If our application has a lot of types and complex dependency relationships among the types, then probably we can switch to Dig or Wire library, otherwise, stay with the current solution.

With the combination of configuration and container, we seldom change any existing code (except for the container code), but only add new code to reduce QA's workload.

**Config database**
![](https://i.imgur.com/BqAhOGR.png)

**Config logger**
![](https://i.imgur.com/RwH1PjE.png)
  
## Watcher
To have an observer/feedback system, I add test cases for every file belong to business code and implement logging, tracing, resilient techniques like rate limiting, retry, timeout, circuit break to make sure the service work consistently.  

**Implement Rate limiting on server**
![](https://i.imgur.com/4pUIMcX.png)

**Implement Timeout on both server and client**
![](https://i.imgur.com/tCWTopW.png)

**Implement Retry on client**
![](https://i.imgur.com/Ne0dkFV.png)

**Implement Circuit Breaker on client**
![](https://i.imgur.com/dTvYPOD.png)

In the future, we can use Netflix’s Hystrix which integrates both the bulkhead isolation technology and the Circuit breaker to achieve isolation by restricting access to a service’s resources (typically Thread). And when we have many functions to implement resilient techniques, it's time for Service Mesh to come to the rescue. Because most of the problems and solutions of service resilience are related to infrastructure, it is better to leave them to infrastructure rather than code. This should relief application code from those burdens and focus back on business logic.  Now we have extracted those features out of application code and passed them to Service Mesh, of which the popular ones are Istio and Linkerd. By manipulating service requests, Service Mesh gained granular control of applications, while on the contrary, the container can only control the service level.

# CONCLUSION

1. The application is isolated from frameworks. So the frameworks/libraries won’t take over the application and we decide when and where to use them.
2. Technical changes are separated from business changes. Business logic code is never touched when making the above changes, via vera. Ex: switch to a better logger/tracer, change database handler (MySQL, MongoDB, Redis ...), change communication method (REST, gRPC..) ...
3. Loose coupled and highly cohesive
4. Easy to write tests (unit, integration, e2e test).
5. Open-closed principle: seldom change any existing code (except for the container code), but only add new code to reduce QA's workload.
6. Support multiple databases ( SQL and NoSql database) on the data persistence layer
7. Support data coming from other Microservices using different protocols such as gRPC or REST
8. Support easy and consistent logging and be able to change it ( for example, logging level and logging provider) without modifying logging statements in every file.
9. Make application/service resilient.
10. Have room to evolve.
