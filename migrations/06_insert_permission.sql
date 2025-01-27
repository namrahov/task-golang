-- +migrate Up

insert into permissions(url, http_method, description)
values ('/v1/users/logout', 'GET', 'Logout'),
       ('/v1/boards', 'POST', 'Create board'),
       ('/v1/boards/{id}/access', 'POST', 'Give access to a user'),
       ('/v1/boards/{userId}', 'GET', 'Get user boards');

insert into roles_permissions(role_id, permission_id)
values (1, 1),
       (2,1),
       (2, 2),
       (2, 3),
       (1, 4),
       (2, 4);
