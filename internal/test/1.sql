DROP TABLE IF EXISTS public.objects;

CREATE TABLE public.objects
(
    id         uuid PRIMARY KEY                  DEFAULT gen_random_uuid(),
    o_uuid_2   uuid                     NOT NULL,
    o_uuid_3   uuid                              DEFAULT gen_random_uuid(),
    o_uuid_4   uuid,
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
    o_int_8    char                     NOT NULL DEFAULT 'A',
    o_int_16   int2                     NOT NULL,
    o_int_32   int4                              DEFAULT 0,
    o_int_64   int8,
    o_time_1   timestamp with time zone NOT NULL DEFAULT now(),
    o_time_2   timestamp with time zone NOT NULL,
    o_time_3   timestamp with time zone          DEFAULT now(),
    o_time_4   timestamp with time zone
);
