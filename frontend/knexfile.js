// Database configuration for MarvelQuant Trading Platform
// This file contains the database connection settings

module.exports = {
  development: {
    client: 'mysql',
    connection: {
      host: 'localhost',
      database: 'trading_platform',
      user: 'trading_user',
      password: 's2oF3i061E1n8u',
      charset: 'utf8'
    },
    pool: {
      min: 2,
      max: 10
    },
    migrations: {
      tableName: 'knex_migrations',
      directory: './migrations'
    },
    seeds: {
      directory: './seeds'
    }
  },
  production: {
    client: 'mysql',
    connection: {
      host: process.env.DB_HOST || 'localhost',
      database: process.env.DB_NAME || 'trading_platform',
      user: process.env.DB_USER || 'trading_user',
      password: process.env.DB_PASSWORD || 's2oF3i061E1n8u',
      charset: 'utf8'
    },
    pool: {
      min: 2,
      max: 10
    },
    migrations: {
      tableName: 'knex_migrations',
      directory: './migrations'
    }
  }
};
