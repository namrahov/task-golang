-- +migrate Up

insert into permissions(url, http_method, description)
values ('/v1/users/demo', 'POST', 'Create demo api');

insert into roles_permissions(role_id, permission_id)
values (2, 1);
