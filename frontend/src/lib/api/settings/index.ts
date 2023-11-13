import { Level } from '../log'

export interface Settings {
  ca_cert: string;
  ca_key: string;
  proxy_port: number;
  proxy_host: string;
  log_level: Level;
  dark_mode: boolean;
}

