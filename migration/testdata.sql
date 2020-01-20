INSERT INTO products(name,product_url) values('FIRST_AWESOME_PRODUCT','https://domain.dmn/wiki/spaces/FIRST_AWESOME_PRODUCT');
INSERT INTO products(name,product_url) values('NOT_SO_AWESOME_PRODUCT','https://domain.dmn/wiki/spaces/NOT_SO_AWESOME_PRODUCT');
INSERT INTO teams(name,team_url) values('DevOps','https://domain.dmn/wiki/spaces/DEVOPS/overview');
INSERT INTO teams(name,team_url) values('SRE','https://domain.dmn/wiki/spaces/sre/overview');
INSERT INTO services(id, name, service_level, service_type, service_status, product_id, team_id)
           values('elb', 'AWS ELB', 'P2', 'cloud_managed', 'LIVE', (select product_id from products where name='Infrastructure'),
                 (select team_id from teams where name='DevOps'));
          
INSERT INTO services(id, name, service_level, service_type, service_status, product_id, team_id)
           values('akamai', 'Akamai CDN', 'P2', 'third_party', 'LIVE', (select product_id from products where name='Infrastructure'),
                 (select team_id from teams where name='SRE'));
          
INSERT INTO services(id, name, service_level, service_type, service_status, product_id, team_id)
           values('cloudfront', 'Amazon CloudFront', 'P2', 'cloud_managed', 'LIVE', (select product_id from products where name='Infrastructure'),
                 (select team_id from teams where name='DevOps'));
          
INSERT INTO services(id, name, service_level, service_type, service_status, product_id, team_id)
           values('cognito', 'Amazon Cognito', 'P2', 'cloud_managed', 'LIVE', (select product_id from products where name='Infrastructure'),
                 (select team_id from teams where name='DevOps'));
          
INSERT INTO services(id, name, service_level, service_type, service_status, product_id, team_id)
           values('route53', 'Amazon Route 53 ', 'P2', 'cloud_managed', 'LIVE', (select product_id from products where name='Infrastructure'),
                 (select team_id from teams where name='DevOps'));
          
INSERT INTO services(id, name, service_level, service_type, service_status, product_id, team_id)
           values('ses', 'Amazon SES', 'P2', 'cloud_managed', 'LIVE', (select product_id from products where name='Infrastructure'),
                 (select team_id from teams where name='DevOps'));
          
INSERT INTO services(id, name, service_level, service_type, service_status, product_id, team_id)
           values('sqs', 'Amazon Simple Queue Service', 'P2', 'cloud_managed', 'LIVE', (select product_id from products where name='FIRST_AWESOME_PRODUCT'),
                 (select team_id from teams where name='DevOps'));
          
INSERT INTO services(id, name, service_level, service_type, service_status, product_id, team_id)
           values('s3', 'Amazon Simple Storage Service', 'P2', 'cloud_managed', 'LIVE', (select product_id from products where name='Infrastructure'),
                 (select team_id from teams where name='DevOps'));
          
