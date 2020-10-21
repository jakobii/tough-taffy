drop table if exists permissions;
drop table if exists access;
drop table if exists permission_types;
drop table if exists resources;
drop table if exists members;
drop table if exists groups;
drop table if exists organizations;
drop table if exists tokens;
drop table if exists passwords;
drop table if exists users;

create table users (
    id uuid primary key not null,
    name text not null
);

create table passwords (
    id uuid primary key not null,
    created timestamptz not null,
    "user" uuid references users(id) not null,
    hash text not null,
    salt text not null
);

create table tokens (
    id uuid primary key not null,
    expires timestamptz not null,
    "user" uuid references users(id) not null
);

create table organizations (
    id uuid primary key not null,
    name text not null
);

-- all permissions should stem from a group.
-- groups are direct children to organizations.
-- the group table has some unconstained logic and should be queried with care.
-- this unconstained logic is bunndled in this table alone, and enables intuitive
-- fk's in other tables. e.g. the user and is_system feilds have special meaning
-- when used.
create table groups (
    id uuid primary key not null,
    
    org uuid references organizations(id) not null,
    name text not null,
    unique(org,name),

    description text not null,

    "user" uuid references users(id) null,
    -- when not null the group only references a single user.
    -- every users gets a unique group, of which memberships to other groups are based on.
    
    is_system bool not null
    -- when true, the group has special significants to the system, and should be managed
    -- more carefully, then a normal group.
    -- e.g.
    --     owner: god rights to the org, choose members wisely!
    --     billing: can manage billing.
    --     admin: can do anything, except delete the org and manage billing.
    --     users: defines org memberships. the only required group for each org.
);

create unique index org_user_unique on groups(org,"user") where "user" is not null;
-- ensures that and organization does not have duplicate user memberships.


create table members (
    id uuid primary key not null,
    "group" uuid references groups(id) not null,
    member uuid references groups(id) not null,
    unique(member,"group")

    -- care should be taken to prevent endless recusion when resolving
    -- memberships or permissions.
); 

create table resources (
    id uuid primary key not null,
    org uuid references organizations(id) not null,
    description text not null,
    name text not null,
    -- the api's asking for permission can have working domains, to prevent 
    -- colliding with other api permissions.
    -- when the value is "*" the expression applies to all domains.
    -- a domain could be an application name, an api name, or just some random
    -- string that only the caller would know.
    expression text not null,
    -- a regex that uniquely identifies a resource in a given domain.
    -- expression can also contain variable that should be filled
    -- before the expression is matched.
    -- vars should begin and end with '#' and conaint only letters, numbers, hypens, and underscores
    -- e.g. #self#, #org_id#, #a1-ff-de#, #1111#
    -- The '#' was chosen because it has no special significants in regex
    -- the '#' does have significants in a uri, but since the variables 
    -- should be evaluated before the match, so this should not be a problem.
    -- in the expression "^/users[?]org=#org_id#$", the "#org_id#" would be replaced, then
    -- the replaced expression would be used to match against. the expression suggestions that
    -- the caller would have some level of access to all users in a given organization.
    unique(name,expression)
);

create table permission_types (
	id int primary key,
	name text unique not null
);

create table access (
    id uuid primary key not null,
    "group" uuid references groups(id),
    resource uuid references resources(id)
);


create table permissions (
	access uuid references access(id) not null,
	type int references permission_types(id) not null,	
	primary key (access, type)
);
