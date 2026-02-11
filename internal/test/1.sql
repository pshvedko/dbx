DROP TABLE IF EXISTS public.objects;

CREATE TABLE public.objects
(
    id         uuid PRIMARY KEY                  DEFAULT gen_random_uuid(),
    o_uuid_2   uuid                     NOT NULL,
    o_uuid_3   uuid                              DEFAULT gen_random_uuid(),
    o_uuid_4   uuid                     CHECK ( id <> o_uuid_4 ),
    o_string_1 text                     NOT NULL DEFAULT 'green',
    o_string_2 text                     NOT NULL,
    o_string_3 text                              DEFAULT 'red',
    o_string_4 text,
    o_bool_1   boolean                  NOT NULL DEFAULT TRUE,
    o_bool_2   boolean                  NOT NULL,
    o_bool_3   boolean                           DEFAULT FALSE,
    o_bool_4   boolean,
    o_float_32 float4                   NOT NULL DEFAULT 0,
    o_float_64 float8,
    o_int_8    bit(8)                   NOT NULL DEFAULT 0::bit(8),
    o_int_16   int2                     NOT NULL,
    o_int_32   int4                              DEFAULT 0,
    o_int_64   int8,
    o_time_1   timestamp with time zone NOT NULL DEFAULT now(),
    o_time_2   timestamp with time zone NOT NULL DEFAULT now(),
    o_time_3   timestamp with time zone          DEFAULT now(),
    o_time_4   timestamp with time zone
);

create unique index if not exists objects_parent_idx
    on objects ((TRUE))
    where (uuid_4 IS NULL);

INSERT INTO public.objects (id, uuid_2, uuid_3, uuid_4, string_1, string_2, string_3, string_4, bool_1, bool_2, bool_3,
                            bool_4, float_32, float_64, int_8, int_16, int_32, int_64, time_1, time_2, time_3, time_4)
VALUES ('26dd056a-fddf-4992-9e7f-1af09d610c5c', '167b738c-d795-44d2-a8d6-2724e9fda43b', null, null, 'green', 'yellow',
        'red', null, true, true, false, null, 100, 3.14, '00000111', 16, 0, null, '2025-02-24 13:08:06.482425 +00:00',
        '2025-02-24 13:08:06.482425 +00:00', '1970-01-01 00:00:00.000000 +00:00', null);
INSERT INTO public.objects (id, uuid_2, uuid_3, uuid_4, string_1, string_2, string_3, string_4, bool_1, bool_2, bool_3,
                            bool_4, float_32, float_64, int_8, int_16, int_32, int_64, time_1, time_2, time_3, time_4)
VALUES ('4fc2afae-80b1-4045-845b-c373d180beda', 'f1824d97-6bc9-4f21-bed0-a09d3b4e80da', null,
        '26dd056a-fddf-4992-9e7f-1af09d610c5c', 'green', 'amber', null, null, true, true, false, null, 0, null,
        '00000111', 17, 0, null, '2025-02-24 15:39:52.238818 +00:00', '2025-02-24 15:39:52.238818 +00:00',
        '1970-01-01 00:00:00.000000 +00:00', null);
INSERT INTO public.objects (id, uuid_2, uuid_3, uuid_4, string_1, string_2, string_3, string_4, bool_1, bool_2, bool_3,
                            bool_4, float_32, float_64, int_8, int_16, int_32, int_64, time_1, time_2, time_3, time_4)
VALUES ('6dc2ed41-03ce-47e3-937e-6b4a7b159ec2', '411f8130-0230-4a7f-9f35-35c5af8e690c', null,
        '26dd056a-fddf-4992-9e7f-1af09d610c5c', 'green', 'white', 'red', null, true, true, false, null, 0, null,
        '00000111', 18, 0, null, '2025-02-24 15:40:16.175158 +00:00', '2025-02-24 15:40:16.175158 +00:00',
        '1970-01-01 00:00:00.000000 +00:00', null);
