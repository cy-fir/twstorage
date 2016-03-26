
# Twitter Storage (experiment)

This is just a concept of using Twitter for a encrypted file storage, not a serious project or usage.

Premise: encrypt and encode file, upload twitter size chunks as replies to the first tweet. Download using Twitter URL and key to decode. Additionally, the first tweet could hold the metadata about the file, such as filename, size, ...

Inspired by a comment on @rauchg's require from Twitter HN post, which the following slide was posted:
https://users.soe.ucsc.edu/~jpagnutt/Elevator_Pitch.pdf


### Setup

Create a new Twitter app and create ~/.twstorage

```
{
  "ClientSecret": "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
  "ClientToken": "abcdefghijklmnopqrstuvwxyz"
}
```


### Build Tool

```
git clone github.com/mkaz/twstorage
go get
go build
```


### Encrypt Usage:

```
twstorage filename.txt
    - will prompt to authorize to twitter
    - stores credentials at ~/.twstorage

 Returns:
   Tweet: <url>
   Key  : <key>
```

### Decrypt Usage:
```
twstorage -k <key> <url>
```


### Test Tweet

You should be able to test decoding using this Tweet and key:

```
Tweet: https://twitter.com/mkaz/status/713567397869817856
Key  : Kgo0BK5JVatO0kDCsJmgT59nQqNaLgkS
```

Decrypt using:

./twstorage -k Kgo0BK5JVatO0kDCsJmgT59nQqNaLgkS https://twitter.com/mkaz/status/713567397869817856


### Credits

Twitter API code borrowed heavily from: https://github.com/mattn/twty

