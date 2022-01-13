create table merge_requests
(
    id                  uuid         not null,
    merge_request_i_id  varchar(255) not null,
    repository_url      varchar(255) not null,
    repository_id       varchar(255) not null,
    created_at          timestamp with time zone,
    updated_at          timestamp with time zone,
    deleted_at          timestamp with time zone,
    constraint merge_requests_pkey
        primary key (id)
);

create index idx_merge_requests_deleted_at
    on merge_requests (deleted_at);
