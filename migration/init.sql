CREATE TABLE IF NOT EXISTS products (
		product_id serial PRIMARY KEY,
		name VARCHAR (200) NOT NULL,
		owner varchar(160),
		product_url varchar
		);
CREATE TABLE IF NOT EXISTS teams (
		team_id serial PRIMARY KEY,
		name VARCHAR (200) NOT NULL,
		team_url varchar,
		teamlead varchar(160)
		);
CREATE TABLE IF NOT EXISTS services (
		service_id serial PRIMARY KEY,
		id VARCHAR (200) UNIQUE NOT NULL,
		name VARCHAR (200) NOT NULL,
		description text,
		team_id int,
		service_type varchar(80),
		service_level varchar(40),
		service_status varchar(40),
		sentry_id varchar,
		newrelic_id varchar,
		dashboard_url varchar,
		docurl varchar,
		svcurl varchar,
		product_id integer,
		CONSTRAINT service_product_id_fkey FOREIGN KEY (product_id)
        REFERENCES products (product_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION,
		CONSTRAINT service_team_id_fkey FOREIGN KEY (team_id)
        REFERENCES teams (team_id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION
);
