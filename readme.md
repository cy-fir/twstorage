
# Twitter Storage (experiment)

This is just a concept of using Twitter for a encrypted file storage, not a serious project or usage.

Premise: encrypt and encode file, upload twitter size chunks as replies to the first tweet. Download using Twitter URL and key to decode. Additionally, the first tweet could hold the metadata about the file, such as filename, size, ...

Inspired by @rauchg's Require from Twitter, and <a href="https://news.ycombinator.com/item?id=11352150">this comment</a> which linked to the slide: <a href="https://users.soe.ucsc.edu/~jpagnutt/Elevator_Pitch.pdf">https://users.soe.ucsc.edu/~jpagnutt/Elevator_Pitch.pdf</a>


### Setup

Create a <a href="https://apps.twitter.com/">new Twitter app</a> and create ~/.twstorage with:

```
{
  "ClientSecret": "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
  "ClientToken": "abcdefghijklmnopqrstuvwxyz"
}
```


### Build Tool using Go

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
   Key  : <key>
   Tweet: <url>
```

### Decrypt Usage:
```
twstorage -k <key> <url>
```


### Test Tweet

You should be able to test decoding using this Tweet and key:

```
Key  : Kgo0BK5JVatO0kDCsJmgT59nQqNaLgkS
Tweet: https://twitter.com/mkaz/status/713567397869817856
```

Decrypt using:

./twstorage -k Kgo0BK5JVatO0kDCsJmgT59nQqNaLgkS https://twitter.com/mkaz/status/713567397869817856


## Credits

* Twitter API code borrowed heavily from: https://github.com/mattn/twty
* Encryption from: https://gist.github.com/kkirsche/e28da6754c39d5e7ea10


