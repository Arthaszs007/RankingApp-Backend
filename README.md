# go + gin + gorm + bcrypt + logrus + jwt+cors(midware) + uuid + redis + redisCloud

server:
support client to check players info to vote and rank list. users register , login, logout and JWT.

framework: gin

security: config method,host and header to check for authorized access.

redis: use for store the rank data with a expire, improve the data reading performance. work with the redis cloud.

auth: create and verify JWT to support automatic login, encript the senstive infomation with bcript.

database: postgreSQL (neon) with orm(gorm)

database design:
![failed](https://github.com/Arthaszs007/RankingApp-Backend/blob/master/DBR.png?raw=true)
