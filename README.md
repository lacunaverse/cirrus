# cirrus

`quick and dirty entity recognition`

Cirrus attempts to serialize arbitrary natural text into a Go value using dictionary matching. At the moment it doesn't really do much, and shouldn't be used by anybody.

Cirrus tries to parse tokens into various values based on pretty simple heuristics like capitalization and the presence of a dollar sign. Some entities, like cardinality ("one", "many", "ten", etc) have to be matched against every token.

## Roadmap

+ grouping of sequential entities
  + e.g. sequential cardinals can be grouped, `two dozen` becomes `Result{Value: 24}`
+ text classification
  + use wordnet?

## Examples

+ `2021-10-22` is a `date`
+ `$20` is a `value` in `USD`
+ `20mph` is a `value` in `miles per hour`
+ `Hong Kong` is a `city`
+ `Charles Dickens` is a `person`
+ `Microsoft` or `MSFT` is a `company`
+ `Australia` is a `country`

### Numbers

```text
1.2
10e-3
1.00001
20E10
15
2/3
one
two thousand three hundred seventy five
```

### Dates

```text
Christmas eve
10/22/15
10.22.15
2005
december 5 2005 // variations thereof
```

### Times

```text
12:10
ten minutes and five seconds
```

### Phone Numbers

```text

```

### Addresses

### Email Addresses

+ probably best to use an established regexp for this

### Web Addresses

+ best to use an established regexp or std's `url.Parse`

### Tables

```text
| name | age |
| john | 2 |
| jane | 12 |

name,age
john,2
jane,12

name age
"john" 2
"jane" 12

name    age
john    2
jane    12
```
