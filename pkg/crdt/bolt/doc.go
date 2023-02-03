package bolt

/*
	This package implements NodeStore/TopicStore using boltdb.

	We're using go.etcd.io/bbolt because the original github.com/boltdb/bolt is archived
	and has a -race error https://github.com/golang/go/issues/54690

	The DB bucket structure is as follows:

	META
	  - Version -> 0x0100  // DB version
	Topic1
	  EVENTS
	    - CID -> Event
	  HEADS
	    - CID -> nil
	  BY_TIME
	    - ByTimeKey -> CID
	Topic2
	Topic3
	...

	At the top level there is a META Bucket and a Bucket for each topic.
	The keys for the topic Buckets are the topic names.
	The META bucket is for any kind of meta information we want to keep in the DB.
	Currently the only key there is Version.

	In each topic Bucket there is EVENTS, HEADS and BY_TIME bucket.
	EVENTS bucket contains marshalled Events, keyed by the CIDs
	HEADS bucket contains CIDs of current DAG head Events as keys with empty values.
	BY_TIME bucket is an index keyed by ByTimeKeys of Events, the values are the Event CIDs.


	The message API query cursor, i.e the messagev1.IndexCursor, is interpreted as follows:
		cursor.SenderTimeNs = the event TimestampNs
		cursor.Digest = the event CID
	It points at the last Event from the previous page of the query result.
	The returned cursor is nil if this is the last page of the result.
	Only queries with non-zero limit parameter can yield a cursor.

*/
