go+gin+gorm+bcrypt+logrus+jwt+cors(midware)+uuid+redis+redisCloud

server:
login to generate jwt and return it to client, visit server include a jwt in header with api to verify, crossing origin resource sharing to solve the cross-domain visit.

framework: gin

security: config method,host and header to check for authorized access.

redis: use for store the rank data with a expire, improve the data reading performance. work with the redis cloud.

auth: create and verify jwt to support automatic login, encript the senstive infomation with bcript.

database: connect to database (postgreSQL) by orm(gorm)

database design:
