# buffStreams

  This is  a solution to the tech challenge around API design and implementation

## Tech Choices
  
### Config
Using [Viper](https://github.com/spf13/viper) for config. Values can be overridded by using either a config.json file in the same folder as the exe or through environment variables. Prepend API_ to the variable name below when using environment variables

Config options include
| Field | Default | Description |
|--|--|--|
| DBUser | root | User for the db |
| DBPassword | password | Password for the DB |
| DBHost | 127.0.0.1:3310 | Host:IP for the DB |
| DBSchmea |buffup | Name of the DB Schema
| DBArgs| "multiStatements=true&parseTime=true" | Additional Args to pass to DB Connection setup |
| PageSize | 10 | This is default page size for API results. Eventually the users would be able to specify their own |

## Setup

  ### Docker compose
  Run `docker-compose up` to start the services
  

## Testing
As well as inbuilt go unit tests I've included some functional tests using postman.
You can download the app or use newman which is a CLI to run postman. The collection is included in the test folder

  

## Loading Sample Data
Run the load.sql file from the test folder via a MYSQL client


## Decisions/Assumptions
Just a few topics of interest for the solution to help explain some reasoning

 - DB modelling for questions. I'm not sure on this tbh. I had initially thought the most straight-forward thing was adding answer columns to question table but from looking at the site etc there seems to be other question types (like Polls) which will come around sooner than later so didn't want to have to restructure DB too heavily to handle them. One downside is enforcing data integrity (ie that there's a correct answer etc but that could be done at ingest time) another is when it comes to validating answers but I think we would use redis or similar to cache answers for quick retrieval anyway
 - Spec mentioned Guests a lot so I've not included auth. I would normally use a JWT or similar
 - I've not allowed user specified sorting or filtering on query params just specifying a page number
 - I've not modelled state of streams but I assume that would be coming soon
 - Choice of MYSQL. Just because I'm used to it and thought it could model the relationships between questions and streams well
 - Link Table. This is to support a many to many relationship. I'm assuming questions could/should be re-used between streams so wanted to give that flexibility
 - Postman : I like postman for API testing


 ## Next Steps
If I had more time I would look to
- Finish tests. More controller tests. Lot of this is covered by postman but nice to get instant feedback
- I hate the fact the correct answer is always last. I would maybe add a function to random the order so someone cant use dev tools and cheat basically. It's why I don't reference a 'correct' answer
- Test larger dataset and use Redis or similar for caching as stream data should change quite slowly
- Benchmark tests
- Replace int array for questions with a HATEOS style URL and integer