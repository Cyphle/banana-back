import convict from "convict";

const applicationConfig = convict({
  ENV: {
    doc: 'The application environment.',
    format: ['prd', 'dev', 'test'],
    default: 'dev',
    env: 'ENV'
  },
  IDP_SERVER_URL: {
    doc: 'The IDP url.',
    format: String,
    default: 'http://localhost:8182/realms/fastifyexample',
    env: 'IDP_SERVER_URL'
  },
  CLIENT_ID: {
    doc: 'The client ID for OIDC authorization code flow.',
    format: String,
    default: 'banana',
    env: 'CLIENT_ID',
  },
  CLIENT_SECRET: {
    doc: 'The client secret for OIDC authorization code flow.',
    format: String,
    default: 'banana-secret',
    env: 'CLIENT_SECRET'
  },
  REDIRECT_URI: {
    doc: 'The redirect URL for OIDC authorization code flow.',
    format: String,
    default: 'http://localhost:3000/callback',
  },
  FRONT_END_URL: {
    doc: 'The URL of the front-end',
    format: String,
    default: 'http://localhost:9000',
  }
});

export const getConfigValueOf = <T>(key: string, defaultValue: T): T => {
  // @ts-ignore
  const value = applicationConfig.get(key);

  if (!value) return defaultValue;

  return value as T;
}

// Load environment dependent configuration
// var env = config.get('env');
// config.loadFile('./config/' + env + '.json');