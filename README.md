# envir_collector_go
Data collector for Envir energy monitor ([Envi kit](http://www.currentcost.net/Monitor%20Details.html). This replaces the [python version](https://github.com/clinstid/energydash) I wrote a while back.

# Database schema
Unlike the original version from my [energydash](https://github.com/clinstid/energydash#data-collection) collector, I'm now using postgres as my database because the collected data fits very nicely into a traditional RDBMS. Mongo was fun to play with, but it definitely was not the right choice for this kind of data.

The data is currently stored in a single table called `energydash` (yeah, I kept the name):

| Column    | Description                                                                                |
|-----------|--------------------------------------------------------------------------------------------|
| src       | source and software version                                                                |
| dsb       | days since "birth"                                                                         |
| time      | 24 hour clock time (this is set by the collector rather than trusting the original source) |
| tmprf     | temperature in Fahrenheit                                                                  |
| sensor    | sensor id                                                                                  |
| device_id | ID for the device                                                                          |
| ch1_watts | Power reading on channel 1                                                                 |
| ch2_watts | Power reading on channel 2                                                                 |
| ch3_watts | Power reading on channel 3                                                                 |
