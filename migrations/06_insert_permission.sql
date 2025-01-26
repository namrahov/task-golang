-- +migrate Up

insert into permissions(url, http_method, description)
values ('/v1/users/logout', 'GET', 'Logout');

insert into roles_permissions(role_id, permission_id)
values (1, 1),
       (2,1);