INSERT INTO public.objects (id, uuid_2, uuid_3, uuid_4, string_1, string_2, string_3, string_4, bool_1, bool_2, bool_3,
                            bool_4, float_32, float_64, int_8, int_16, int_32, int_64, time_1, time_2, time_3, time_4)
VALUES ('e8d3d9a4-9aba-4303-95bf-4dec730131e6', 'f7161ca6-4a03-4267-be30-99a2a5a78aa7', null,
        '26dd056a-fddf-4992-9e7f-1af09d610c5c', 'green', 'white', 'red', null, true, true, false, null, 0, null,
        '00000111', 19, 0, null, '2025-02-24 15:40:23.390318 +00:00', '2025-02-24 15:40:23.390318 +00:00',
        '1970-01-01 00:00:00.000000 +00:00', null);
INSERT INTO public.objects (id, uuid_2, uuid_3, uuid_4, string_1, string_2, string_3, string_4, bool_1, bool_2, bool_3,
                            bool_4, float_32, float_64, int_8, int_16, int_32, int_64, time_1, time_2, time_3, time_4)
VALUES ('62549d95-ba11-450b-a4a2-0911afd3a73f', '80240ad8-9eb8-4cae-ac20-f3e597ec61ba', null,
        '26dd056a-fddf-4992-9e7f-1af09d610c5c', 'green', 'black', 'red', null, true, true, false, null, 0, null,
        '00000111', 19, 0, null, '2025-02-24 15:40:37.743583 +00:00', '2025-02-24 15:40:37.743583 +00:00',
        '1970-01-01 00:00:00.000000 +00:00', null);
INSERT INTO public.objects (id, uuid_2, uuid_3, uuid_4, string_1, string_2, string_3, string_4, bool_1, bool_2, bool_3,
                            bool_4, float_32, float_64, int_8, int_16, int_32, int_64, time_1, time_2, time_3, time_4)
VALUES ('9c0c9146-cb6d-4184-80e4-acd26c2c1881', '3ad771d3-f617-4814-bda4-abb071ead70f', null,
        '26dd056a-fddf-4992-9e7f-1af09d610c5c', 'green', 'black', 'red', null, true, true, false, null, 0, null,
        '00000111', 0, 0, null, '2025-02-24 15:40:50.622572 +00:00', '2025-02-24 15:40:50.622572 +00:00',
        '1970-01-01 00:00:00.000000 +00:00', '2025-03-10 14:31:10.411012 +00:00');
INSERT INTO public.objects (id, uuid_2, uuid_3, uuid_4, string_1, string_2, string_3, string_4, bool_1, bool_2, bool_3,
                            bool_4, float_32, float_64, int_8, int_16, int_32, int_64, time_1, time_2, time_3, time_4)
VALUES ('1ad724f7-7d0e-4ac2-8723-7ec2ed36ad2f', '8366762c-e1a0-4bd1-bd43-c61c0605f5d9', null,
        '62549d95-ba11-450b-a4a2-0911afd3a73f', 'green', 'orange', 'red', null, true, false, false, null, 0, null,
        '00000001', 1, 0, null, '2025-03-10 16:24:31.217768 +00:00', '2025-03-10 16:24:31.217768 +00:00',
        '2025-03-10 16:24:31.217768 +00:00', '2025-03-10 21:00:00.000000 +00:00');

create table domains
(
    id      uuid                     default gen_random_uuid() not null
        primary key,
    name    text                                               not null,
    active  boolean                  default true              not null,
    created timestamp with time zone default now()             not null,
    updated timestamp with time zone default now()             not null,
    deleted timestamp with time zone
);

alter table domains
    owner to test;

create unique index domains_name_idx
    on domains (name)
    where (deleted IS NULL);

INSERT INTO public.domains (id, name, active, created, updated, deleted)
VALUES ('2eb40782-cf1d-4777-a5ac-fad81ea4f5a2', 'domain1', true, '2025-03-08 12:14:03.758198 +00:00',
        '2025-03-08 12:14:03.758198 +00:00', null);
