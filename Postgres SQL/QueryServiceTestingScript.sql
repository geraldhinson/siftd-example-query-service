-- unittests table with a variety of data types

drop table if exists "PG-KitchenSink"

create table if not exists "PG-KitchenSink" 
(
 "aBigInt" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY,
 "aVarchar" varchar(50) NOT NULL,
 "aJson" json NOT NULL,
 "aFloat4"	float4 NOT NULL,
 "aFloat8"	float8 NOT NULL,
 "anInteger"	integer default 1,
 "aTimestamp" timestamp  NOT null,
 "aTimestamp-NoTimezone" timestamp without time zone NOT null,
 "aTimestamp-WithTimezone" timestamp with time zone NOT null,
 "aDate"	date NOT null,
 "aGuid"	UUID not null,
 "aBoolean"	boolean not NULL,
 "anIntegerArray" integer ARRAY NOT NULL,
 "aDateArray" date ARRAY NOT NULL,
 "aVarcharArray" varchar ARRAY NOT null,
 "aNullableVarchar" varchar,
 "ownerId" varchar(50)
)

insert into public."PG-KitchenSink"
("aVarchar", "aJson", "aFloat4", "aFloat8", "aTimestamp", "aTimestamp-NoTimezone", 
 "aTimestamp-WithTimezone", "aDate", "aGuid", "aBoolean", "anIntegerArray", "aDateArray", "aVarcharArray", "ownerId")
values
(
	'12345678901234567890123456789012345678901234567890',
	'{"menu":{"id":"file","value":"File","popup":{"menuitem":[{"value":"New","onclick":"CreateNewDoc()"},{"value":"Open","onclick":"OpenDoc()"},{"value":"Close","onclick":"CloseDoc()"}]}}}',
	3.14,
	3.14,
	now(),
	now(),
	now(),
	now(),
	gen_random_uuid(),
	true,
	array[1, 2, 3, 5],
	array[now(), now()  + INTERVAL '1 DAY', now() + INTERVAL '2 DAY'],
	array['abc', 'def'],
	'GUID-fake-member-GUID'
)

select * from public."PG-KitchenSink" where "ownerId" = '1dc478eb-a76a-4ddb-a37f-9c4deeb1e768'
select * from public."PG-KitchenSink" pks
select * from public."PG-KitchenSink" where "aGuid" = '1dc478eb-a76a-4ddb-a37f-9c4deeb1e768'
select * from public."PG-KitchenSink" where "aTimestamp" < to_timestamp('2025-02-26 12:20:02', 'YYYY-MM-DD HH24:MI:SS')
select * from public."PG-KitchenSink" where "aTimestamp" between '2025-02-26 12:20:01.385' and '2025-02-26 12:20:01.386'
select * from public."PG-KitchenSink" where "aDate" in ('2025-02-26', '2025-02-27')
select * from public."PG-KitchenSink" where "anIntegerArray" <@ ARRAY[2,1,3,4]
select * from public."PG-KitchenSink" where "anIntegerArray" @> '{2,1,3}'
select * from public."PG-KitchenSink" where "anIntegerArray" && '{2,1,3,4}'
select * from public."PG-KitchenSink" where "aVarcharArray" @> '{"abc", "def"}'
select * from public."PG-KitchenSink" where '{"2025-03-01", "2025-02-28", "2025-02-27"}' <@ "aDateArray"  
select * from public."PG-KitchenSink" where '{"2025-03-15"}' <@ "aDateArray"  


-- resource & journal creates
drop table if exists "Journal";
drop table if exists "Resources";
drop index if exists "IX_Resources_OwnerId";

CREATE TABLE IF NOT EXISTS "Journal" (
                                         "Clock" bigint NOT NULL GENERATED BY DEFAULT AS IDENTITY,
                                         "Resource" text NULL,
                                         "CreatedAt" timestamp without time zone NOT NULL,
                                         "PartitionName" varchar(20) NULL,
                                         CONSTRAINT "PK_Journal" PRIMARY KEY ("Clock", "PartitionName");
);

CREATE TABLE IF NOT EXISTS "Resources" (
                                           "Id" varchar(50) NOT NULL,
                                           "OwnerId" varchar(50) NOT NULL,
                                           "Version" integer NOT NULL,
                                           "CreatedAt" timestamp without time zone NOT NULL,
                                           "UpdatedAt" timestamp without time zone NOT NULL,
                                           "Deleted" boolean NOT NULL,
                                           "Resource" text NULL,
                                           CONSTRAINT "PK_Resources" PRIMARY KEY ("Id");
);

CREATE INDEX IF NOT EXISTS "IX_Resources_OwnerId" ON "Resources" ("OwnerId");

select * from public."Journal";
select * from public."Resources";

delete from public."Journal" where "Clock" = 55;
delete from public."Resources" where "Id" = '094b6128-35db-425b-be40-2bec672d26b5';

SELECT MAX("Clock") AS "Clock" FROM public."Journal";


