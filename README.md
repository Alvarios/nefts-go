# N1QL Easy Full Text Search (Go Edition)

A tool for an efficient configuration of N1QL full text search queries.

## Summary

- [About](#about)
- [NEFTS usage](#nefts-usage)
- [Options](#options)
    - [Config](#config)
        - [Cluster](#cluster)
        - [Bucket](#bucket)
        - [Parameters](#parameters)
    - [Fields](#fields)
    - [Where](#where)
    - [Joins](#joins)
        - [Join parameters](#join-parameters)
        - [Example explanation](#example-explanation)
        - [JoinQuery](#joinquery)
    - [Labels](#labels)
            - [Label Example](#label-example)
            - [Aliases](#aliases)
            - [Nested fields](#nested-fields)
            - [Labels with joins](#labels-with-joins)
                - [Labels on computed tuple](#labels-on-computed-tuple)
                - [Use the Bucket option](#use-the-bucket-option)
    - [LabelOptions](#labeloptions)
        - [Analyzer](#analyzer)
        - [Fuzziness](#fuzziness)
        - [Out](#out)
        - [Weight](#weight)
        - [Bucket (LabelOptions)](#bucket-labeloptions)
        - [PhraseMode](#phrasemode)
        - [RegexpMode](#regexpmode)
    - [Order](#order)
    - [QueryString](#querystring)
- [Results](#results)
- [Error handling](#error-handling)
- [Copyright](#copyright)

## About

N1QL Easy Full Text Search (NEFTS) allow to configure a N1QL query from an 
object oriented perspective. It aims to complete the efficiency of N1QL string
queries in a more configurable way.

N1QL string queries are pretty efficient, however they are very specific. For
Full Text Search, it can sometimes be tricky to write, mostly when they are
various parameters to consider. You want to have big strings, with many
variable parts that switch on and off depending on some condition.

NEFTS is useful in case of advanced full text search, where user can add
labels (`key:variable`) for more precise search, operators (such as OR, AND, ...),
and so on. In that situation, FTS may occur with a highly variable configuration,
leading the N1QL search string generation to be pretty hard.

NEFTS was created to make this configuration way more easy.

## NEFTS usage

NEFTS has a main `Thread()` function to use for your queries. It takes the
following parameters :

```go
queryResults, err := nefts.Thread(start, end, options)
```

Each parameter is required.

| Parameter | Type | Description |
| :--- | :--- | :--- |
| start | uint64 | Index of the first element to retrieve, in the result set. |
| end | uint64 | Index of the last element to retrieve, in the result set. |
| options | [nefts.config.Options](#options) | Options to pass to your query. The difference with config is that options parameters may be controlled by client and vary on each request. |

This function returns two results :

| Return argument | Type |
| :--- | :--- |
| queryResults | [nefts.config.QueryResults](#results) |
| err | [nefts.config.Error](#error-handling) |

## Options

Full example:
```go
options := nefts.config.Options{
    Config: nefts.config.Config{
        Cluster: myCouchbaseCluster,
        Bucket: "myBucket",
    },
    Fields: []string{"*", "meta.id"},
    Where: []string{fmt.Sprintf("b.last_update < %q", timestamp)},
    Joins: []nefts.config.Join{
        Join{
            Bucket: "users",
            ForeignKey: "authorId",
            DestinationKey: "author",
            Fields: ["username", "avatar_url", "sex"],
        },
        Join{
            Bucket: "answers",
            ForeignKey: "META(b).id",
            DestinationKey: "answers",
            Fields: ["*"],
            JoinKey: "postId"
        },
        Join{
            Bucket: "users",
            ForeignKey: "authorId",
            DestinationKey: "answers",
            Fields: ["username", "avatar_url", "sex"],
            ForeignParent: "answers"
        },
    },
    Labels: map[string][]string{
        "genre": ["g", "genre"],
        "user.username": ["username"],
    },
    LabelOptions: map[string]nefts.config.LabelOption{
        "genre": LabelOption{
            Analyzer: "standard",
            Fuzziness: "2",
            Out: "{key}_query_score",
            Weight: "0.8",
            PrefixLength: "0",
            PhraseMode: false,
            RegexpMode: false,
        },
        "user.username": LabelOption{
            Weight: "1",
        },
    },
    Order: map[string]string{
        "last_update": "desc"
    },
    QueryString: "Borderlands 3 cartels event secret puzzle author:lilux",
}
```

### Config

First, create a configuration object. You can use a specific configuration
for each of your methods, depending on your needs.

Minimal example:

```go
Config : nefts.config.Config{
    Cluster: myCouchbaseCluster,
    Bucket: "myBucket",
}
```

Full example:

```go
Config : nefts.config.Config{
    Cluster: myCouchbaseCluster,
    Bucket: "myBucket",
    Parameters: nefts.config.Parameters{
        MaxQueryLength: 1000,
        Debug: false,
    },
}
```

#### Cluster

This is the only required parameter. Pass it the Couchbase Cluster you set
up for your application (see this [setup guide from couchbase](https://docs.couchbase.com/go-sdk/2.1/hello-world/start-using-sdk.html)).

#### Bucket

The bucket in which to perform the search. NEFTS currently doesn't support
buckets protected by an individual password.

#### Parameters

Global parameters for the function.

| Parameter | Default | Description |
| :--- | :--- | :--- |
| MaxQueryLength | 1000 | Cap the client query string length to a maximum value. For no maximum, set it to a negative value. |
| Debug | false | Print the generated N1QL string to log. |

### Fields

A list of fields to filter in the output object. Defaults value is `[]string{"*", "meta.id"}`.

> 💡 Tip : * selects all non meta fields. To select meta fields such as document
id, use the 'meta.' prefix.

### Where

Write additional filter conditions in N1QL format.

### Joins

Join tables in the query.

#### Join parameters

A minimal join contains the following parameters:
```go
Join{
    Bucket: "users",
    ForeignKey: "authorId",
    DestinationKey: "author",
    Fields: ["username", "avatar_url", "sex"],
}
```

| Parameter | Type | Default | Required | Description |
| :--- | :--- | :--- | :--- | :--- |
| Bucket | string | - | true | The bucket to join. Has to be in the same Cluster. |
| ForeignKey | string | - | true | Foreign key in the current table referring to the distant bucket. |
| DestinationKey | string | - | true | Key in the resulting tuple in which the joined information will be inserted. |
| Fields | []string | - | true | Fields to filter in the distant document. **(1)** |
| JoinKey | string | "meta.id" | - | Reference key in the distant document. |
| ForeignParent | string | - | - | Specify a reference bucket, different from the current one. Useful for nested joins. |

**(1)** You can select all fields using either "all" or "*" selectors. If you
do so, they will be nested in a data key. For example :

```go
Join{
    Bucket: "users",
    ForeignKey: "authorId",
    DestinationKey: "author",
    Fields: ["*", "meta.id"],
}
```

will generate the following key in the tuple sent to client :

```json
{
  "author": {
    "id": "my_user_id",
    "data": {...}
  }
}
```

#### Example explanation

Let's look back at the join we performed. It is a pretty good example about
the possibilities of a configurable join.

```go
Joins: []nefts.config.Join{
    Join{
        Bucket: "users",
        ForeignKey: "authorId",
        DestinationKey: "author",
        Fields: ["username", "avatar_url", "sex"],
    },
    Join{
        Bucket: "answers",
        ForeignKey: "meta.id",
        DestinationKey: "answers",
        Fields: ["*"],
        JoinKey: "postId"
    },
    Join{
        Bucket: "users",
        ForeignKey: "authorId",
        DestinationKey: "answers",
        Fields: ["username", "avatar_url", "sex"],
        ForeignParent: "answers"
    },
}
```

Our database may have, in this case, 3 document schemas :

*posts*
```json
{
  "title": "string",
  "content": "string",
  "authorId": "userId"
}
```

*answers*
```json
{
  "title": "string",
  "content": "string",
  "authorId": "userId",
  "postId": "postId",
  "score": "number"
}
```

*users*
```json
{
  "username": "string",
  "avatar_url": "url",
  "sex": "string"
}
```

We want to retrieve our post, all the answers, and link those answers to
their author.

> ⚠️ Warning : this example shows some bad practice. Please refer to the
[below section](#joinquery) for a better practice. Following example just serves as a
demonstration for some specific options.

```go
Join{
    Bucket: "users",
    ForeignKey: "authorId",
    DestinationKey: "author",
    Fields: ["username", "avatar_url", "sex"],
}
```

The first join is pretty easy to understand. Our post data holds the id of
its author. But if we want some more information about her or him, we need
to search in the users table.

Here, we tell NEFTS to look for the document in users table with the ID
provided by authorId in our document. Then, if found, we want to add the
username, avatar_url and sex fields to the author key in our result tuple.

```go
Join{
    Bucket: "answers",
    ForeignKey: "meta.id",
    DestinationKey: "answers",
    Fields: ["*"],
    JoinKey: "postId"
}
```

Now, we want to retrieve every answer related to our post. It is almost the
same structure, instead we provide a JoinKey rather than a ForeignKey.

In this scenario, the link information is not holded by our post, but by the
children, the answers. Instead of updating a large array referencing every
answer id, we provide the answers an id for their parent post.

So the request is : select every answers document where postId match our current
document id.

```go
Join{
    Bucket: "users",
    ForeignKey: "authorId",
    DestinationKey: "answers",
    Fields: ["username", "avatar_url", "sex"],
    ForeignParent: "answers"
}
```

Finally, we want to retrieve the same author information for each answer as
we did for the post. So the join configuration is the same as the first one,
except we specify a ForeignParent bucket.

#### JoinQuery

The above example is a bad practice. If our post has thousands of answers,
they will all be retrieved at once, leading to poor performances and
eventually crash of the server.

While above parameters can be used for limited sets of joined data, our
case would rather require, for example, to only retrieve the top post,
based on some custom criterias.

Here comes the power of N1QL joins : a Join clause can hide a full nested query
inside !

Let's take again the above example, and this time try to only join the most
rated post. Since we don't perfom nested FTS, we will go with minimal
configuration :

```go
config := nefts.config.Config{
    Cluster: myCouchbaseCluster,
    Bucket: "answers",
}
```

Let's now build a minimal Options parameter :
```go
options := nefts.config.Options{
    Joins: []nefts.config.Join{
        Join{
            Bucket: "users",
            ForeignKey: "authorId",
            DestinationKey: "author",
            Fields: ["username", "avatar_url", "sex"],
        },
    },
    Order: map[string]string{
        "score": "desc"
    },
}
```

Our query will just look up for the answers, ordered by descending score,
and link them with their author.

We just want the most rated post. So let's just set the limit to 1.
But we don't actually want to run a full thread calls : we just need to build
the nested query string. Then, for better performance, run the whole thing in
one single call. This is where the JoinQuery option comes.

JoinQuery takes every NEFTS Thread arguments as parameters. We can now
transform our original declaration:

```go
baseOptions := nefts.config.Options{
    Fields: []string{"*", "meta.id"},
    Where: []string{fmt.Sprintf("b.last_update < %q", timestamp)},
    Joins: []nefts.config.Join{
        Join{
            Bucket: "users",
            ForeignKey: "authorId",
            DestinationKey: "author",
            Fields: ["username", "avatar_url", "sex"],
        },
        Join{
            Bucket: "answers",
            ForeignKey: "META(b).id",
            DestinationKey: "answers",
            Fields: ["*"],
            JoinKey: "postId",
            JoinQuery: &nefts.config.JoinQuery{
                Config: config,
                Options: options,
                Start: 0,
                End: 1,
            }
        },
    },
    Order: map[string]string{
        "last_update": "desc"
    },
    QueryString: "Borderlands 3 cartels event secret puzzle author:lilux",
}
```

And we're good ! Now our result will hold an array in the "answers" key,
with only the most rated post.

> 💡 Tip : JoinQuery key holds a pointer.

### Labels

Add optional label marks to your client query string. Label mark restricts a
term or a group of terms to a specific subset of your documents.

#### Label Example

Let's take an example : say you have the following document in your database.
```json
{
  "title": "string",
  "description": "string",
  "labels": [
    {
      "labelName": "string",
      "labelCategory": "string"
    }
  ],
  "authorId": "string id"
}
```

Without any consideration for the below sections, a query will look for every
field in each document.

Now let's say you want to add the possibility to restrict some terms only to
a specific field. For example, add a filter that only applies to labels. An
easy way for your user to do it would be :

`my awesome search label:cooking`

Thus, all terms will look up for the whole document, except 'cooking', which
is only going to match the labels field.

> 💡 Tip : NEFTS supports usage of quotation marks to keep terms linked together.
    For example : `label:cooking tips` will only insert `cooking` in the label
    mark, but `label:"cooking tip"` will keep both terms.
>
> Note NEFTS is very permissive with quotation marks. `lab"el:cooking t"ip`,
    `label:cookin"g tip"` or `label:cooking" "tip` will all produce the same
    result as above example.

Let's take a look at how to implement this behavior. The configuration for it
would be :

```go
Labels: map[string][]string {
    "labels": ["label"]
}
```

Thus, you tell NEFTS that for this particular query, the `label:` mark is
only going to look inside the labels fields.

#### Aliases

In the configuration map, key will point to the actual field in the fetched
tuples, and value is an array of aliases. Thus, you can assign multiple marks
to the same field, for example, a long version and a shorter one :

```go
Labels: map[string][]string {
    // "l:" and "label:" will both point to the "labels" field.
    "labels": ["label", "l"]
}
```

> 💡 Tip : Since label marks are variables, you can also add localization (different
        labels depending on the client language). There are many use cases here.

#### Nested fields

You can also specify a label for a nested field. For example, if you want
the above example to target only the `labelName` field. You just have to
change the configuration to :

```go
Labels: map[string][]string {
    // "l:" and "label:" will both point to the "labels" field.
    "labels.labelName": ["label", "l"]
}
```

> 💡 Tip : Nested field declaration works for both nested fields and arrays (as in the
        above example).

#### Labels with joins

You may want to perform FTS on some joined buckets. There are two methods.

##### Labels on computed tuple

By default, labels will filter the tuple produced by the query.

Let's take again the document example:

```json
{
  "title": "string",
  "description": "string",
  "labels": [
    {
      "labelName": "string",
      "labelCategory": "string"
    }
  ],
  "authorId": "string id"
}
```

Now, you may want to send some information about the author, which are stored
in another bucket. Let's perform the following join (more details in the
[joins section](#joins)) :

```go
// More details in the options section.
options := nefts.config.Options{
    //...
    Joins: []nefts.config.Join{
        Bucket: "users",
        // The important parameter.
        ForeignKey: "authorId",
        // All joined data will be stored under the 'author' key in results.
        DestinationKey: "author",
        Fields: ["*"],
    }
}
```

Your user document looks like this :
```json
{
  "username": "string",
  "reputation": 1000,
  "description": "string",
  "achievements": {
    "trophies": []
  }
}
```

The above join will produce the following result:

```json
{
  "title": "string",
  "description": "string",
  "labels": [
    {
      "labelName": "string",
      "labelCategory": "string"
    }
  ],
  "authorId": "string id",
  "author": {
    "data": {
      "username": "string",
      "reputation": 1000,
      "description": "string",
      "achievements": {
        "trophies": []
      }
    }
  }
}
```

And that's the magic. By default, label will filter the joined tuple rather
than the documents in table. So knowing this, you can simply use :

```go
Labels: map[string][]string {
    // Match author username.
    "author.data.username": ["author", "a"]
}
```

##### Use the Bucket option

This option is part of below [LabelOptions](#labeloptions) section.

The con of above methods is it will ignore fields that are filtered in the
result. You may want, for example, to search among user with a specific
achievement, but don't need to send achievements back to client.

Bucket option will specify in which bucket to apply the filter. Keep in
mind that Couchbase implementation of FTS will only attribute a match score,
so you don't need the field to be actually returned.

For example :
```go
LabelOptions: map[string]nefts.config.LabelOption{
    "author.data.achievements": nefts.config.LabelOption{
        Bucket: "users"
    }
}
```

So in case you don't send the achievements to the client, you can still
take it in account for your search. More details are provided in the
[below section](#labeloptions).

### LabelOptions

You can fully control a specific FTS operation with label options. Just add
a key value pair to this parameter : key is going to tell NEFTS which
label key to apply the options to, and value lists those options.

> 💡 Tip : you can also control options for the non labelled terms. Non labelled
        terms are actually considered as labelled under the "general" key.
        Use this key to control them.

> 💡 Tip : to apply some parameter to multiple labels, use the "global" key. Thus,
        each non specified key will comply to "global" key parameters.

Let's take a look at a typical set of options :

```go
options := LabelOption{
    Analyzer: "standard",
    Fuzziness: "2",
    Out: "{key}_query_score",
    PrefixLength: "0",
    Weight: "0.8",
    Bucket: "b",
    PhraseMode: false,
    RegexpMode: false,
}

//...

LabelOptions: map[string]nefts.config.LabelOption{
    "my.key": options
}
```

#### Analyzer

Analyzers are provided by Couchbase to attribute match score on a FTS query.
More details can be found in [official documentation](https://docs.couchbase.com/server/6.5/fts/fts-using-analyzers.html).

Couchbase provides 4 prebuilt analyzers:

**keyword**: Creates a single token representing the entire input, which is otherwise unchanged. This forces exact matches, and preserves characters such as spaces.

**simple**: Analysis by means of the Unicode tokenizer and the to_lower token filter.

**standard**: Analysis by means of the Unicode tokenizer, the to_lower token filter, and the stop token filter.

**web**: Analysis by means of the Web tokenizer and the to_lower token filter.

By default, NEFTS uses "standard" analyzer.

#### Fuzziness

From [Couchbase documentation](https://docs.couchbase.com/server/6.5/fts/fts-query-types.html#match-query):

When fuzzy matching is used, if the single parameter is set to a non-zero 
integer, the analyzed text is matched with a corresponding level of 
fuzziness. The maximum supported fuzziness is 2.

#### Out

Specify an alias for the match score. This score is returned to the results.
By default, each score is returned as `{key}_query_score`.

> 💡 Tip : When naming score, {key} will be replaced by the current label key
    (or "general" for non labelled terms). If multiple labels are provided,
    it is highly recommended to use the {key} syntax in your score naming,
    to avoid unexpected behaviors.

PrefixLength

From [Couchbase documentation](https://docs.couchbase.com/server/6.5/fts/fts-query-types.html#match-query):

When a prefix match is used, the prefix_length parameter specifies that for
a match to occur, a prefix of specified length must be shared by the 
input-term and the target text-element.

#### Weight

Attribute a weight to balance each match score. A higher score means a higher
importance for a given label.

#### Bucket (LabelOptions)

Specify a Bucket to search for labelled terms. This is useful when some fields
are filtered in a document, but still need to be taken in account in the query.

#### PhraseMode

From [Couchbase documentation](https://docs.couchbase.com/server/6.5/fts/fts-query-types.html#match-phrase-query):

The input text is analyzed, and a phrase query is built with the terms 
resulting from the analysis. This type of query searches for terms in the
target that occur in the positions and offsets indicated by the input: this 
depends on term vectors, which must have been included in the creation of the
index used for the search.

For example, a match phrase query for location for functions is matched with 
locate the function, if the standard analyzer is used: this analyzer uses a 
stemmer, which tokenizes location and locate to locat, and reduces
functions and function to function. Additionally, the analyzer employs stop
removal, which removes small and less significant words from input and target
text, so that matches are attempted on only the more significant elements of
vocabulary: in this case for and the are removed. Following this processing,
the tokens locat and function are recognized as common to both input and
target; and also as being both in the same sequence as, and at the same
distance from one another; and therefore a match is made.

#### RegexpMode

Consider the query string as a regular expression for search. It has to follow
the [Go regexp syntax](https://golang.org/pkg/regexp/syntax/).

### Order

Key-Value pairs. Key is the field, and value is either "desc" or "asc".

### QueryString

The string to perform FTS from.

## Results

Results is a custom object that contains useful information to help you with
threading. It is of type `nefts.config.QueryResults`.

| Key | Type | Description |
| :--- | :--- | :--- |
| Results | []interface{} | An array of json objects. It contains the tuples generated by query. |
| Boundaries | Boundaries | Hold information about the actual limits of the returned set. |
| Boundaries.start | uint64 | Low boundary of the returned set. |
| Boundaries.end | uint64 | High boundary of the returned set. |
| Flags | Flags | Indication about the last request status. |
| Flags.BeginningOfResults | boolean | No data left below the returned set. |
| Flags.EndOfResults | boolean | No data left after the returned set. |

## Error handling

NEFTS returns a pointer to an error object adapted to web servers.

| Key | Type | Description |
| :--- | :--- | :--- |
| Status | int | The http status of the error. |
| Message | string | Describes the nature of the error. |

## Copyright
2020 Kushuh - [MIT license](https://nefts/blob/master/LICENSE)
